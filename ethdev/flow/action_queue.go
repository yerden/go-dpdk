package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

var _ Action = (*ActionQueue)(nil)

// ActionQueue implements Action which assigns packets to a given
// queue index.
type ActionQueue struct {
	Index uint16
}

// Transform implements Action interface.
func (action *ActionQueue) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	s := &C.struct_rte_flow_action_queue{
		index: C.ushort(action.Index),
	}
	return common.TransformPOD(alloc, s)
}

// ActionType implements Action interface.
func (action *ActionQueue) ActionType() ActionType {
	return ActionTypeQueue
}
