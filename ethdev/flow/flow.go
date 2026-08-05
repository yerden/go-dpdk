package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"

import (
	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/ethdev"
)

// static test that ItemTypeEnd is 0.
const _ uintptr = -uintptr(ItemTypeEnd)

// static test that ActionTypeEnd is 0.
const _ uintptr = -uintptr(ActionTypeEnd)

// Flow is the opaque flow handle.
type Flow C.struct_rte_flow

type cArgs struct {
	pid  C.ushort
	attr C.struct_rte_flow_attr
	pat  *C.struct_rte_flow_item
	act  *C.struct_rte_flow_action
	e    *C.struct_rte_flow_error
}

func doFancy(port ethdev.Port, attr *Attr, pattern []Item, actions []Action, flowErr *Error, fn func(*cArgs)) {
	alloc := common.NewAllocatorSession(&common.StdAlloc{})
	defer alloc.Flush()

	// patterns
	var pat []C.struct_rte_flow_item
	for i := range pattern {
		p := &pattern[i]
		cPat := C.struct_rte_flow_item{}
		cPat._type = uint32(p.Spec.ItemType())
		cPat.spec, _ = pattern[i].Spec.Transform(alloc)
		cPat.last, _ = pattern[i].Last.Transform(alloc)
		cPat.mask, _ = pattern[i].Mask.Transform(alloc)
		pat = append(pat, cPat)
	}

	// patterns finalizer
	pat = append(pat, C.struct_rte_flow_item{})

	// actions
	var act []C.struct_rte_flow_action
	for _, p := range actions {
		cAction := C.struct_rte_flow_action{}
		cAction._type = uint32(p.ActionType())
		cAction.conf, _ = p.Transform(alloc)
		act = append(act, cAction)
	}

	// actions finalizer
	act = append(act, C.struct_rte_flow_action{})

	args := &cArgs{
		pid:  C.ushort(port),
		attr: attr.cvtAttr(),
		pat:  &pat[0],
		act:  &act[0],
		e:    (*C.struct_rte_flow_error)(flowErr),
	}

	fn(args)
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
func Create(port ethdev.Port, attr *Attr, pattern []Item, actions []Action, flowErr *Error) (f *Flow, err error) {
	doFancy(port, attr, pattern, actions, flowErr, func(args *cArgs) {
		if p := C.rte_flow_create(args.pid, &args.attr, args.pat, args.act, args.e); p != nil {
			f = (*Flow)(p)
		} else {
			err = common.RteErrno()
		}
	})
	return
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
func Validate(port ethdev.Port, attr *Attr, pattern []Item, actions []Action, flowErr *Error) (err error) {
	doFancy(port, attr, pattern, actions, flowErr, func(args *cArgs) {
		rc := C.rte_flow_validate(args.pid, &args.attr, args.pat, args.act, args.e)
		err = common.IntErr(int64(rc))
	})
	return
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
