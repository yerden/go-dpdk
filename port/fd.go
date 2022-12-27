package port

/*
#include <rte_config.h>
#include <rte_errno.h>

#include <rte_port_fd.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
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

// InOps implements InParams interface.
func (rd *FdRx) InOps() *InOps {
	return (*InOps)(&C.rte_port_fd_reader_ops)
}

// Transform implements common.Transformer interface.
func (rd *FdRx) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return common.TransformPOD(alloc, &C.struct_rte_port_fd_reader_params{
		fd:      C.int(rd.Fd),
		mtu:     C.uint32_t(rd.MTU),
		mempool: (*C.struct_rte_mempool)(unsafe.Pointer(rd.Mempool)),
	})
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

// OutOps implements OutParams interface.
func (wr *FdTx) OutOps() *OutOps {
	if wr.NoDrop {
		return (*OutOps)(&C.rte_port_fd_writer_nodrop_ops)
	}

	return (*OutOps)(&C.rte_port_fd_writer_ops)
}

// Transform implements common.Transformer interface.
func (wr *FdTx) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return common.TransformPOD(alloc, &C.struct_rte_port_fd_writer_nodrop_params{
		fd:        C.int(wr.Fd),
		n_retries: C.uint32_t(wr.Retries),
	})
}
