package main

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/stats/v4"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mbuf"
	"github.com/yerden/go-dpdk/util"
)

type PacketBytes struct {
	Packets uint64 `metric:"packets" type:"counter"`
	Bytes   uint64 `metric:"bytes" type:"counter"`
}

type QueueCounter struct {
	PortID  string `tag:"index"`
	Name    string `tag:"name"`
	QueueID string `tag:"queue"`

	RX PacketBytes `metric:"rx"`
	TX PacketBytes `metric:"tx"`
}

type QueueCounterReporter struct {
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

	qc := qcr.getLcoresCounter(eal.LcoreID())
	name, err := pid.Name()
	if err != nil {
		panic(util.ErrWrapf(err, "no name for pid=%d", pid))
	}
	qc.PortID = strconv.FormatUint(uint64(pid), 10)
	qc.QueueID = strconv.FormatUint(uint64(qid), 10)
	qc.Name = name

	qcr.mtx.Unlock()

	return qc
}

func (qc *QueueCounter) Incr(pkts []*mbuf.Mbuf) {
	dataLen := uint64(0)
	for i := range pkts {
		dataLen += uint64(len(pkts[i].Data()))
	}
	atomic.AddUint64(&qc.RX.Packets, uint64(len(pkts)))
	atomic.AddUint64(&qc.RX.Bytes, uint64(dataLen))
}

func (qcr *QueueCounterReporter) ReportAt(t time.Time, eng *stats.Engine) {
	data := []QueueCounter{}

	qcr.mtx.Lock()

	for _, qc := range qcr.lcores {
		if qc == nil {
			continue
		}

		var newQC QueueCounter
		newQC.PortID = qc.PortID
		newQC.QueueID = qc.QueueID
		newQC.Name = qc.Name

		newQC.RX.Packets = atomic.SwapUint64(&qc.RX.Packets, 0)
		newQC.RX.Bytes = atomic.SwapUint64(&qc.RX.Bytes, 0)

		data = append(data, newQC)
	}

	qcr.mtx.Unlock()

	eng.ReportAt(t, data)
}
