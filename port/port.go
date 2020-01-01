/*
Package port wraps RTE port library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package port

import (
	"github.com/yerden/go-dpdk/common"
)

type opaqueData [0]byte

// In is an input port instance.
type In opaqueData

// Out is an output port instance.
type Out opaqueData

// InArg is an opaque argument to be used at input port creation.
type InArg opaqueData

// OutArg is an opaque argument to be used at output port creation.
type OutArg opaqueData

// ConfigIn implements reader port capability which allows to read
// packets from it.
type ConfigIn interface {
	Ops() *InOps
	Arg(common.Allocator) *InArg
}

// ConfigOut implements writer port capability which allows to
// write packets to it.
type ConfigOut interface {
	Ops() *OutOps
	Arg(common.Allocator) *OutArg
}

// XXX: we need to wrap calls which are not performance bottlenecks.

func err(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

// CreateIn creates input ops and input port instance.
func CreateIn(c ConfigIn, socket int) (*InOps, *In) {
	mem := common.NewAllocatorSession(&common.StdAlloc{})
	defer mem.Flush()
	ops := c.Ops()
	return ops, ops.Create(socket, c.Arg(mem))
}

// CreateOut creates output ops and output port instance.
func CreateOut(c ConfigOut, socket int) (*OutOps, *Out) {
	mem := common.NewAllocatorSession(&common.StdAlloc{})
	defer mem.Flush()
	ops := c.Ops()
	return ops, ops.Create(socket, c.Arg(mem))
}
