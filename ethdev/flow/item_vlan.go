package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_ether.h>
#include <rte_flow.h>

enum {
	ITEM_VLAN_OFF_HDR = offsetof(struct rte_flow_item_vlan, hdr),
};

static const void *get_item_vlan_mask() {
	return &rte_flow_item_vlan_mask;
}

static void set_has_more_vlan(struct rte_flow_item_vlan *item, uint32_t d) {
	item->has_more_vlan = d;
}

*/
import "C"
import (
	"runtime"
	"unsafe"
)

var _ ItemStruct = (*ItemVlan)(nil)

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
	cPointer
	HasMoreVlan bool
	TCI         uint16
	InnerType   uint16
}

// Reload implements ItemStruct interface.
func (item *ItemVlan) Reload() {
	cptr := (*C.struct_rte_flow_item_vlan)(item.createOrRet(C.sizeof_struct_rte_flow_item_vlan))

	if item.HasMoreVlan {
		C.set_has_more_vlan(cptr, 1)
	} else {
		C.set_has_more_vlan(cptr, 0)
	}

	hdr := (*C.struct_rte_vlan_hdr)(off(unsafe.Pointer(cptr), C.ITEM_VLAN_OFF_HDR))
	beU16(item.TCI, unsafe.Pointer(&hdr.vlan_tci))
	beU16(item.InnerType, unsafe.Pointer(&hdr.eth_proto))

	runtime.SetFinalizer(item, (*ItemVlan).free)
}

// Type implements ItemStruct interface.
func (item *ItemVlan) Type() ItemType {
	return ItemTypeVlan
}

// Mask implements ItemStruct interface.
func (item *ItemVlan) Mask() unsafe.Pointer {
	return C.get_item_vlan_mask()
}
