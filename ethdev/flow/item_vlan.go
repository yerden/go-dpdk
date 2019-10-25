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

type ItemVlan struct {
	cPointer
	HasMoreVlan bool
	TCI         uint16
	InnerType   uint16
}

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

func (item *ItemVlan) Type() ItemType {
	return ItemTypeVlan
}

func (item *ItemVlan) Mask() unsafe.Pointer {
	return C.get_item_vlan_mask()
}
