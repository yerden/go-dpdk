package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

var _ Action = (*ActionRSS)(nil)

// ActionRSS implements Receive-Side Scaling feature.
//
// Similar to QUEUE, except RSS is additionally performed on packets
// to spread them among several queues according to the provided
// parameters.
//
// Unlike global RSS settings used by other DPDK APIs, unsetting the
// types field does not disable RSS in a flow rule. Doing so instead
// requests safe unspecified "best-effort" settings from the
// underlying PMD, which depending on the flow rule, may result in
// anything ranging from empty (single queue) to all-inclusive RSS.
//
// Note: RSS hash result is stored in the hash.rss mbuf field which
// overlaps hash.fdir.lo. Since the MARK action sets the hash.fdir.hi
// field only, both can be requested simultaneously.
type ActionRSS struct {
	Func   HashFunction
	Queues []uint16
	Key    []byte
	Level  uint32
	Types  uint64

	cptr *C.struct_rte_flow_action_rss
}

func (action *ActionRSS) free() {
	cptr := action.cptr
	C.free(unsafe.Pointer(cptr.key))
	C.free(unsafe.Pointer(cptr.queue))
	C.free(unsafe.Pointer(cptr))
}

// Reload implements Action interface.
func (action *ActionRSS) Reload() {
	// allocate if needed
	cptr := action.cptr
	if cptr == nil {
		cptr = (*C.struct_rte_flow_action_rss)(C.malloc(C.sizeof_struct_rte_flow_action_rss))
		*cptr = C.struct_rte_flow_action_rss{}
		action.cptr = cptr
	}

	// set queues
	if len(action.Queues) > 0 {
		sz := C.size_t(len(action.Queues)) * C.size_t(unsafe.Sizeof(action.Queues[0]))
		cQueues := C.malloc(sz)
		C.memcpy(cQueues, unsafe.Pointer(&action.Queues[0]), sz)
		C.free(unsafe.Pointer(cptr.queue))
		cptr.queue_num = C.uint32_t(len(action.Queues))
		cptr.queue = (*C.uint16_t)(cQueues)
	}

	// set key
	if len(action.Key) > 0 {
		sz := C.size_t(len(action.Key))
		cKey := C.malloc(sz)
		C.memcpy(cKey, unsafe.Pointer(&action.Key[0]), sz)
		C.free(unsafe.Pointer(cptr.key))
		cptr.key_len = C.uint32_t(len(action.Key))
		cptr.key = (*C.uint8_t)(cKey)
	}

	cptr.level = C.uint32_t(action.Level)
	cptr.types = C.uint64_t(action.Types)
	cptr._func = uint32(action.Func)

	runtime.SetFinalizer(action, nil)
	runtime.SetFinalizer(action, (*ActionRSS).free)
}

// Pointer implements Action interface.
func (action *ActionRSS) Pointer() unsafe.Pointer {
	return unsafe.Pointer(action.cptr)
}

// Type implements Action interface.
func (action *ActionRSS) Type() ActionType {
	return ActionTypeRss
}
