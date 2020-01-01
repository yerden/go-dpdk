package port

/*
#include <rte_config.h>
#include <rte_port.h>
#include <rte_port_ethdev.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// compile time checks
var _ = []ConfigIn{
	&EthdevIn{},
}

var _ = []ConfigOut{
	&EthdevOut{},
}

// EthdevIn is an input port built on top of pre-initialized NIC
// RX queue.
type EthdevIn struct {
	// Configured Ethernet port and RX queue ID.
	PortID, QueueID uint16
}

// Ops implements ConfigIn interface.
func (rd *EthdevIn) Ops() *InOps {
	return (*InOps)(&C.rte_port_ethdev_reader_ops)
}

// Arg implements ConfigIn interface.
func (rd *EthdevIn) Arg(mem common.Allocator) *InArg {
	var rc *C.struct_rte_port_ethdev_reader_params
	common.MallocT(mem, &rc)
	rc.port_id = C.uint16_t(rd.PortID)
	rc.queue_id = C.uint16_t(rd.QueueID)
	return (*InArg)(unsafe.Pointer(rc))
}

// EthdevOut is an output port built on top of pre-initialized NIC
// TX queue.
type EthdevOut struct {
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

// Ops implements ConfigOut interface.
func (wr *EthdevOut) Ops() *OutOps {
	if !wr.NoDrop {
		return (*OutOps)(&C.rte_port_ethdev_writer_ops)
	}
	return (*OutOps)(&C.rte_port_ethdev_writer_nodrop_ops)
}

// Arg implements ConfigOut interface.
func (wr *EthdevOut) Arg(mem common.Allocator) *OutArg {
	// NOTE: struct rte_port_ethdev_writer_params is a subset of struct
	// rte_port_ethdev_writer_nodrop_params, so we may simply use the latter
	// for it would fit regardless of NoDrop flag.
	var rc *C.struct_rte_port_ethdev_writer_nodrop_params
	common.MallocT(mem, &rc)
	rc.port_id = C.uint16_t(wr.PortID)
	rc.queue_id = C.uint16_t(wr.QueueID)
	rc.tx_burst_sz = C.uint32_t(wr.TxBurstSize)
	rc.n_retries = C.uint32_t(wr.Retries)
	return (*OutArg)(unsafe.Pointer(rc))
}
