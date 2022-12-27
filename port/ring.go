package port

/*
#include <rte_config.h>
#include <rte_errno.h>

#include <rte_port_ring.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/ring"
)

// compile time checks
var _ = []InParams{
	&RingRx{},
}

var _ = []OutParams{
	&RingTx{},
}

// RingRx is an input port built on top of pre-initialized RTE ring.
type RingRx struct {
	// Underlying ring
	*ring.Ring

	// Set if specified ring is multi consumer.
	Multi bool
}

// InOps implements InParams interface.
func (rd *RingRx) InOps() *InOps {
	if !rd.Multi {
		return (*InOps)(&C.rte_port_ring_reader_ops)
	}
	return (*InOps)(&C.rte_port_ring_multi_reader_ops)
}

// Transform implements common.Transformer interface.
func (rd *RingRx) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return common.TransformPOD(alloc, &C.struct_rte_port_ring_reader_params{
		ring: (*C.struct_rte_ring)(unsafe.Pointer(rd.Ring)),
	})
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

// OutOps implements OutParams interface.
func (wr *RingTx) OutOps() *OutOps {
	ops := []*C.struct_rte_port_out_ops{
		&C.rte_port_ring_writer_ops,
		&C.rte_port_ring_multi_writer_ops,
		&C.rte_port_ring_writer_nodrop_ops,
		&C.rte_port_ring_multi_writer_nodrop_ops,
	}

	if wr.Multi {
		ops = ops[1:]
	}

	if wr.NoDrop {
		ops = ops[2:]
	}

	return (*OutOps)(ops[0])
}

// Transform implements common.Transformer interface.
func (wr *RingTx) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return common.TransformPOD(alloc, &C.struct_rte_port_ring_writer_nodrop_params{
		ring:        (*C.struct_rte_ring)(unsafe.Pointer(wr.Ring)),
		tx_burst_sz: C.uint32_t(wr.TxBurstSize),
		n_retries:   C.uint32_t(wr.Retries),
	})
}
