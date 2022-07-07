package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>

static const struct rte_flow_item_udp *get_item_udp_mask() {
	return &rte_flow_item_udp_mask;
}

*/
import "C"
import (
	"runtime"
	"unsafe"
)

// UDPHeader represents UDP header format.
type UDPHeader struct {
	SrcPort  uint16 /* UDP source port. */
	DstPort  uint16 /* UDP destination port. */
	Length   uint16 /* UDP datagram length */
	Checksum uint16 /* UDP datagram checksum */
}

// ItemUDP matches an UDP header.
type ItemUDP struct {
	cPointer

	Header UDPHeader
}

var _ ItemStruct = (*ItemUDP)(nil)

// Reload implements ItemStruct interface.
func (item *ItemUDP) Reload() {
	cptr := (*C.struct_rte_flow_item_udp)(item.createOrRet(C.sizeof_struct_rte_flow_item_udp))
	cvtUDPHeader(&cptr.hdr, &item.Header)
	runtime.SetFinalizer(item, (*ItemUDP).free)
}

func cvtUDPHeader(dst *C.struct_rte_udp_hdr, src *UDPHeader) {
	beU16(uint16(src.SrcPort), unsafe.Pointer(&dst.src_port))
	beU16(uint16(src.DstPort), unsafe.Pointer(&dst.dst_port))
	beU16(src.Length, unsafe.Pointer(&dst.dgram_len))
	beU16(src.Checksum, unsafe.Pointer(&dst.dgram_cksum))
}

// Type implements ItemStruct interface.
func (item *ItemUDP) Type() ItemType {
	return ItemTypeUDP
}

// Mask implements ItemStruct interface.
func (item *ItemUDP) Mask() unsafe.Pointer {
	return unsafe.Pointer(C.get_item_udp_mask())
}
