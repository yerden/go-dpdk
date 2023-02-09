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
}

// Transform implements Action interface.
func (action *ActionRSS) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	cptr := (*C.struct_rte_flow_action_rss)(alloc.Malloc(C.sizeof_struct_rte_flow_action_rss))
	*cptr = C.struct_rte_flow_action_rss{}

	// set queues
	if len(action.Queues) > 0 {
		var x *C.uint16_t
		common.CallocT(alloc, &x, len(action.Queues))
		queues := unsafe.Slice(x, len(action.Queues))
		for i := range queues {
			queues[i] = C.uint16_t(action.Queues[i])
		}
		cptr.queue_num = C.uint32_t(len(action.Queues))
	}

	// set key
	if len(action.Key) > 0 {
		cptr.key_len = C.uint32_t(len(action.Key))
		cptr.key = (*C.uchar)(common.CBytes(alloc, action.Key))
	}

	cptr.level = C.uint32_t(action.Level)
	cptr.types = C.uint64_t(action.Types)
	cptr._func = uint32(action.Func)

	return unsafe.Pointer(cptr), func(p unsafe.Pointer) {
		cptr = (*C.struct_rte_flow_action_rss)(p)
		alloc.Free(unsafe.Pointer(cptr.key))
		alloc.Free(unsafe.Pointer(cptr.queue))
		alloc.Free(p)
	}
}

// ActionType implements Action interface.
func (action *ActionRSS) ActionType() ActionType {
	return ActionTypeRss
}
