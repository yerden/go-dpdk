package table

/*
#include <rte_table_lpm.h>
#include <rte_table_lpm_ipv6.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

var (
	// LPMOps is the LPM IPv4 ops.
	LPMOps = (*Ops)(&C.rte_table_lpm_ops)

	// LPM6Ops is the LPM IPv6 ops.
	LPM6Ops = (*Ops)(&C.rte_table_lpm_ipv6_ops)
)

// LPMParams is the Hash table parameters.
type LPMParams struct {
	TableOps *Ops

	// Name.
	Name string

	// Number of rules.
	Rules uint32

	// Field is currently unused.
	NumTBL8 uint32

	// Number of bytes at the start of the table entry that uniquely
	// identify the entry. Cannot be bigger than table entry size.
	EntryUniqueSize uint32

	// Byte offset within input packet meta-data where lookup key
	// (i.e. the destination IP address) is located.
	Offset uint32
}

// Ops implements Params interface.
func (p *LPMParams) Ops() *Ops {
	return p.TableOps
}

// Transform implements common.Transformer interface.
func (p *LPMParams) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	switch p.TableOps {
	case LPMOps:
		return p.tformTo4(alloc)
	case LPM6Ops:
		return p.tformTo6(alloc)
	}

	panic("unsupported LPM Table ops")
}

func (p *LPMParams) tformTo4(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	var params *C.struct_rte_table_lpm_params
	params = (*C.struct_rte_table_lpm_params)(alloc.Malloc(unsafe.Sizeof(*params)))
	params.name = (*C.char)(common.CString(alloc, p.Name))
	params.number_tbl8s = C.uint32_t(p.NumTBL8)
	params.n_rules = C.uint32_t(p.Rules)
	params.entry_unique_size = C.uint32_t(p.EntryUniqueSize)
	params.offset = C.uint32_t(p.Offset)

	return unsafe.Pointer(params), func(p unsafe.Pointer) {
		params := (*C.struct_rte_table_lpm_params)(p)
		alloc.Free(unsafe.Pointer(params.name))
		alloc.Free(p)
	}
}

func (p *LPMParams) tformTo6(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	var params *C.struct_rte_table_lpm_ipv6_params
	params = (*C.struct_rte_table_lpm_ipv6_params)(alloc.Malloc(unsafe.Sizeof(*params)))
	params.name = (*C.char)(common.CString(alloc, p.Name))
	params.number_tbl8s = C.uint32_t(p.NumTBL8)
	params.n_rules = C.uint32_t(p.Rules)
	params.entry_unique_size = C.uint32_t(p.EntryUniqueSize)
	params.offset = C.uint32_t(p.Offset)

	return unsafe.Pointer(params), func(p unsafe.Pointer) {
		params := (*C.struct_rte_table_lpm_ipv6_params)(p)
		alloc.Free(unsafe.Pointer(params.name))
		alloc.Free(p)
	}
}

// LPMKey is the insert key for LPM IPv4 table.
type LPMKey struct {
	// IP address.
	IP uint32

	// IP address depth. The most significant "depth" bits of the IP
	// address specify the network part of the IP address, while the
	// rest of the bits specify the host part of the address and are
	// ignored for the purpose of route specification.
	Depth uint8

	// padding to fix alignment
	_ [3]byte
}

type cLPMKey C.struct_rte_table_lpm_key

// LPM6Key is the insert key for LPM IPv4 table.
type LPM6Key struct {
	// IP address.
	IP [C.RTE_LPM_IPV6_ADDR_SIZE]byte

	// IP address depth. The most significant "depth" bits of the IP
	// address specify the network part of the IP address, while the
	// rest of the bits specify the host part of the address and are
	// ignored for the purpose of route specification.
	Depth uint8
}

type cLPM6Key C.struct_rte_table_lpm_ipv6_key

var _ = []uintptr{
	unsafe.Sizeof(LPMKey{}) - unsafe.Sizeof(cLPMKey{}),
	unsafe.Sizeof(cLPMKey{}) - unsafe.Sizeof(LPMKey{}),

	unsafe.Sizeof(LPM6Key{}) - unsafe.Sizeof(cLPM6Key{}),
	unsafe.Sizeof(cLPM6Key{}) - unsafe.Sizeof(LPM6Key{}),
}
