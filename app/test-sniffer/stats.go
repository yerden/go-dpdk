package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yerden/go-dpdk/ethdev"
)

var statsInt = flag.Duration("statsInt", 10*time.Second, "Specify interval between collecting statistics")

const (
	// namespace for metrics
	statsNamespace = "dpdk"
)

type portStat struct {
	pid ethdev.Port

	// ethdev basic statistics
	basic struct {
		prevStats, curStats ethdev.Stats

		iPackets prometheus.Counter
		oPackets prometheus.Counter
		iBytes   prometheus.Counter
		oBytes   prometheus.Counter
		iMissed  prometheus.Counter
		iErrors  prometheus.Counter
		oErrors  prometheus.Counter
		rxNoMbuf prometheus.Counter
	}
	// ethdev extended statistics
	extended struct {
		ids               []uint64
		prevVals, curVals []uint64

		counters []prometheus.Counter
	}
}

type Stats struct {
	ports []portStat
}

func NewStats(reg prometheus.Registerer, pids []ethdev.Port) (*Stats, error) {
	xVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: statsNamespace,
			Name:      "xstats",
		},
		[]string{"index", "driver", "name", "xstat_name"},
	)
	reg.MustRegister(xVec)

	var portStats []portStat
	for _, pid := range pids {
		var devInfo ethdev.DevInfo
		if err := pid.InfoGet(&devInfo); err != nil {
			return nil, err
		}

		devName, err := pid.Name()
		if err != nil {
			return nil, err
		}

		labels := prometheus.Labels{
			"index":  fmt.Sprint(pid),
			"driver": devInfo.DriverName(),
			"name":   devName,
		}
		newCounter := func(name string) prometheus.Counter {
			c := prometheus.NewCounter(prometheus.CounterOpts{
				Namespace:   statsNamespace,
				Name:        name,
				ConstLabels: labels,
			})
			reg.MustRegister(c)
			return c
		}

		var ps portStat
		ps.pid = pid

		// basic stats
		b := &ps.basic
		b.iPackets = newCounter("ipackets")
		b.oPackets = newCounter("opackets")
		b.iBytes = newCounter("ibytes")
		b.oBytes = newCounter("obytes")
		b.iMissed = newCounter("imissed")
		b.iErrors = newCounter("ierrors")
		b.oErrors = newCounter("oerrors")
		b.rxNoMbuf = newCounter("rxnombuf")

		// xstats
		xstatNameIDs, err := pid.XstatNameIDs()
		if err != nil {
			return nil, err
		}

		x := &ps.extended
		x.prevVals = make([]uint64, len(xstatNameIDs))
		x.curVals = make([]uint64, len(xstatNameIDs))
		for id, xstatName := range xstatNameIDs {
			// here labels changes, so can't be used for basic stats later, but
			// it's created at each iteration
			labels["xstat_name"] = xstatName
			x.counters = append(x.counters, xVec.With(labels))
			x.ids = append(x.ids, id)
		}

		portStats = append(portStats, ps)
	}

	return &Stats{portStats}, nil
}

func (s *Stats) Report() error {
	var diffStats ethdev.Stats
	add := func(c prometheus.Counter, delta uint64) {
		if delta > 0 {
			c.Add(float64(delta))
		}
	}

	for i := range s.ports {
		ps := &s.ports[i]
		pid := ps.pid

		// ethdev basic statistics
		b := &ps.basic
		if err := pid.StatsGet(&b.curStats); err != nil {
			return err
		}
		b.curStats.Diff(&b.prevStats, &diffStats)
		gd := diffStats.Cast()
		add(b.iPackets, gd.Ipackets)
		add(b.oPackets, gd.Opackets)
		add(b.iBytes, gd.Ibytes)
		add(b.oBytes, gd.Obytes)
		add(b.iMissed, gd.Imissed)
		add(b.iErrors, gd.Ierrors)
		add(b.oErrors, gd.Oerrors)
		add(b.rxNoMbuf, gd.RxNoMbuf)
		b.prevStats, b.curStats = b.curStats, b.prevStats

		// ethdev extended statistics
		x := &ps.extended
		if _, err := pid.XstatGetByID(x.ids, x.curVals); err != nil {
			return err
		}
		for i, cur := range x.curVals {
			if prev := x.prevVals[i]; prev < cur {
				x.counters[i].Add(float64(cur - prev))
			}
		}
		x.prevVals, x.curVals = x.curVals, x.prevVals
	}

	return nil
}
