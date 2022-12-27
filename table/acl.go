package table

/*
#include <rte_table.h>
#include <rte_table_acl.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// ACLFieldDef is the ACL Field definition. Each field in the ACL rule
// has an associate definition. It defines the type of field, its
// size, its offset in the input buffer, the field index, and the
// input index. For performance reasons, the inner loop of the search
// function is unrolled to process four input bytes at a time. This
// requires the input to be grouped into sets of 4 consecutive bytes.
// The loop processes the first input byte as part of the setup and
// then subsequent bytes must be in groups of 4 consecutive bytes.
type ACLFieldDef struct {
	Type       uint8  // type - RTE_ACL_FIELD_TYPE_*.
	Size       uint8  // size of field 1,2,4, or 8.
	FieldIndex uint8  // index of field inside the rule.
	InputIndex uint8  // 0-N input index.
	Offset     uint32 // offset to start of field.
}

// ACLParams is the ACL table parameters.
type ACLParams struct {
	Name        string        // Name.
	Rules       uint32        // Maximum number of ACL rules in the table.
	FieldFormat []ACLFieldDef // Format specification of the fields of the ACL rule.
}

var _ Params = (*ACLParams)(nil)

// Ops implements Params interface.
func (p *ACLParams) Ops() *Ops {
	return (*Ops)(&C.rte_table_acl_ops)
}

var alloc = &common.StdAlloc{}

// Transform implements common.Transformer interface.
func (p *ACLParams) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	var params *C.struct_rte_table_acl_params
	params = (*C.struct_rte_table_acl_params)(alloc.Malloc(unsafe.Sizeof(*params)))

	if len(p.FieldFormat) > len(params.field_format) {
		panic("excessive field format entries in ACL params")
	}

	params.name = (*C.char)(common.CString(alloc, p.Name))
	params.n_rules = C.uint(p.Rules)
	params.n_rule_fields = C.uint(len(p.FieldFormat))
	for i := range p.FieldFormat {
		src := &p.FieldFormat[i]
		dst := &params.field_format[i]
		dst._type = C.uint8_t(src.Type)
		dst.size = C.uint8_t(src.Size)
		dst.field_index = C.uint8_t(src.FieldIndex)
		dst.input_index = C.uint8_t(src.InputIndex)
		dst.offset = C.uint32_t(src.Offset)
	}

	return unsafe.Pointer(params), func(p unsafe.Pointer) {
		params := (*C.struct_rte_table_acl_params)(p)
		alloc.Free(unsafe.Pointer(params.name))
		alloc.Free(unsafe.Pointer(params))
	}
}
