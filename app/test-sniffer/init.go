package main

/*
#include <rte_ethdev.h>
*/
import "C"

import (
	"bytes"
	"flag"
	"fmt"

	"github.com/segmentio/stats/v4"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mempool"
)

type CmdMempool struct{}

var mbufElts = flag.Int("poolMbufs", 100000, "Specify number of mbufs in mempool")
var mbufSize = flag.Int("dataRoomSize", 3000, "Specify mbuf size in mempool")
var poolCache = flag.Int("poolCache", 512, "Specify amount of mbufs in per-lcore cache")
var rxQueues = flag.Int("nrxq", 32, "Specify number of RX queues per port")
var rxDesc = flag.Int("ndesc", 256, "Specify number of RX queues per port")

func (*CmdMempool) NewMempool(name string, opts []mempool.Option) (*mempool.Mempool, error) {
	return mempool.CreateMbufPool(name, uint32(*mbufElts), uint16(*mbufSize), opts...)
}

type App struct {
	RxqMempooler
	Stats *Stats
	Ports []ethdev.Port
	Work  map[uint]PortQueue
	QCR   *QueueCounterReporter
}

func NewApp(eng *stats.Engine) (*App, error) {
	var app *App
	return app, doOnMain(func() error {
		var err error
		app, err = newApp(eng)
		return err
	})
}

func newApp(eng *stats.Engine) (*App, error) {
	rxqPools, err := NewMempoolPerPort("mbuf_pool", &CmdMempool{},
		mempool.OptCacheSize(uint32(*poolCache)),
		mempool.OptOpsName("lf_stack"),
		mempool.OptPrivateDataSize(64), // for each mbuf
	)

	if err != nil {
		return nil, err
	}

	rssConf := &ethdev.RssConf{
		Key: bytes.Repeat([]byte{0x6D, 0x5A}, 20),
		Hf:  C.RTE_ETH_RSS_IP,
	}

	ethdevCfg := &EthdevConfig{
		Options: []ethdev.Option{
			ethdev.OptRss(*rssConf),
			ethdev.OptRxMode(ethdev.RxMode{
				MqMode: C.RTE_ETH_MQ_RX_RSS,
			}),
		},
		RxQueues: uint16(*rxQueues),
		OnConfig: []EthdevCallback{
			EthdevCallbackFunc((ethdev.Port).Start),
			// &RssConfig{rssConf},
		},
		Pooler:        rxqPools,
		RxDescriptors: uint16(*rxDesc),
		FcMode:        fcMode.Mode,
	}

	ports := make([]ethdev.Port, 0, ethdev.CountTotal())

	for i := 0; i < cap(ports); i++ {
		if pid := ethdev.Port(i); pid.IsValid() {
			ports = append(ports, pid)
		}
	}

	for i := range ports {
		fmt.Printf("configuring port %d: %s... ", ports[i], ifaceName(ports[i]))
		if err := ethdevCfg.Configure(ports[i]); err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println("OK")
	}

	metrics, err := NewStats(eng, ports)
	if err != nil {
		return nil, err
	}

	work, err := DistributeQueues(ports, eal.LcoresWorker())
	if err != nil {
		return nil, err
	}

	app := &App{
		RxqMempooler: rxqPools,
		Ports:        ports,
		Stats:        metrics,
		Work:         work,
		QCR:          &QueueCounterReporter{},
	}

	return app, nil
}
