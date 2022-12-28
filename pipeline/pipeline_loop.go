package pipeline

/*
#include "pipeline_loop.h"

static int
simple_control(void *params, struct rte_pipeline *p)
{
	volatile uint8_t *stop = (uint8_t *)params;
	return stop[0];
}

pipeline_op_ctrl go_simple_control = simple_control;

*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// OpsCtrlF returns 0 if the pipeline continues, >0 if the pipeline
// should be stopped and <0 in case of error.
//
// The implementation should satisfy the signature:
//   int (*pipeline_op_ctrl)(void *params, struct rte_pipeline *p);
type OpsCtrlF C.pipeline_op_ctrl

// Ops is the control operations of a pipeline tight loop.
type Ops struct {
	Ctrl OpsCtrlF
}

// Controller controls execution of pipeline loop.
type Controller interface {
	// Transformer is implemented to reflect the configuration into C
	// memory.
	common.Transformer

	// Ops returns the C implementation of hooks for controlling
	// pipeline execution.
	Ops() *Ops
}

// RunLoop runs pipeline continuously, flushing it every 'flush'
// iterations. RunLoop returns if value referenced by 'stop' is set to
// non-zero value.
//
// flush must be power of 2.
func (pl *Pipeline) RunLoop(flush uint32, c Controller) error {
	if flush&(flush-1) != 0 {
		panic("flush should be power of 2")
	}

	params := &C.struct_lcore_arg{}

	{
		arg, dtor := c.Transform(alloc)
		defer dtor(arg)
		params.ops_arg = arg
	}

	ops := c.Ops()

	params.ops = C.struct_pipeline_ops{
		f_ctrl: (C.pipeline_op_ctrl)(ops.Ctrl),
	}

	params.p = (*C.struct_rte_pipeline)(pl)
	params.flush = C.uint32_t(flush)

	return common.IntErr(int64(C.run_pipeline_loop(unsafe.Pointer(params))))
}

// SimpleController implements Controller for pipeline control.
//
// It allocates a single byte which serves as a stop flag for pipeline
// loop. Empty SimpleController is not usable, please use
// NewSimpleController to create SimpleController instance.
type SimpleController struct {
	stop unsafe.Pointer
}

// NewSimpleController allocates new SimpleController.
func NewSimpleController() *SimpleController {
	c := &SimpleController{}
	c.stop = alloc.Malloc(1)
	*(*uint8)(c.stop) = 0

	runtime.SetFinalizer(c, func(c *SimpleController) {
		alloc.Free(c.stop)
	})

	return c
}

// Transform implements common.Transformer interface.
func (c *SimpleController) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	return c.stop, func(unsafe.Pointer) {}
}

// Ops implements Controller interface.
func (c *SimpleController) Ops() *Ops {
	return &Ops{
		Ctrl: OpsCtrlF(C.go_simple_control),
	}
}

// Stop stops the execution of a pipeline that watches this
// controller. It may be used by multiple pipelines.
func (c *SimpleController) Stop() {
	*(*uint8)(c.stop) = 1
}
