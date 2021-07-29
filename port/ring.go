package port

/*
#include <rte_config.h>
#include <rte_errno.h>

#include <rte_port_ring.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/ring"
)

// compile time checks
var _ = []RxFactory{
	&RingRx{},
}

var _ = []TxFactory{
	&RingTx{},
}

// RingRx is an input port built on top of pre-initialized single
// consumer ring.
type RingRx struct {
	// Underlying ring
	*ring.Ring

	// Set if specified ring is multi consumer.
	Multi bool
}

// CreateRx implements RxFactory interface.
func (rd *RingRx) CreateRx(socket int) (*Rx, error) {
	rx := &Rx{}

	if !rd.Multi {
		rx.ops = &C.rte_port_ring_reader_ops
	} else {
		rx.ops = &C.rte_port_ring_multi_reader_ops
	}

	params := &C.struct_rte_port_ring_reader_params{
		ring: (*C.struct_rte_ring)(unsafe.Pointer(rd.Ring)),
	}

	return rx, rx.doCreate(socket, unsafe.Pointer(params))
}

// RingTx is an output port built on top of pre-initialized single
// producer ring.
type RingTx struct {
	// Underlying ring
	*ring.Ring

	// Recommended burst size for ring operations.
	TxBurstSize uint32

	// Set if specified ring is multi producer.
	Multi bool

	// If NoDrop set writer makes Retries attempts to write packets to
	// ring.
	NoDrop bool

	// If NoDrop set and Retries is 0, number of retries is unlimited.
	Retries uint32
}

// CreateTx implements TxFactory interface.
func (wr *RingTx) CreateTx(socket int) (*Tx, error) {
	tx := &Tx{}

	var err error
	if wr.NoDrop {
		if wr.Multi {
			tx.ops = &C.rte_port_ring_multi_writer_nodrop_ops
		} else {
			tx.ops = &C.rte_port_ring_writer_nodrop_ops
		}

		params := &C.struct_rte_port_ring_writer_nodrop_params{
			ring:        (*C.struct_rte_ring)(unsafe.Pointer(wr.Ring)),
			tx_burst_sz: C.uint32_t(wr.TxBurstSize),
			n_retries:   C.uint32_t(wr.Retries),
		}
		err = tx.doCreate(socket, unsafe.Pointer(params))
	} else {
		if wr.Multi {
			tx.ops = &C.rte_port_ring_multi_writer_ops
		} else {
			tx.ops = &C.rte_port_ring_writer_ops
		}

		params := &C.struct_rte_port_ring_writer_params{
			ring:        (*C.struct_rte_ring)(unsafe.Pointer(wr.Ring)),
			tx_burst_sz: C.uint32_t(wr.TxBurstSize),
		}
		err = tx.doCreate(socket, unsafe.Pointer(params))
	}

	return tx, err
}
