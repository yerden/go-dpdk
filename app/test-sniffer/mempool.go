package main

import (
	"fmt"
	"os"

	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mempool"
)

// RxqMempooler returns a mempool for specified port and rx queue.
type RxqMempooler interface {
	GetRxMempool(pid ethdev.Port, qid uint16) (*mempool.Mempool, error)
	FreeMempools()
}

// mpPerPort implements RxqMempooler.
type mpPerPort struct {
	pools []*mempool.Mempool
}

// NewMempoolPerPort creates new RxqMempooler. It creates a mempool
// per each valid port number.
//
// prefix is used to construct a name for created mempools.
//
// Mempools are created from src. Closest to port NUMA node is
// prepended to the list of specified mempool options. If you
// deliberately want to set NUMA node, enlist it in opts.
func NewMempoolPerPort(prefix string, src mempool.Factory, opts ...mempool.Option) (RxqMempooler, error) {
	mpp := &mpPerPort{
		pools: make([]*mempool.Mempool, ethdev.CountTotal()),
	}

	var err error
	for i := range mpp.pools {
		if pid := ethdev.Port(i); pid.IsValid() {
			name := fmt.Sprintf("%s_%d", prefix, i)

			// initially, we set the socket of mempool to the closest one
			// user may specify desired socket if needed
			options := append([]mempool.Option{
				mempool.OptSocket(pid.SocketID()),
			}, opts...)

			if mpp.pools[i], err = src.NewMempool(name, options); err != nil {
				return nil, err
			}
		}
	}

	return mpp, nil
}

// GetRxMempool implements RxqMempooler interface.
func (mpp *mpPerPort) GetRxMempool(pid ethdev.Port, qid uint16) (*mempool.Mempool, error) {
	if !pid.IsValid() {
		return nil, os.ErrInvalid
	}

	return mpp.pools[pid], nil
}

// FreeMempools implements RxqMempooler interface.
func (mpp *mpPerPort) FreeMempools() {
	for _, mp := range mpp.pools {
		mp.Free()
	}
}
