package main

import (
	"errors"
	"log"
	"syscall"

	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/util"
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

	// Flow Control Mode
	FcMode uint32
}

func ifaceName(pid ethdev.Port) string {
	name, _ := pid.Name()
	return name
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

	var fc ethdev.FcConf

	if err := pid.FlowCtrlGet(&fc); err == nil {
		fc.SetMode(conf.FcMode)
		if err := pid.FlowCtrlSet(&fc); err != nil {
			return util.ErrWrapf(err, "FlowCtrlSet")
		}

		log.Printf("pid=%d: Flow Control set to %d", pid, conf.FcMode)
	} else if !errors.Is(err, syscall.ENOTSUP) {
		return util.ErrWrapf(err, "FlowCtrlGet")
	}

	for i := range conf.OnConfig {
		if err := conf.OnConfig[i].EthdevCall(pid); err != nil {
			return util.ErrWrapf(err, "OnConfig %d: %v", i, conf.OnConfig[i])
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
