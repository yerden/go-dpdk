package common

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import "unsafe"

// StdAlloc wraps system malloc/free memory allocation.
type StdAlloc struct{}

var _ Allocator = (*StdAlloc)(nil)

// Malloc implements Allocator.
func (mem *StdAlloc) Malloc(size uintptr) unsafe.Pointer {
	return C.malloc(C.size_t(size))
}

// Free implements Allocator.
func (mem *StdAlloc) Free(p unsafe.Pointer) {
	C.free(p)
}

// Realloc implements Allocator.
func (mem *StdAlloc) Realloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	return C.realloc(p, C.size_t(size))
}
