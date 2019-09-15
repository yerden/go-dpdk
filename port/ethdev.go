package port

/*
#include <rte_config.h>
#include <rte_port.h>
#include <rte_port_ethdev.h>
*/
import "C"

import (
	"unsafe"
)

// compile time checks
var _ = []ReaderParams{
	&EthdevReader{},
}

var _ = []WriterParams{
	&EthdevWriter{},
}

// EthdevReader is an input port built on top of pre-initialized NIC
// RX queue.
type EthdevReader struct {
	// Configured Ethernet port and RX queue ID.
	PortID, QueueID uint16
}

// ReaderOps implements ReaderParams interface.
func (rd *EthdevReader) ReaderOps() *ReaderOps {
	return (*ReaderOps)(&C.rte_port_ethdev_reader_ops)
}

// NewArg implements ReaderParams interface.
func (rd *EthdevReader) NewArg() unsafe.Pointer {
	rc := &C.struct_rte_port_ethdev_reader_params{
		port_id:  C.uint16_t(rd.PortID),
		queue_id: C.uint16_t(rd.QueueID),
	}
	return unsafe.Pointer(rc)
}

// EthdevWriter is an output port built on top of pre-initialized NIC
// TX queue.
type EthdevWriter struct {
	// Configured Ethernet port and TX queue ID.
	PortID, QueueID uint16

	// Recommended burst size for NIC TX queue.
	TxBurstSize uint32

	// If NoDrop set writer makes Retries attempts to write packets to
	// NIC TX queue.
	NoDrop bool

	// If NoDrop set and Retries is 0, number of retries is unlimited.
	Retries uint32
}

// WriterOps implements WriterParams interface.
func (wr *EthdevWriter) WriterOps() *WriterOps {
	if !wr.NoDrop {
		return (*WriterOps)(&C.rte_port_ethdev_writer_ops)
	}
	return (*WriterOps)(&C.rte_port_ethdev_writer_nodrop_ops)
}

// NewArg implements WriterParams interface.
func (wr *EthdevWriter) NewArg() unsafe.Pointer {
	// NOTE: struct rte_port_ethdev_writer_params is a subset of struct
	// rte_port_ethdev_writer_nodrop_params, so we may simply use the latter
	// for it would fit regardless of NoDrop flag.
	return unsafe.Pointer(&C.struct_rte_port_ethdev_writer_nodrop_params{
		port_id:     C.uint16_t(wr.PortID),
		queue_id:    C.uint16_t(wr.QueueID),
		tx_burst_sz: C.uint32_t(wr.TxBurstSize),
		n_retries:   C.uint32_t(wr.Retries),
	})
}
