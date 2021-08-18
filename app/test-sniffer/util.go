package main

/*
#include <rte_ethdev.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"log"
	"os"

	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/ethdev/flow"
)

func doOnMain(fn func() error) error {
	var err error
	e := eal.ExecOnMain(func(*eal.LcoreCtx) {
		err = fn()
	})
	if e != nil {
		return e
	}

	return err
}

// generate queues array [0..n)
func queuesSeq(n int) []uint16 {
	q := make([]uint16, n)

	for i := range q {
		q[i] = uint16(i)
	}

	return q
}

func failOnErr(err error) {
	if err != nil {
		var e *eal.ErrLcorePanic
		if errors.As(err, &e) {
			e.FprintStack(os.Stdout)
		}
		log.Fatal(err)
	}
}

func driverName(pid ethdev.Port) string {
	var devInfo ethdev.DevInfo
	if err := pid.InfoGet(&devInfo); err != nil {
		return ""
	}
	return devInfo.DriverName()
}

func rssEthVlanIPv4(pid ethdev.Port) (*flow.Flow, error) {
	attr := &flow.Attr{Ingress: true}

	pattern := []flow.Item{
		{Spec: flow.ItemTypeEth},  // Ethernet
		{Spec: flow.ItemTypeVlan}, // VLAN
		{Spec: flow.ItemTypeIPv4}, // IPv4
	}

	actions := []flow.Action{
		&flow.ActionRSS{
			Types: C.ETH_RSS_IPV4,
			Func:  flow.HashFunctionSymmetricToeplitz,
		},
	}

	e := &flow.Error{}
	if err := flow.Validate(pid, attr, pattern, actions, e); err == nil {
		if f, err := flow.Create(pid, attr, pattern, actions, e); err == nil {
			return f, nil
		}
	}

	return nil, e
}

func mlxRssEthVlanIPv4(pid ethdev.Port) (*flow.Flow, error) {
	attr := &flow.Attr{Ingress: true}

	pattern := []flow.Item{
		{Spec: flow.ItemTypeEth},  // Ethernet
		{Spec: flow.ItemTypeVlan}, // VLAN
		{Spec: flow.ItemTypeIPv4}, // IPv4
	}

	var info ethdev.DevInfo
	if err := pid.InfoGet(&info); err != nil {
		return nil, err
	}

	actions := []flow.Action{
		&flow.ActionRSS{
			Types:  C.ETH_RSS_IPV4,
			Key:    bytes.Repeat([]byte{0x6D, 0x5A}, 20),
			Queues: queuesSeq(int(info.NbRxQueues())),
			Func:   flow.HashFunctionToeplitz,
		},
	}

	e := &flow.Error{}
	if err := flow.Validate(pid, attr, pattern, actions, e); err == nil {
		if f, err := flow.Create(pid, attr, pattern, actions, e); err == nil {
			return f, nil
		}
	}

	return nil, e
}

type RssConfig struct {
}

func (c *RssConfig) EthdevCall(pid ethdev.Port) error {
	var err error
	switch driverName(pid) {
	case "mlx5_pci":
		_, err = mlxRssEthVlanIPv4(pid)
	case "net_af_packet":
		fallthrough
	case "net_ice":
		_, err = rssEthVlanIPv4(pid)
	default:
	}

	return err
}
