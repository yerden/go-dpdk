package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>

void set_flow_attr_ingress(struct rte_flow_attr *attr) {
	attr->ingress = 1;
}

void set_flow_attr_egress(struct rte_flow_attr *attr) {
	attr->egress = 1;
}

void set_flow_attr_transfer(struct rte_flow_attr *attr) {
	attr->transfer = 1;
}
*/
import "C"

// Attr is Flow rule attributes.
//
// Priorities are set on a per rule based within groups.
//
// Lower values denote higher priority, the highest priority for a
// flow rule is 0, so that a flow that matches for than one rule, the
// rule with the lowest priority value will always be matched.
//
// Although optional, applications are encouraged to group similar
// rules as much as possible to fully take advantage of hardware
// capabilities (e.g. optimized matching) and work around limitations
// (e.g. a single pattern type possibly allowed in a given group).
// Applications should be aware that groups are not linked by default,
// and that they must be explicitly linked by the application using
// the JUMP action.
//
// Priority levels are arbitrary and up to the application, they do
// not need to be contiguous nor start from 0, however the maximum
// number varies between devices and may be affected by existing flow
// rules.
//
// If a packet is matched by several rules of a given group for a
// given priority level, the outcome is undefined. It can take any
// path, may be duplicated or even cause unrecoverable errors.
//
// Note that support for more than a single group and priority level
// is not guaranteed.
//
// Flow rules can apply to inbound and/or outbound traffic
// (ingress/egress).
//
// Several pattern items and actions are valid and can be used in both
// directions. Those valid for only one direction are described as
// such.
//
// At least one direction must be specified.
//
// Specifying both directions at once for a given rule is not
// recommended but may be valid in a few cases (e.g. shared counter).
type Attr struct {
	// Priority group.
	Group uint32

	// Rule priority level within group.
	Priority uint32

	// Rule applies to ingress traffic.
	Ingress bool

	// Rule applies to egress traffic.
	Egress bool

	// Instead of simply matching the properties of traffic as it
	// would appear on a given DPDK port ID, enabling this attribute
	// transfers a flow rule to the lowest possible level of any
	// device endpoints found in the pattern.
	//
	// When supported, this effectively enables an application to
	// re-route traffic not necessarily intended for it (e.g. coming
	// from or addressed to different physical ports, VFs or
	// applications) at the device level.
	//
	// It complements the behavior of some pattern items such as
	// RTE_FLOW_ITEM_TYPE_PHY_PORT and is meaningless without them.
	//
	// When transferring flow rules, ingress and egress attributes
	// keep their original meaning, as if processing traffic emitted
	// or received by the application.
	Transfer bool
}

func (a *Attr) cvtAttr() (out C.struct_rte_flow_attr) {
	out.group = C.uint32_t(a.Group)
	out.priority = C.uint32_t(a.Priority)
	if a.Ingress {
		C.set_flow_attr_ingress(&out)
	}
	if a.Egress {
		C.set_flow_attr_egress(&out)
	}
	if a.Transfer {
		C.set_flow_attr_transfer(&out)
	}
	return
}
