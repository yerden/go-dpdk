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
var _ = []ConfigIn{
	&RingIn{},
}

var _ = []ConfigOut{
	&RingOut{},
}

// RingIn is an input port built on top of pre-initialized single
// consumer ring.
type RingIn struct {
	// Underlying ring
	*ring.Ring

	// Set if specified ring is multi consumer.
	Multi bool
}

// Ops implements ConfigIn interface.
func (rd *RingIn) Ops() *InOps {
	if !rd.Multi {
		return (*InOps)(&C.rte_port_ring_reader_ops)
	}
	return (*InOps)(&C.rte_port_ring_multi_reader_ops)
}

// Arg implements ConfigIn interface.
func (rd *RingIn) Arg(mem common.Allocator) *InArg {
	var rc *C.struct_rte_port_ring_reader_params
	common.MallocT(mem, &rc)
	rc.ring = (*C.struct_rte_ring)(unsafe.Pointer(rd.Ring))
	return (*InArg)(unsafe.Pointer(rc))
}

// RingOut is an output port built on top of pre-initialized single
// producer ring.
type RingOut struct {
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

// Ops implements ConfigOut interface.
func (wr *RingOut) Ops() *OutOps {
	switch {
	case wr.Multi && wr.NoDrop:
		return (*OutOps)(&C.rte_port_ring_multi_writer_nodrop_ops)
	case wr.Multi:
		return (*OutOps)(&C.rte_port_ring_multi_writer_ops)
	case wr.NoDrop:
		return (*OutOps)(&C.rte_port_ring_writer_nodrop_ops)
	default:
		return (*OutOps)(&C.rte_port_ring_writer_ops)
	}
}

// Arg implements ConfigOut interface.
func (wr *RingOut) Arg(mem common.Allocator) *OutArg {
	// NOTE: struct rte_port_ring_writer_params is a subset of struct
	// rte_port_ring_writer_nodrop_params, so we may simply use the
	// latter for it would fit regardless of NoDrop flag.
	var rc *C.struct_rte_port_ring_writer_nodrop_params
	common.MallocT(mem, &rc)
	rc.ring = (*C.struct_rte_ring)(unsafe.Pointer(wr.Ring))
	rc.tx_burst_sz = C.uint32_t(wr.TxBurstSize)
	rc.n_retries = C.uint32_t(wr.Retries)
	return (*OutArg)(unsafe.Pointer(rc))
}
