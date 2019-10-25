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

	"github.com/google/gopacket/layers"
)

// ItemUDP matches an UDP header.
type ItemUDP struct {
	cPointer

	Header layers.UDP
}

var _ ItemStruct = (*ItemUDP)(nil)

// Reload implements ItemStruct interface.
func (item *ItemUDP) Reload() {
	cptr := (*C.struct_rte_flow_item_udp)(item.createOrRet(C.sizeof_struct_rte_flow_item_udp))
	cvtUDPHeader(&cptr.hdr, &item.Header)
	runtime.SetFinalizer(item, (*ItemUDP).free)
}

func cvtUDPHeader(dst *C.struct_rte_udp_hdr, src *layers.UDP) {
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
