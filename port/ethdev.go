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
var _ = []RxFactory{
	&EthdevRx{},
}

var _ = []TxFactory{
	&EthdevTx{},
}

// EthdevRx is an input port built on top of pre-initialized NIC
// RX queue.
type EthdevRx struct {
	// Configured Ethernet port and RX queue ID.
	PortID, QueueID uint16
}

// CreateRx implements RxFactory interface.
func (rd *EthdevRx) CreateRx(socket int) (*Rx, error) {
	rx := &Rx{
		ops: &C.rte_port_ethdev_reader_ops,
	}

	// port
	params := &C.struct_rte_port_ethdev_reader_params{
		port_id:  C.uint16_t(rd.PortID),
		queue_id: C.uint16_t(rd.QueueID),
	}

	return rx, rx.doCreate(socket, unsafe.Pointer(params))
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

// CreateTx implements TxFactory interface.
func (wr *EthdevTx) CreateTx(socket int) (*Tx, error) {
	tx := &Tx{}

	// port
	var params unsafe.Pointer

	if wr.NoDrop {
		tx.ops = &C.rte_port_ethdev_writer_nodrop_ops
		params = unsafe.Pointer(&C.struct_rte_port_ethdev_writer_nodrop_params{
			port_id:     C.uint16_t(wr.PortID),
			queue_id:    C.uint16_t(wr.QueueID),
			tx_burst_sz: C.uint32_t(wr.TxBurstSize),
			n_retries:   C.uint32_t(wr.Retries),
		})
	} else {
		tx.ops = &C.rte_port_ethdev_writer_ops
		params = unsafe.Pointer(&C.struct_rte_port_ethdev_writer_params{
			port_id:     C.uint16_t(wr.PortID),
			queue_id:    C.uint16_t(wr.QueueID),
			tx_burst_sz: C.uint32_t(wr.TxBurstSize),
		})
	}

	return tx, tx.doCreate(socket, params)
}
