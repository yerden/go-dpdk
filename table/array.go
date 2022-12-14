package table

/*
#include <rte_table_array.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// ArrayParams is the Array table parameters.
type ArrayParams struct {
	// Number of array entries. Has to be a power of two.
	Entries uint32

	// Byte offset within input packet meta-data where lookup key
	// (i.e. the array entry index) is located.
	Offset uint32
}

// Transform implements common.Transformer interface.
func (p *ArrayParams) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return common.TransformPOD(alloc, &C.struct_rte_table_array_params{
		n_entries: C.uint32_t(p.Entries),
		offset:    C.uint32_t(p.Offset),
	})
}

// Ops implements Params interface.
func (p *ArrayParams) Ops() *Ops {
	return (*Ops)(&C.rte_table_array_ops)
}

// ArrayKey is the key for adding/deleting key from table.
type ArrayKey struct {
	Pos uint32
}
