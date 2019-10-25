package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"
import (
	"runtime"
)

var _ Action = (*ActionQueue)(nil)

// ActionQueue implements Action which assigns packets to a given
// queue index.
type ActionQueue struct {
	cPointer
	Index uint16
}

// Reload implements Action interface.
func (action *ActionQueue) Reload() {
	cptr := (*C.struct_rte_flow_action_queue)(action.createOrRet(C.sizeof_struct_rte_flow_action_queue))

	cptr.index = C.uint16_t(action.Index)
	runtime.SetFinalizer(action, (*ActionQueue).free)
}

// Type implements Action interface.
func (action *ActionQueue) Type() ActionType {
	return ActionTypeQueue
}
