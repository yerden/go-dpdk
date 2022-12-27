package table

/*
#include <rte_table_stub.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// StubParams is the empty parameters for Stub table.
type StubParams struct{}

// Transform implements common.Transformer interface.
func (p *StubParams) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return nil, alloc.Free
}

// Ops implements Params interface.
func (p *StubParams) Ops() *Ops {
	return (*Ops)(&C.rte_table_stub_ops)
}
