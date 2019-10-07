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
var _ = []ReaderParams{
	&RingReader{},
}

var _ = []WriterParams{
	&RingWriter{},
}

// RingReader is an input port built on top of pre-initialized single
// consumer ring.
type RingReader struct {
	// Underlying ring
	*ring.Ring

	// Set if specified ring is multi consumer.
	Multi bool
}

// ReaderOps implements ReaderParams interface.
func (rd *RingReader) ReaderOps() (ops *ReaderOps, arg unsafe.Pointer) {
	if !rd.Multi {
		ops = (*ReaderOps)(&C.rte_port_ring_reader_ops)
	} else {
		ops = (*ReaderOps)(&C.rte_port_ring_multi_reader_ops)
	}
	rc := &C.struct_rte_port_ring_reader_params{
		ring: (*C.struct_rte_ring)(unsafe.Pointer(rd.Ring)),
	}
	return ops, unsafe.Pointer(rc)
}

// RingWriter is an output port built on top of pre-initialized single
// producer ring.
type RingWriter struct {
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

// WriterOps implements WriterParams interface.
func (wr *RingWriter) WriterOps() (ops *WriterOps, arg unsafe.Pointer) {
	switch {
	case wr.Multi && wr.NoDrop:
		ops = (*WriterOps)(&C.rte_port_ring_multi_writer_nodrop_ops)
	case wr.Multi:
		ops = (*WriterOps)(&C.rte_port_ring_multi_writer_ops)
	case wr.NoDrop:
		ops = (*WriterOps)(&C.rte_port_ring_writer_nodrop_ops)
	default:
		ops = (*WriterOps)(&C.rte_port_ring_writer_ops)
	}
	// NOTE: struct rte_port_ring_writer_params is a subset of struct
	// rte_port_ring_writer_nodrop_params, so we may simply use the
	// latter for it would fit regardless of NoDrop flag.
	arg = unsafe.Pointer(&C.struct_rte_port_ring_writer_nodrop_params{
		ring:        (*C.struct_rte_ring)(unsafe.Pointer(wr.Ring)),
		tx_burst_sz: C.uint32_t(wr.TxBurstSize),
		n_retries:   C.uint32_t(wr.Retries),
	})
	return
}
