package main

import (
	"log"

	"github.com/yerden/go-dpdk/ethdev"
)

// EthdevCallback specifies callback to call on ethdev.Port.
type EthdevCallback interface {
	EthdevCall(ethdev.Port) error
}

type EthdevCallbackFunc func(ethdev.Port) error

func (f EthdevCallbackFunc) EthdevCall(pid ethdev.Port) error {
	return f(pid)
}

// EthdevConfig specifies information on how to configure ethdev.Port.
type EthdevConfig struct {
	Options  []ethdev.Option
	RxQueues uint16

	// Hooks to call after configuration
	OnConfig []EthdevCallback

	// RX queue config
	Pooler        RxqMempooler
	RxDescriptors uint16
	RxOptions     []ethdev.QueueOption
}

func ifaceName(pid ethdev.Port) string {
	var info ethdev.DevInfo
	if err := pid.InfoGet(&info); err != nil {
		panic(err)
	}
	return info.InterfaceName()
}

// Configure must be called on main lcore to configure ethdev.Port.
func (conf *EthdevConfig) Configure(pid ethdev.Port) error {
	if err := pid.DevConfigure(conf.RxQueues, 0, conf.Options...); err != nil {
		return err
	}

	if err := pid.PromiscEnable(); err != nil {
		return err
	}

	for qid := uint16(0); qid < conf.RxQueues; qid++ {
		//fmt.Printf("configuring rxq: %d@%d\n", pid, qid)
		mp, err := conf.Pooler.GetRxMempool(pid, qid)
		if err != nil {
			return err
		}
		if err := pid.RxqSetup(qid, conf.RxDescriptors, mp, conf.RxOptions...); err != nil {
			return err
		}
	}

	for i := range conf.OnConfig {
		if err := conf.OnConfig[i].EthdevCall(pid); err != nil {
			return err
		}
	}

	return nil
}

func printPortConfig(pid ethdev.Port) error {
	var info ethdev.DevInfo
	if err := pid.InfoGet(&info); err != nil {
		return err
	}

	log.Printf("port %d: nrxq=%d\n", pid, info.NbRxQueues())
	return nil
}
