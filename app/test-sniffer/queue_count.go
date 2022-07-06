package main

import (
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mbuf"
	"github.com/yerden/go-dpdk/util"
)

type PacketBytes struct {
	Packets prometheus.Counter
	Bytes   prometheus.Counter
}

type QueueCounter struct {
	RX PacketBytes
}

type QueueCounterReporter struct {
	reg    prometheus.Registerer
	mtx    sync.Mutex
	lcores []*QueueCounter
}

func (qcr *QueueCounterReporter) getLcoresCounter(lcore uint) *QueueCounter {
	if extra := int(lcore) - len(qcr.lcores) + 1; extra > 0 {
		qcr.lcores = append(qcr.lcores, make([]*QueueCounter, extra)...)
	}

	p := new(QueueCounter)
	qcr.lcores[lcore] = p
	return p
}

func (qcr *QueueCounterReporter) Register(pid ethdev.Port, qid uint16) *QueueCounter {
	qcr.mtx.Lock()
	defer qcr.mtx.Unlock()

	qc := qcr.getLcoresCounter(eal.LcoreID())
	name, err := pid.Name()
	if err != nil {
		panic(util.ErrWrapf(err, "no name for pid=%d", pid))
	}

	labels := prometheus.Labels{
		"index": strconv.FormatUint(uint64(pid), 10),
		"queue": strconv.FormatUint(uint64(qid), 10),
		"name":  name,
	}
	newCounter := func(name string) prometheus.Counter {
		c := prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:   statsNamespace,
			Subsystem:   "rxq",
			Name:        name,
			ConstLabels: labels,
		})
		qcr.reg.MustRegister(c)
		return c
	}

	qc.RX.Packets = newCounter("rx_packets")
	qc.RX.Bytes = newCounter("rx_bytes")

	return qc
}

func (qc *QueueCounter) Incr(pkts []*mbuf.Mbuf) {
	dataLen := uint64(0)
	for i := range pkts {
		dataLen += uint64(len(pkts[i].Data()))
	}
	qc.RX.Packets.Add(float64(len(pkts)))
	qc.RX.Bytes.Add(float64(dataLen))
}
