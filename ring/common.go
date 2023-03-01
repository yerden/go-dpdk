package ring

/*
#include <stdlib.h>
#include <rte_config.h>
#include <rte_ring.h>
#include "ring.h"
*/
import "C"

import (
	"unsafe"
)

func args(r *Ring, obj []unsafe.Pointer) (*C.struct_rte_ring,
	C.uintptr_t, C.uint) {
	return (*C.struct_rte_ring)(r), C.uintptr_t(uintptr(unsafe.Pointer(&obj[0]))), C.uint(len(obj))
}

func argsElem(r *Ring, objtable unsafe.Pointer, esize uintptr, n int) (*C.struct_rte_ring,
	C.uint, C.uintptr_t, C.uint) {
	return (*C.struct_rte_ring)(r), C.uint(esize), C.uintptr_t(uintptr(unsafe.Pointer(objtable))), C.uint(n)
}

func ret(out C.struct_compound_int) (rc, n uint32) {
	return uint32(out.rc), uint32(out.n)
}
