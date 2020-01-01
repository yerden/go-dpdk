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

// FdIn input port built on top of valid non-blocking file
// descriptor.
type FdIn struct {
	// Pre-initialized buffer pool.
	*mempool.Mempool

	// File descriptor.
	Fd uintptr

	// Maximum Transfer Unit (MTU)
	MTU uint32
}

// Ops implements ConfigIn interface.
func (rd *FdIn) Ops() *InOps {
	return (*InOps)(&C.rte_port_fd_reader_ops)
}

// Arg implements ConfigIn interface.
func (rd *FdIn) Arg(mem common.Allocator) *InArg {
	var rc *C.struct_rte_port_fd_reader_params
	common.MallocT(mem, &rc)
	rc.fd = C.int(rd.Fd)
	rc.mtu = C.uint32_t(rd.MTU)
	rc.mempool = (*C.struct_rte_mempool)(unsafe.Pointer(rd.Mempool))
	return (*InArg)(unsafe.Pointer(rc))
}

// FdOut is an output port built on top of valid non-blocking file
// descriptor.
type FdOut struct {
	// File descriptor.
	Fd uintptr

	// If NoDrop set writer makes Retries attempts to write packets to
	// ring.
	NoDrop bool

	// If NoDrop set and Retries is 0, number of retries is unlimited.
	Retries uint32
}

// Ops implements ConfigOut interface.
func (wr *FdOut) Ops() *OutOps {
	if !wr.NoDrop {
		return (*OutOps)(&C.rte_port_fd_writer_ops)
	}
	return (*OutOps)(&C.rte_port_fd_writer_nodrop_ops)
}

// Arg implements ConfigOut interface.
func (wr *FdOut) Arg(mem common.Allocator) *OutArg {
	var rc *C.struct_rte_port_fd_writer_nodrop_params
	common.MallocT(mem, &rc)
	rc.fd = C.int(wr.Fd)
	rc.n_retries = C.uint32_t(wr.Retries)
	return (*OutArg)(unsafe.Pointer(rc))
}
