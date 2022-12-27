// Package pipeline wraps RTE Pipeline.
//
// This tool is part of the DPDK Packet Framework tool suite and
//provides a standard methodology (logically similar to OpenFlow) for
// rapid development of complex packet processing pipelines out of
//ports, tables and actions.
//
//Basic operation. A pipeline is constructed by connecting its input
// ports to its output ports through a chain of lookup tables. As
//result of lookup operation into the current table, one of the table
// entries (or the default table entry, in case of lookup miss) is
// identified to provide the actions to be executed on the current
// packet and the associated action meta-data. The behavior of user
// actions is defined through the configurable table action handler,
// while the reserved actions define the next hop for the current
// packet (either another table, an output port or packet drop) and
// are handled transparently by the framework.
//
// Initialization and run-time flows. Once all the pipeline elements
// (input ports, tables, output ports) have been created, input ports
// connected to tables, table action handlers configured, tables
// populated with the initial set of entries (actions and action
// meta-data) and input ports enabled, the pipeline runs
// automatically, pushing packets from input ports to tables and
// output ports. At each table, the identified user actions are being
// executed, resulting in action meta-data (stored in the table entry)
// and packet meta-data (stored with the packet descriptor) being
// updated. The pipeline tables can have further updates and input
// ports can be disabled or enabled later on as required.
//
// Multi-core scaling. Typically, each CPU core will run its own
// pipeline instance. Complex application-level pipelines can be
// implemented by interconnecting multiple CPU core-level pipelines in
// tree-like topologies, as the same port devices (e.g. SW rings) can
// serve as output ports for the pipeline running on CPU core A, as
// well as input ports for the pipeline running on CPU core B. This
// approach enables the application development using the pipeline
// (CPU cores connected serially), cluster/run-to-completion (CPU
// cores connected in parallel) or mixed (pipeline of CPU core
// clusters) programming models.
//
// Thread safety. It is possible to have multiple pipelines running on
// the same CPU core, but it is not allowed (for thread safety
// reasons) to have multiple CPU cores running the same pipeline
// instance.
package pipeline

/*
#cgo pkg-config: libdpdk
#include <rte_pipeline.h>
*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Params configures the pipeline.
type Params struct {
	// Name of the pipeline.
	Name string

	// SocketID for memory allocation.
	SocketID int

	// Offset within packet meta-data to port_id to be used by action
	// "Send packet to output port read from packet meta-data". Has to
	// be 4-byte aligned.
	OffsetPortID uint32
}

// Transform implements common.Transformer interface.
func (p *Params) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	var params *C.struct_rte_pipeline_params
	params = (*C.struct_rte_pipeline_params)(alloc.Malloc(unsafe.Sizeof(*params)))

	params.name = (*C.char)(common.CString(alloc, p.Name))
	params.socket_id = C.int(p.SocketID)
	params.offset_port_id = C.uint(p.OffsetPortID)

	return unsafe.Pointer(params), func(p unsafe.Pointer) {
		params := (*C.struct_rte_pipeline_params)(p)
		alloc.Free(unsafe.Pointer(params.name))
		alloc.Free(p)
	}
}

// Pipeline object.
type Pipeline C.struct_rte_pipeline

var alloc = &common.StdAlloc{}

// Create pipeline with specified configuration.
//
// Returns nil in case of an error.
func Create(p *Params) *Pipeline {
	arg, dtor := p.Transform(alloc)
	defer dtor(arg)
	return (*Pipeline)(C.rte_pipeline_create((*C.struct_rte_pipeline_params)(arg)))
}

// Free destroys the pipeline.
func (pl *Pipeline) Free() error {
	return common.IntErr(int64(C.rte_pipeline_free((*C.struct_rte_pipeline)(pl))))
}

// Check the consistency of the pipeline.
func (pl *Pipeline) Check() error {
	return common.IntErr(int64(C.rte_pipeline_check((*C.struct_rte_pipeline)(pl))))
}

// Flush the pipeline.
func (pl *Pipeline) Flush() error {
	return common.IntErr(int64(C.rte_pipeline_flush((*C.struct_rte_pipeline)(pl))))
}

// Run the pipeline.
func (pl *Pipeline) Run() int {
	return int(C.rte_pipeline_run((*C.struct_rte_pipeline)(pl)))
}
