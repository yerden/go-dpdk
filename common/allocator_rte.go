package common

/*
#include <stdlib.h>
#include <string.h>

#include <rte_config.h>
#include <rte_malloc.h>
*/
import "C"

import "unsafe"

// RteAlloc implements allocator based on DPDK rte_malloc.h.
type RteAlloc struct {
	// Requested alignment.
	Align uint

	// Requested NUMA node. Set to SocketIDAny if meaningless.
	Socket int
}

var _ Allocator = (*RteAlloc)(nil)

// Malloc implements Allocator.
func (mem *RteAlloc) Malloc(size uintptr) unsafe.Pointer {
	return C.rte_malloc_socket(nil, C.size_t(size), C.uint(mem.Align), C.int(mem.Socket))
}

// Free implements Allocator.
func (mem *RteAlloc) Free(p unsafe.Pointer) {
	C.rte_free(p)
}

// Realloc implements Allocator.
//
// Note: rte_realloc() may not reside on the same NUMA node.
func (mem *RteAlloc) Realloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	return C.rte_realloc(p, C.size_t(size), C.uint(mem.Align))
}
