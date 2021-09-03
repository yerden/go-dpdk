package main

import (
	"fmt"
	"time"

	"github.com/segmentio/stats/v4"
	"github.com/yerden/go-dpdk/ethdev"
)

type portStat struct {
	pid                ethdev.Port
	eng                *stats.Engine
	oldStats           *ethdev.Stats
	newStats           *ethdev.Stats
	xstatItems         []statItem
	delta, old, xstats []ethdev.Xstat
}

type Stats struct {
	ports []portStat
}

type statItem struct {
	Name  string `tag:"xstat_name"`
	Value uint64 `metric:"xstats" type:"counter"`
}

func NewStats(eng *stats.Engine, pids []ethdev.Port) (*Stats, error) {
	s := &Stats{}

	// xstat
	for _, pid := range pids {
		ps := portStat{pid: pid}

		names, err := pid.XstatNames()

		if err != nil {
			return nil, err
		}

		var devInfo ethdev.DevInfo
		if err := pid.InfoGet(&devInfo); err != nil {
			return nil, err
		}

		devName, err := pid.Name()
		if err != nil {
			return nil, err
		}

		ps.newStats = new(ethdev.Stats)
		ps.oldStats = new(ethdev.Stats)

		ps.xstats = make([]ethdev.Xstat, len(names))
		ps.xstatItems = make([]statItem, len(names))

		for i := range ps.xstatItems {
			item := &ps.xstatItems[i]
			item.Name = names[i].String()
		}

		ps.eng = eng.WithTags(
			stats.T("index", fmt.Sprint(pid)),
			stats.T("driver", devInfo.DriverName()),
			stats.T("name", devName),
		)

		s.ports = append(s.ports, ps)
	}

	return s, nil
}

func (s *Stats) ReportAt(now time.Time) error {
	var diffStats ethdev.Stats
	for i := range s.ports {
		ps := &s.ports[i]
		pid := ps.pid

		if err := pid.StatsGet(ps.newStats); err != nil {
			return err
		}

		ps.newStats.Diff(ps.oldStats, &diffStats)
		ps.eng.ReportAt(now, diffStats.Cast())
		ps.oldStats, ps.newStats = ps.newStats, ps.oldStats

		// reset delta
		for j := range ps.xstatItems {
			ps.xstatItems[j].Value = 0
		}

		for j := range ps.delta {
			ps.delta[j].Value = 0
		}

		// get new stats
		n, err := pid.XstatsGet(ps.xstats)
		if err != nil {
			return err
		}

		newXstat := ps.xstats[:n]
		ps.old, ps.delta = ethdev.XstatDiff(newXstat, ps.old, ps.delta)

		for j := range ps.delta {
			id := ps.delta[j].Index
			ps.xstatItems[id].Value = ps.delta[j].Value
		}

		ps.eng.ReportAt(now, ps.xstatItems)
	}

	return nil
}
