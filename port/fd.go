package port

/*
#include <rte_config.h>
#include <rte_errno.h>

#include <rte_port_fd.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/mempool"
)

// FdRx input port built on top of valid non-blocking file
// descriptor.
type FdRx struct {
	// Pre-initialized buffer pool.
	*mempool.Mempool

	// File descriptor.
	Fd uintptr

	// Maximum Transfer Unit (MTU)
	MTU uint32
}

// CreateRx implements RxFactory interface.
func (rd *FdRx) CreateRx(socket int) (*Rx, error) {
	rx := &Rx{
		ops: &C.rte_port_fd_reader_ops,
	}

	// port
	params := &C.struct_rte_port_fd_reader_params{
		fd:      C.int(rd.Fd),
		mtu:     C.uint32_t(rd.MTU),
		mempool: (*C.struct_rte_mempool)(unsafe.Pointer(rd.Mempool)),
	}

	return rx, rx.doCreate(socket, unsafe.Pointer(params))
}

// FdTx is an output port built on top of valid non-blocking file
// descriptor.
type FdTx struct {
	// File descriptor.
	Fd uintptr

	// If NoDrop set writer makes Retries attempts to write packets to
	// ring.
	NoDrop bool

	// If NoDrop set and Retries is 0, number of retries is unlimited.
	Retries uint32
}

// CreateTx implements TxFactory interface.
func (wr *FdTx) CreateTx(socket int) (*Tx, error) {
	tx := &Tx{}

	var params unsafe.Pointer
	if wr.NoDrop {
		tx.ops = &C.rte_port_fd_writer_nodrop_ops
		params = unsafe.Pointer(&C.struct_rte_port_fd_writer_nodrop_params{
			fd:        C.int(wr.Fd),
			n_retries: C.uint32_t(wr.Retries),
		})
	} else {
		tx.ops = &C.rte_port_fd_writer_ops
		params = unsafe.Pointer(&C.struct_rte_port_fd_writer_params{
			fd: C.int(wr.Fd),
		})
	}

	return tx, tx.doCreate(socket, params)
}
