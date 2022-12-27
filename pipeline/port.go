package pipeline

/*
#include <rte_pipeline.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/port"
)

// PortIn represents an input port created in a pipeline.
type PortIn C.uint32_t

// PortOut represents an output port created in a pipeline.
type PortOut C.uint32_t

// PortInAction is the pipeline input port action handler.
//
// The action handler can decide to drop packets by resetting the
// associated packet bit in the pkts_mask parameter. In this case, the
// action handler is required not to free the packet buffer, which
// will be freed eventually by the pipeline.
type PortInAction C.rte_pipeline_port_in_action_handler

// PortOutAction is the pipeline output port action handler.
//
// The action handler can decide to drop packets by resetting the
// associated packet bit in the pkts_mask parameter. In this case, the
// action handler is required not to free the packet buffer, which
// will be freed eventually by the pipeline.
type PortOutAction C.rte_pipeline_port_out_action_handler

// PortInParams configures input port creation for pipeline.
type PortInParams struct {
	// Port configuration. See port package.
	Params port.InParams

	// Input port action. Implemented in C.
	Action PortInAction
	// Input port action argument.
	ActionArg unsafe.Pointer

	// Amount of packet burst.
	BurstSize int
}

// PortOutParams configures output port creation for pipeline.
type PortOutParams struct {
	// Port configuration. See port package.
	Params port.OutParams

	// Output port action. Implemented in C.
	Action PortOutAction
	// Output port action argument.
	ActionArg unsafe.Pointer
}

// PortInCreate creates new input port in the pipeline with specified
// configuration. Returns id of the created port or an error.
func (pl *Pipeline) PortInCreate(p *PortInParams) (id PortIn, err error) {
	params := &C.struct_rte_pipeline_port_in_params{}
	params.ops = (*C.struct_rte_port_in_ops)(unsafe.Pointer(p.Params.InOps()))

	{
		x, dtor := p.Params.Transform(alloc)
		defer dtor(x)
		params.arg_create = x
	}

	params.f_action = (C.rte_pipeline_port_in_action_handler)(p.Action)
	params.arg_ah = p.ActionArg
	params.burst_size = C.uint(p.BurstSize)

	rc := C.rte_pipeline_port_in_create(
		(*C.struct_rte_pipeline)(pl),
		params,
		(*C.uint32_t)(&id))

	return id, common.IntErr(int64(rc))
}

// PortOutCreate creates new output port in the pipeline with specified
// configuration. Returns id of the created port or an error.
func (pl *Pipeline) PortOutCreate(p *PortOutParams) (id PortOut, err error) {
	params := &C.struct_rte_pipeline_port_out_params{}
	params.ops = (*C.struct_rte_port_out_ops)(unsafe.Pointer(p.Params.OutOps()))

	{
		x, dtor := p.Params.Transform(alloc)
		defer dtor(x)
		params.arg_create = x
	}

	params.f_action = (C.rte_pipeline_port_out_action_handler)(p.Action)
	params.arg_ah = p.ActionArg

	rc := C.rte_pipeline_port_out_create(
		(*C.struct_rte_pipeline)(pl),
		params,
		(*C.uint32_t)(&id))

	return id, common.IntErr(int64(rc))
}

// ConnectToTable connects specified portID to created tableID.
func (pl *Pipeline) ConnectToTable(port PortIn, table Table) error {
	return common.IntErr(int64(C.rte_pipeline_port_in_connect_to_table(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(port),
		C.uint32_t(table))))
}

// Disable input port in the pipeline.
func (pl *Pipeline) Disable(port PortIn) error {
	return common.IntErr(int64(C.rte_pipeline_port_in_disable(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(port))))
}

// Enable input port in the pipeline.
func (pl *Pipeline) Enable(port PortIn) error {
	return common.IntErr(int64(C.rte_pipeline_port_in_enable(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(port))))
}
