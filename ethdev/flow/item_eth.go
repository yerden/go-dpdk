package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_ether.h>
#include <rte_flow.h>
#include <rte_version.h>

enum {
	ITEM_ETH_OFF_HDR = offsetof(struct rte_flow_item_eth, hdr),
};

// Ethernet header fields renamed since commit: 04d43857ea3acbd4db4b28939dc2807932b85e72.
#if RTE_VERSION < RTE_VERSION_NUM(21, 11, 0, 0)
enum {
	ETHER_HDR_OFF_SRC = offsetof(struct rte_ether_hdr, s_addr),
	ETHER_HDR_OFF_DST = offsetof(struct rte_ether_hdr, d_addr),
};
#else
enum {
	ETHER_HDR_OFF_SRC = offsetof(struct rte_ether_hdr, src_addr),
	ETHER_HDR_OFF_DST = offsetof(struct rte_ether_hdr, dst_addr),
};
#endif

void set_has_vlan(struct rte_flow_item_eth *item, uint32_t b) {
	item->has_vlan = b;
}

*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

var _ ItemValue = (*ItemEth)(nil)

// ItemEth matches an Ethernet header.
//
// Inside hdr field, the sub-field ether_type stands either for
// EtherType or TPID, depending on whether the item is followed by a
// VLAN item or not. If two VLAN items follow, the sub-field refers to
// the outer one, which, in turn, contains the inner TPID in the
// similar header field. The innermost VLAN item contains a layer-3
// EtherType. All of that follows the order seen on the wire.
//
// If the field in question contains a TPID value, only tagged packets
// with the specified TPID will match the pattern. Alternatively, it's
// possible to match any type of tagged packets by means of the field
// has_vlan rather than use the EtherType/TPID field. Also, it's
// possible to leave the two fields unused. If this is the case, both
// tagged and untagged packets will match the pattern.
type ItemEth struct {
	HasVlan   bool
	Src, Dst  [6]byte
	EtherType uint16
}

func boolU32(b bool) (x C.uint32_t) {
	if b {
		x = 1
	}
	return
}

// Transform implements Action interface.
func (item *ItemEth) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	cptr := &C.struct_rte_flow_item_eth{}
	hdr := (*C.struct_rte_ether_addr)(unsafe.Add(unsafe.Pointer(cptr), C.ITEM_ETH_OFF_HDR))
	*(*[6]byte)(unsafe.Add(unsafe.Pointer(hdr), C.ETHER_HDR_OFF_SRC)) = item.Src
	*(*[6]byte)(unsafe.Add(unsafe.Pointer(hdr), C.ETHER_HDR_OFF_DST)) = item.Dst
	C.set_has_vlan(cptr, boolU32(item.HasVlan))
	return common.TransformPOD(alloc, cptr)
}

// ItemType implements ItemValue interface.
func (item *ItemEth) ItemType() ItemType {
	return ItemTypeEth
}

// DefaultMask implements ItemStruct interface.
func (item *ItemEth) DefaultMask() unsafe.Pointer {
	return unsafe.Pointer(&C.rte_flow_item_eth_mask)
}
