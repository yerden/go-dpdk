package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_ether.h>
#include <rte_flow.h>

enum {
	ITEM_VLAN_OFF_HDR = offsetof(struct rte_flow_item_vlan, hdr),
};

static void set_has_more_vlan(struct rte_flow_item_vlan *item, uint32_t d) {
	item->has_more_vlan = d;
}

*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

var _ ItemValue = (*ItemVlan)(nil)

// ItemVlan matches an 802.1Q/ad VLAN tag.
//
// The corresponding standard outer EtherType (TPID) values are
// RTE_ETHER_TYPE_VLAN or RTE_ETHER_TYPE_QINQ. It can be overridden by the
// preceding pattern item. If a VLAN item is present in the pattern, then only
// tagged packets will match the pattern. The field has_more_vlan can be used to
// match any type of tagged packets, instead of using the eth_proto field of hdr.
// If the eth_proto of hdr and has_more_vlan fields are not specified, then any
// tagged packets will match the pattern.
type ItemVlan struct {
	HasMoreVlan bool
	TCI         uint16
	InnerType   uint16
}

// Transform implements Action interface.
func (item *ItemVlan) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	cptr := &C.struct_rte_flow_item_vlan{}

	if item.HasMoreVlan {
		C.set_has_more_vlan(cptr, 1)
	} else {
		C.set_has_more_vlan(cptr, 0)
	}

	hdr := (*C.struct_rte_vlan_hdr)(unsafe.Add(unsafe.Pointer(cptr), C.ITEM_VLAN_OFF_HDR))
	hdr.vlan_tci = C.ushort(item.TCI)
	hdr.eth_proto = C.ushort(item.InnerType)
	return common.TransformPOD(alloc, cptr)
}

// ItemType implements ItemValue interface.
func (item *ItemVlan) ItemType() ItemType {
	return ItemTypeVlan
}

// DefaultMask implements ItemStruct interface.
func (item *ItemVlan) DefaultMask() unsafe.Pointer {
	return unsafe.Pointer(&C.rte_flow_item_vlan_mask)
}
