package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
#include <rte_udp.h>

*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// UDPHeader represents UDP header format.
type UDPHeader struct {
	SrcPort  uint16 /* UDP source port. */
	DstPort  uint16 /* UDP destination port. */
	Length   uint16 /* UDP datagram length */
	Checksum uint16 /* UDP datagram checksum */
}

const _ uintptr = unsafe.Sizeof(UDPHeader{}) - C.sizeof_struct_rte_udp_hdr
const _ uintptr = C.sizeof_struct_rte_udp_hdr - unsafe.Sizeof(UDPHeader{})

// ItemUDP matches an UDP header.
type ItemUDP struct {
	Header UDPHeader
}

var _ ItemValue = (*ItemUDP)(nil)

// Transform implements Action interface.
func (item *ItemUDP) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	ptr := &C.struct_rte_flow_item_udp{}
	ptr.hdr = *(*C.struct_rte_udp_hdr)(unsafe.Pointer(&item.Header))
	return common.TransformPOD(alloc, ptr)
}

// ItemType implements ItemValue interface.
func (item *ItemUDP) ItemType() ItemType {
	return ItemTypeUDP
}

// DefaultMask implements ItemStruct interface.
func (item *ItemUDP) DefaultMask() unsafe.Pointer {
	return unsafe.Pointer(&C.rte_flow_item_udp_mask)
}
