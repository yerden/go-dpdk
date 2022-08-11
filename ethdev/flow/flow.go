package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"

import (
	"runtime"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/ethdev"
)

// static test that ItemTypeEnd is 0.
const _ uintptr = -uintptr(ItemTypeEnd)

// static test that ActionTypeEnd is 0.
const _ uintptr = -uintptr(ActionTypeEnd)

// Flow is the opaque flow handle.
type Flow C.struct_rte_flow

// allocate c-style list of rte_flow_item's.
func cPattern(pattern []Item) []C.struct_rte_flow_item {
	pat := make([]C.struct_rte_flow_item, len(pattern)+1)

	for i := range pattern {
		typ := pattern[i].Spec.Type()
		pat[i]._type = uint32(typ)
		pattern[i].Spec.Reload()
		pat[i].spec = pattern[i].Spec.Pointer()
		if item := pattern[i].Mask; item != nil {
			item.Reload()
			pat[i].mask = item.Pointer()
		}
		if item := pattern[i].Last; item != nil {
			item.Reload()
			pat[i].last = item.Pointer()
		}
	}

	return pat
}

// allocate c-style list of rte_flow_action's.
func cActions(actions []Action) []C.struct_rte_flow_action {
	act := make([]C.struct_rte_flow_action, len(actions)+1)

	for i := range actions {
		typ := actions[i].Type()
		act[i]._type = uint32(typ)
		actions[i].Reload()
		act[i].conf = actions[i].Pointer()
	}

	return act
}

// Create a flow rule on a given port.
//
// port is port identifier of Ethernet device. attr is the Flow
// rule attributes.  pattern is a pattern specification (list
// terminated by the END pattern item).  actions is Associated
// actions (list terminated by the END action).  error is to
// Perform verbose error reporting if not NULL. PMDs initialize this
// structure in case of error only.
//
// Returns a valid handle in case of success, NULL otherwise and
// rte_errno is set to the positive version of one of the error codes
// defined for rte_flow_validate().
func Create(port ethdev.Port, attr *Attr, pattern []Item, actions []Action, flowErr *Error) (*Flow, error) {
	pat := cPattern(pattern)
	act := cActions(actions)
	cAttr := attr.cvtAttr()
	f := C.rte_flow_create(C.ushort(port), &cAttr, &pat[0], &act[0], (*C.struct_rte_flow_error)(flowErr))
	runtime.KeepAlive(pattern)
	runtime.KeepAlive(actions)
	if f == nil {
		return nil, common.RteErrno()
	}

	return (*Flow)(f), nil
}

// Validate checks whether a flow rule can be created on a given port.
//
// The flow rule is validated for correctness and whether it could be
// accepted by the device given sufficient resources. The rule is
// checked against the current device mode and queue configuration.
// The flow rule may also optionally be validated against existing
// flow rules and device resources. This function has no effect on the
// target device.
//
// The returned value is guaranteed to remain valid only as long as no
// successful calls to rte_flow_create() or rte_flow_destroy() are
// made in the meantime and no device parameter affecting flow rules
// in any way are modified, due to possible collisions or resource
// limitations (although in such cases EINVAL should not be returned).
func Validate(port ethdev.Port, attr *Attr, pattern []Item, actions []Action, flowErr *Error) error {
	pat := cPattern(pattern)
	act := cActions(actions)
	cAttr := attr.cvtAttr()
	ret := C.rte_flow_validate(C.ushort(port), &cAttr, &pat[0], &act[0], (*C.struct_rte_flow_error)(flowErr))
	runtime.KeepAlive(pattern)
	runtime.KeepAlive(actions)
	return common.IntToErr(ret)
}

// Destroy a flow rule on a given port.
//
// Failure to destroy a flow rule handle may occur when other flow
// rules depend on it, and destroying it would result in an
// inconsistent state.
//
// This function is only guaranteed to succeed if handles are
// destroyed in reverse order of their creation.
func Destroy(port ethdev.Port, flow *Flow, flowErr *Error) error {
	return common.IntToErr(C.rte_flow_destroy(C.ushort(port), (*C.struct_rte_flow)(flow),
		(*C.struct_rte_flow_error)(flowErr)))
}

// Flush destroys all flow rules associated with a port.
//
// In the unlikely event of failure, handles are still considered
// destroyed and no longer valid but the port must be assumed to be in
// an inconsistent state.
func Flush(port ethdev.Port, flowErr *Error) error {
	return common.IntToErr(C.rte_flow_flush(C.ushort(port), (*C.struct_rte_flow_error)(flowErr)))
}

// Isolate restricts ingress traffic to the defined flow rules.
//
// Isolated mode guarantees that all ingress traffic comes from
// defined flow rules only (current and future).
//
// Besides making ingress more deterministic, it allows PMDs to safely
// reuse resources otherwise assigned to handle the remaining traffic,
// such as global RSS configuration settings, VLAN filters, MAC
// address entries, legacy filter API rules and so on in order to
// expand the set of possible flow rule types.
//
// Calling this function as soon as possible after device
// initialization, ideally before the first call to
// rte_eth_dev_configure(), is recommended to avoid possible failures
// due to conflicting settings.
//
// Once effective, leaving isolated mode may not be possible depending
// on PMD implementation.
//
// port is the identifier of Ethernet device. Specify set to nonzero
// value to enter isolated mode, otherwise it marks the attempt to
// leave it. flowErr performs verbose error reporting if not NULL.
// PMDs initialize this structure in case of error only.
func Isolate(port ethdev.Port, set int, flowErr *Error) error {
	return common.IntToErr(C.rte_flow_isolate(C.ushort(port), C.int(set), (*C.struct_rte_flow_error)(flowErr)))
}
