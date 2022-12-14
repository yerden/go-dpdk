package table

/*
#include <rte_table_acl.h>

enum {
	ACL_FIELD_DIM_64 = sizeof(struct rte_acl_field) / sizeof(uint64_t),
};

*/
import "C"

import (
	"unsafe"
)

// ACLField describes a field to match.
type ACLField [C.ACL_FIELD_DIM_64]uint64
type cACLField C.struct_rte_acl_field

// NewACLField8 returns ACLField with 8-bit value and range mask.
func NewACLField8(value, maskRange uint8) ACLField {
	var ret ACLField
	dst := (*C.struct_rte_acl_field)(unsafe.Pointer(&ret))
	*(*uint8)(unsafe.Pointer(&dst.value)) = value
	*(*uint8)(unsafe.Pointer(&dst.mask_range)) = maskRange
	return ret
}

// NewACLField16 returns ACLField with 16-bit value and range mask.
func NewACLField16(value, maskRange uint16) ACLField {
	var ret ACLField
	dst := (*C.struct_rte_acl_field)(unsafe.Pointer(&ret))
	*(*uint16)(unsafe.Pointer(&dst.value)) = value
	*(*uint16)(unsafe.Pointer(&dst.mask_range)) = maskRange
	return ret
}

// NewACLField32 returns ACLField with 32-bit value and range mask.
func NewACLField32(value, maskRange uint32) ACLField {
	var ret ACLField
	dst := (*C.struct_rte_acl_field)(unsafe.Pointer(&ret))
	*(*uint32)(unsafe.Pointer(&dst.value)) = value
	*(*uint32)(unsafe.Pointer(&dst.mask_range)) = maskRange
	return ret
}

// NewACLField64 returns ACLField with 64-bit value and range mask.
func NewACLField64(value, maskRange uint64) ACLField {
	var ret ACLField
	dst := (*C.struct_rte_acl_field)(unsafe.Pointer(&ret))
	*(*uint64)(unsafe.Pointer(&dst.value)) = value
	*(*uint64)(unsafe.Pointer(&dst.mask_range)) = maskRange
	return ret
}

// ACLRuleAdd is the key for adding a key/value to a table.
type ACLRuleAdd struct {
	Priority   int32
	FieldValue [C.RTE_ACL_MAX_FIELDS]ACLField
}
type cACLRuleAdd C.struct_rte_table_acl_rule_add_params

// ACLRuleDelete is the key for deleting a key/value from a table.
type ACLRuleDelete struct {
	FieldValue [C.RTE_ACL_MAX_FIELDS]ACLField
}
type cACLRuleDelete C.struct_rte_table_acl_rule_delete_params

var _ = []uintptr{
	unsafe.Sizeof(ACLField{}) - unsafe.Sizeof(cACLField{}),
	unsafe.Sizeof(cACLField{}) - unsafe.Sizeof(ACLField{}),

	unsafe.Sizeof(ACLRuleAdd{}) - unsafe.Sizeof(cACLRuleAdd{}),
	unsafe.Sizeof(cACLRuleAdd{}) - unsafe.Sizeof(ACLRuleAdd{}),

	unsafe.Sizeof(ACLRuleDelete{}) - unsafe.Sizeof(cACLRuleDelete{}),
	unsafe.Sizeof(cACLRuleDelete{}) - unsafe.Sizeof(ACLRuleDelete{}),
}
