package main

// #include <rte_ring.h>
import "C"

import (
	"context"
	"log"
	"strings"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/memzone"
	"github.com/yerden/go-dpdk/ring"
)

const (
	portNameLbl  = "port_name"
	memzoneLbl   = "memzone_name"
	macAddrLbl   = "mac_addr"
	drvNameLbl   = "driver_name"
	ifaceNameLbl = "interface_name"
	xstatNameLbl = "xstat_name"
)
const (
	namespace    = "dpdk_exporter"
	ringMzPrefix = "RG_"
)

type Metrics struct {
	EthDev *EthDevMetrics
	Ring   *RingMetrics
}

func NewMetrics() (m *Metrics, err error) {
	ethDev, err := NewEthDevMetrics()
	if err != nil {
		return nil, err
	}

	m = &Metrics{
		EthDev: ethDev,
		Ring:   NewRingMetrics(),
	}
	return
}

func (m *Metrics) Collect() {
	if err := m.EthDev.Collect(); err != nil {
		log.Printf("collect eth dev metrics: %v", err)
	}
	if err := m.Ring.Collect(); err != nil {
		log.Printf("collect ring metrics: %v", err)
	}
}

func (m *Metrics) StartCollecting(ctx context.Context) {
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			eal.ExecOnMain(func(lc *eal.LcoreCtx) {
				m.Collect()
			})
		}
	}
}

func setBooleanGauge(g prometheus.Gauge, v bool) {
	if v {
		g.Set(1)
	} else {
		g.Set(0)
	}
}

type EthDevMetrics struct {
	EthXstats *EthXstats

	AvailablePorts prometheus.Gauge

	// link-level metrics of an Ethernet port.
	Link *LinkMetrics

	SocketID   *prometheus.GaugeVec
	NbRxQueues *prometheus.GaugeVec
	NbTxQueues *prometheus.GaugeVec
	RetaSize   *prometheus.GaugeVec
	Info       *prometheus.GaugeVec
}

func NewEthDevMetrics() (*EthDevMetrics, error) {
	var m EthDevMetrics

	var err error
	m.EthXstats, err = NewEthXstats()
	if err != nil {
		return nil, err
	}

	labelNames := []string{portNameLbl}
	const subsystem = "eth_dev"
	m.AvailablePorts = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "available_ports",
	})
	m.Link = NewLinkMetrics(labelNames)
	m.SocketID = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "socket_id",
	}, labelNames)
	m.NbRxQueues = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "nb_rx_queues",
	}, labelNames)
	m.NbTxQueues = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "nb_tx_queues",
	}, labelNames)
	m.RetaSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "redirection_table_size",
	}, labelNames)
	m.Info = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "info",
	}, append(labelNames, macAddrLbl, drvNameLbl, ifaceNameLbl))

	return &m, nil
}

func (m *EthDevMetrics) Collect() error {
	m.AvailablePorts.Set(float64(ethdev.CountAvail()))

	for _, port := range ethdev.ValidPorts() {
		name, err := port.Name()
		if err != nil {
			log.Printf("get device name: %v, port id: %d", err, port)
			continue
		}
		labels := prometheus.Labels{portNameLbl: name}

		if err := m.Link.Collect(port, labels); err != nil {
			log.Printf("retrieve link status: %v, port id: %d, name: %s", err, port, name)
		}

		m.SocketID.With(labels).Set(float64(port.SocketID()))

		var macAddr ethdev.MACAddr
		if err := port.MACAddrGet(&macAddr); err != nil {
			log.Printf("retrieve mac address: %v, port id: %d, name: %s", err, port, name)
			continue
		}
		var devInfo ethdev.DevInfo
		if err := port.InfoGet(&devInfo); err != nil {
			log.Printf("retrieve device info: %v, port id: %d, name: %s", err, port, name)
			continue
		}

		m.NbRxQueues.With(labels).Set(float64(devInfo.NbRxQueues()))
		m.NbTxQueues.With(labels).Set(float64(devInfo.NbTxQueues()))
		m.RetaSize.With(labels).Set(float64(devInfo.RetaSize()))

		m.Info.With(prometheus.Labels{
			portNameLbl:  name,
			macAddrLbl:   macAddr.String(),
			drvNameLbl:   devInfo.DriverName(),
			ifaceNameLbl: devInfo.InterfaceName(),
		}).Set(1)
	}

	if err := m.EthXstats.Collect(); err != nil {
		log.Printf("retrieve extended statistics of eth dev: %v", err)
	}

	return nil
}

type LinkMetrics struct {
	AutoNeg *prometheus.GaugeVec
	Duplex  *prometheus.GaugeVec
	Speed   *prometheus.GaugeVec
	Status  *prometheus.GaugeVec
}

func NewLinkMetrics(labelNames []string) *LinkMetrics {
	var m LinkMetrics

	const subsystem = "eth_dev_link"
	m.AutoNeg = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "auto_negotiation",
	}, labelNames)
	m.Duplex = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "full_duplex",
	}, labelNames)
	m.Speed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "speed_mbps",
	}, labelNames)
	m.Status = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "status",
	}, labelNames)

	return &m
}

func (m *LinkMetrics) Collect(port ethdev.Port, labels prometheus.Labels) error {
	link, err := port.EthLinkGet()
	if err == syscall.ENOTSUP {
		// simply not exporting
		return nil
	} else if err != nil {
		return err
	}

	setBooleanGauge(m.AutoNeg.With(labels), link.AutoNeg())
	setBooleanGauge(m.Duplex.With(labels), link.Duplex())
	m.Speed.With(labels).Set(float64(link.Speed()))
	setBooleanGauge(m.Status.With(labels), link.Status())

	return nil
}

type ethXstat struct {
	cnt             []prometheus.Counter
	ids, prev, next []uint64
}

type EthXstats struct {
	stats map[ethdev.Port]*ethXstat
}

func NewEthXstats() (*EthXstats, error) {
	labelNames := []string{portNameLbl, xstatNameLbl}
	vec := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "eth_dev_xstats",
		Help:      "Extended NIC statistic counters",
	}, labelNames)

	stats := make(map[ethdev.Port]*ethXstat)
	for _, port := range ethdev.ValidPorts() {
		pname, err := port.Name()
		if err != nil {
			return nil, err
		}

		nameIDs, err := port.XstatNameIDs()
		if err != nil {
			return nil, err
		}

		ex := &ethXstat{
			prev: make([]uint64, len(nameIDs)),
			next: make([]uint64, len(nameIDs)),
		}
		for id, xname := range nameIDs {
			ex.ids = append(ex.ids, id)
			ex.cnt = append(ex.cnt, vec.WithLabelValues(pname, xname))
		}

		stats[port] = ex
	}

	return &EthXstats{stats}, nil
}

func (exs *EthXstats) Collect() error {
	for port, ex := range exs.stats {
		if _, err := port.XstatGetByID(ex.ids, ex.next); err != nil {
			return err
		}

		// calculate delta
		for i, p := range ex.prev {
			if n := ex.next[i]; n > p {
				ex.cnt[i].Add(float64(n - p))
			}
		}

		// swap old and new stats
		ex.prev, ex.next = ex.next, ex.prev
	}

	return nil
}

type RingMetrics struct {
	MemzoneLen      *prometheus.GaugeVec
	MemzonePageSize *prometheus.GaugeVec
	MemzoneSocketID *prometheus.GaugeVec
	Cap             *prometheus.GaugeVec
	Count           *prometheus.GaugeVec
}

func NewRingMetrics() *RingMetrics {
	var m RingMetrics

	labelNames := []string{memzoneLbl}
	const subsystem = "ring"
	m.MemzoneLen = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "memzone_len",
	}, labelNames)
	m.MemzonePageSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "memzone_page_size",
	}, labelNames)
	m.MemzoneSocketID = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "memzone_socket_id",
	}, labelNames)
	m.Cap = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "capacity",
	}, labelNames)
	m.Count = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "count",
	}, labelNames)

	return &m
}

func (m *RingMetrics) Collect() error {
	memzone.Walk(func(mz *memzone.Memzone) {
		name := mz.Name()
		if !strings.HasPrefix(name, C.RTE_RING_MZ_PREFIX) {
			return
		}

		labels := prometheus.Labels{memzoneLbl: name}

		m.MemzoneLen.With(labels).Set(float64(mz.Len()))
		m.MemzonePageSize.With(labels).Set(float64(mz.HugePageSz()))
		m.MemzoneSocketID.With(labels).Set(float64(mz.SocketID()))

		r := (*ring.Ring)(mz.Addr())
		m.Cap.With(labels).Set(float64(r.Cap()))
		m.Count.With(labels).Set(float64(r.Count()))
	})
	return nil
}
