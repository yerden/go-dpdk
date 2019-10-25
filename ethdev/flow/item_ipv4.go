package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>

static const struct rte_flow_item_ipv4 *get_item_ipv4_mask() {
	return &rte_flow_item_ipv4_mask;
}

*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/google/gopacket/layers"
)

// ItemIPv4 matches an IPv4 header.
//
// Note: IPv4 options are handled by dedicated pattern items.
type ItemIPv4 struct {
	cPointer

	Header layers.IPv4
}

var _ ItemStruct = (*ItemIPv4)(nil)

// Reload implements ItemStruct interface.
func (item *ItemIPv4) Reload() {
	cptr := (*C.struct_rte_flow_item_ipv4)(item.createOrRet(C.sizeof_struct_rte_flow_item_ipv4))
	cvtIPv4Header(&cptr.hdr, &item.Header)
	runtime.SetFinalizer(item, (*ItemIPv4).free)
}

func cvtIPv4Header(dst *C.struct_rte_ipv4_hdr, src *layers.IPv4) {
	dst.version_ihl = C.uint8_t(src.Version<<4 + src.IHL)
	dst.type_of_service = C.uint8_t(src.TOS)
	beU16(src.Length, unsafe.Pointer(&dst.total_length))
	beU16(src.Id, unsafe.Pointer(&dst.packet_id))
	beU16(src.FragOffset, unsafe.Pointer(&dst.fragment_offset))
	dst.time_to_live = C.uint8_t(src.TTL)
	dst.next_proto_id = C.uint8_t(src.Protocol)
	beU16(src.Checksum, unsafe.Pointer(&dst.hdr_checksum))

	if addr := src.SrcIP.To4(); addr != nil {
		dst.src_addr = *(*C.rte_be32_t)(unsafe.Pointer(&addr[0]))
	}

	if addr := src.DstIP.To4(); addr != nil {
		dst.dst_addr = *(*C.rte_be32_t)(unsafe.Pointer(&addr[0]))
	}
}

// Type implements ItemStruct interface.
func (item *ItemIPv4) Type() ItemType {
	return ItemTypeIpv4
}

// Mask implements ItemStruct interface.
func (item *ItemIPv4) Mask() unsafe.Pointer {
	return unsafe.Pointer(C.get_item_ipv4_mask())
}
