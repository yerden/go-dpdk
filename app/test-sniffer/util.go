package main

/*
#include <rte_ethdev.h>
*/
import "C"

import (
	"errors"
	"log"
	"os"

	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/ethdev/flow"
)

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
