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

static const struct rte_flow_item_eth *get_item_eth_mask() {
	return &rte_flow_item_eth_mask;
}

*/
import "C"
import (
	"net"
	"reflect"
	"runtime"
	"unsafe"
)

var _ ItemStruct = (*ItemEth)(nil)

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
	cPointer
	HasVlan   bool
	Src, Dst  net.HardwareAddr
	EtherType uint16
}

// Reload implements ItemStruct interface.
func (item *ItemEth) Reload() {
	cptr := (*C.struct_rte_flow_item_eth)(item.createOrRet(C.sizeof_struct_rte_flow_item_eth))

	var u uint32
	if item.HasVlan {
		u = 1
	}
	C.set_has_vlan(cptr, C.uint32_t(u))

	hdr := (*C.struct_rte_ether_hdr)(off(unsafe.Pointer(cptr), C.ITEM_ETH_OFF_HDR))

	if len(item.Src) > 0 {
		p := off(unsafe.Pointer(hdr), C.ETHER_HDR_OFF_SRC)
		setAddr((*C.struct_rte_ether_addr)(p), item.Src)
	}

	if len(item.Dst) > 0 {
		p := off(unsafe.Pointer(hdr), C.ETHER_HDR_OFF_DST)
		setAddr((*C.struct_rte_ether_addr)(p), item.Dst)
	}

	beU16(item.EtherType, unsafe.Pointer(&hdr.ether_type))

	runtime.SetFinalizer(item, (*ItemEth).free)
}

func setAddr(p *C.struct_rte_ether_addr, addr net.HardwareAddr) {
	var hwaddr []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&hwaddr))
	sh.Data = uintptr(unsafe.Pointer(&p.addr_bytes[0]))
	sh.Len = len(p.addr_bytes)
	sh.Cap = sh.Len
	copy(hwaddr, addr)
}

// Type implements ItemStruct interface.
func (item *ItemEth) Type() ItemType {
	return ItemTypeEth
}

// Mask implements ItemStruct interface.
func (item *ItemEth) Mask() unsafe.Pointer {
	return unsafe.Pointer(C.get_item_eth_mask())
}
