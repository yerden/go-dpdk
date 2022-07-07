package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>

enum {
	IPv4_HDR_OFF_DST_VERSION_IHL = offsetof(struct rte_ipv4_hdr, version_ihl),
};

static const struct rte_flow_item_ipv4 *get_item_ipv4_mask() {
	return &rte_flow_item_ipv4_mask;
}

*/
import "C"
import (
	"runtime"
	"unsafe"
)

// IPv4 represents a raw IPv4 address.
type IPv4 [4]byte

// IPv4Header is the IPv4 header raw format.
type IPv4Header struct {
	VersionIHL     uint8  /* Version and header length. */
	ToS            uint8  /* Type of service. */
	TotalLength    uint16 /* Length of packet. */
	ID             uint16 /* Packet ID. */
	FragmentOffset uint16 /* Fragmentation offset. */
	TTL            uint8  /* Time to live. */
	Proto          uint8  /* Protocol ID. */
	Checksum       uint16 /* Header checksum. */
	SrcAddr        IPv4   /* Source address. */
	DstAddr        IPv4   /* Destination address. */
}

// ItemIPv4 matches an IPv4 header.
//
// Note: IPv4 options are handled by dedicated pattern items.
type ItemIPv4 struct {
	cPointer

	Header IPv4Header
}

var _ ItemStruct = (*ItemIPv4)(nil)

// Reload implements ItemStruct interface.
func (item *ItemIPv4) Reload() {
	cptr := (*C.struct_rte_flow_item_ipv4)(item.createOrRet(C.sizeof_struct_rte_flow_item_ipv4))
	cvtIPv4Header(&cptr.hdr, &item.Header)
	runtime.SetFinalizer(item, (*ItemIPv4).free)
}

func cvtIPv4Header(dst *C.struct_rte_ipv4_hdr, src *IPv4Header) {
	setIPv4HdrVersionIHL(dst, src)

	dst.type_of_service = C.uint8_t(src.ToS)
	beU16(src.TotalLength, unsafe.Pointer(&dst.total_length))
	beU16(src.ID, unsafe.Pointer(&dst.packet_id))
	beU16(src.FragmentOffset, unsafe.Pointer(&dst.fragment_offset))
	dst.time_to_live = C.uint8_t(src.TTL)
	dst.next_proto_id = C.uint8_t(src.Proto)
	beU16(src.Checksum, unsafe.Pointer(&dst.hdr_checksum))

	dst.src_addr = *(*C.rte_be32_t)(unsafe.Pointer(&src.SrcAddr[0]))
	dst.dst_addr = *(*C.rte_be32_t)(unsafe.Pointer(&src.DstAddr[0]))
}

func setIPv4HdrVersionIHL(dst *C.struct_rte_ipv4_hdr, src *IPv4Header) {
	p := off(unsafe.Pointer(dst), C.IPv4_HDR_OFF_DST_VERSION_IHL)
	*(*C.uint8_t)(p) = (C.uchar)(src.VersionIHL)
}

// Type implements ItemStruct interface.
func (item *ItemIPv4) Type() ItemType {
	return ItemTypeIPv4
}

// Mask implements ItemStruct interface.
func (item *ItemIPv4) Mask() unsafe.Pointer {
	return unsafe.Pointer(C.get_item_ipv4_mask())
}
