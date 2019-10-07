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

// FdReader input port built on top of valid non-blocking file
// descriptor.
type FdReader struct {
	// Pre-initialized buffer pool.
	*mempool.Mempool

	// File descriptor.
	Fd uintptr

	// Maximum Transfer Unit (MTU)
	MTU uint32
}

// ReaderOps implements ReaderParams interface.
func (rd *FdReader) ReaderOps() (*ReaderOps, unsafe.Pointer) {
	ops := (*ReaderOps)(&C.rte_port_fd_reader_ops)
	rc := &C.struct_rte_port_fd_reader_params{
		fd:      C.int(rd.Fd),
		mtu:     C.uint32_t(rd.MTU),
		mempool: (*C.struct_rte_mempool)(unsafe.Pointer(rd.Mempool)),
	}
	return ops, unsafe.Pointer(rc)
}

// FdWriter is an output port built on top of valid non-blocking file
// descriptor.
type FdWriter struct {
	// File descriptor.
	Fd uintptr

	// If NoDrop set writer makes Retries attempts to write packets to
	// ring.
	NoDrop bool

	// If NoDrop set and Retries is 0, number of retries is unlimited.
	Retries uint32
}

// WriterOps implements WriterParams interface.
func (wr *FdWriter) WriterOps() (ops *WriterOps, arg unsafe.Pointer) {
	if !wr.NoDrop {
		ops = (*WriterOps)(&C.rte_port_fd_writer_ops)
	} else {
		ops = (*WriterOps)(&C.rte_port_fd_writer_nodrop_ops)
	}
	arg = unsafe.Pointer(&C.struct_rte_port_fd_writer_nodrop_params{
		fd:        C.int(wr.Fd),
		n_retries: C.uint32_t(wr.Retries),
	})
	return
}
