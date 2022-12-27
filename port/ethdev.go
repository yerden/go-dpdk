package port

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_port.h>
#include <rte_port_ethdev.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

var (
	_ InParams  = (*EthdevRx)(nil)
	_ OutParams = (*EthdevTx)(nil)
)

// EthdevRx is an input port built on top of pre-initialized NIC
// RX queue.
type EthdevRx struct {
	// Configured Ethernet port and RX queue ID.
	PortID, QueueID uint16
}

var _ InParams = (*EthdevRx)(nil)

// InOps implements InParams interface.
func (p *EthdevRx) InOps() *InOps {
	return (*InOps)(&C.rte_port_ethdev_reader_ops)
}

// Transform implements common.Transformer interface.
func (p *EthdevRx) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return common.TransformPOD(alloc, &C.struct_rte_port_ethdev_reader_params{
		port_id:  C.uint16_t(p.PortID),
		queue_id: C.uint16_t(p.QueueID),
	})
}

// EthdevTx is an output port built on top of pre-initialized NIC
// TX queue.
type EthdevTx struct {
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

// OutOps implements OutParams interface.
func (p *EthdevTx) OutOps() *OutOps {
	if !p.NoDrop {
		return (*OutOps)(&C.rte_port_ethdev_writer_ops)
	}
	return (*OutOps)(&C.rte_port_ethdev_writer_nodrop_ops)
}

// Transform implements common.Transformer interface.
func (p *EthdevTx) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return common.TransformPOD(alloc, &C.struct_rte_port_ethdev_writer_nodrop_params{
		port_id:     C.uint16_t(p.PortID),
		queue_id:    C.uint16_t(p.QueueID),
		tx_burst_sz: C.uint32_t(p.TxBurstSize),
		n_retries:   C.uint32_t(p.Retries),
	})
}
