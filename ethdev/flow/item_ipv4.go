package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_ip.h>
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
	"unsafe"

	"github.com/yerden/go-dpdk/common"
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

const _ uintptr = unsafe.Sizeof(IPv4Header{}) - C.sizeof_struct_rte_ipv4_hdr
const _ uintptr = C.sizeof_struct_rte_ipv4_hdr - unsafe.Sizeof(IPv4Header{})

// ItemIPv4 matches an IPv4 header.
//
// Note: IPv4 options are handled by dedicated pattern items.
type ItemIPv4 struct {
	Header IPv4Header
}

var _ ItemValue = (*ItemIPv4)(nil)

// Transform implements Action interface.
func (item *ItemIPv4) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	ptr := &C.struct_rte_flow_item_ipv4{}
	ptr.hdr = *(*C.struct_rte_ipv4_hdr)(unsafe.Pointer(&item.Header))
	return common.TransformPOD(alloc, ptr)
}

// ItemType implements ItemValue interface.
func (item *ItemIPv4) ItemType() ItemType {
	return ItemTypeIPv4
}

// DefaultMask implements ItemStruct interface.
func (item *ItemIPv4) DefaultMask() unsafe.Pointer {
	return unsafe.Pointer(&C.rte_flow_item_ipv4_mask)
}
