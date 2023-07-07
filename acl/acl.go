package acl

/*
#include <rte_acl.h>

enum {
	FIELD_OFF_U8  = offsetof(union rte_acl_field_types, u8),
	FIELD_OFF_U16 = offsetof(union rte_acl_field_types, u16),
	FIELD_OFF_U32 = offsetof(union rte_acl_field_types, u32),
	FIELD_OFF_U64 = offsetof(union rte_acl_field_types, u64),
};

RTE_ACL_RULE_DEF(sample_rule, 1);

static inline size_t
get_rule_size(int n_defs)
{
	return RTE_ACL_RULE_SZ(n_defs);
}

*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

const (
	FieldTypeMask    uint8 = C.RTE_ACL_FIELD_TYPE_MASK
	FieldTypeBitmask uint8 = C.RTE_ACL_FIELD_TYPE_BITMASK
	FieldTypeRange   uint8 = C.RTE_ACL_FIELD_TYPE_RANGE
)

func RuleSize(numDefs int) uint32 {
	return uint32(C.get_rule_size(C.int(numDefs)))
}

type FieldDef struct {
	Type, Size uint8
	FieldIndex uint8
	InputIndex uint8
	Offset     uint32
}

type Config struct {
	Categories uint32
	Defs       []FieldDef
	MaxSize    int
}

type Field struct {
	// a 1,2,4, or 8 byte value of the field.
	Value any

	// depending on field type:
	// mask -> 1.2.3.4/32 value=0x1020304, mask_range=32,
	// range -> 0 : 65535 value=0, mask_range=65535,
	// bitmask -> 0x06/0xff value=6, mask_range=0xff.
	MaskRange any
}

func setValue(data *[8]byte, v any) {
	switch v := v.(type) {
	case uint8:
		*(*uint8)(unsafe.Pointer(data)) = v
	case uint16:
		*(*uint16)(unsafe.Pointer(data)) = v
	case uint32:
		*(*uint32)(unsafe.Pointer(data)) = v
	case uint64:
		*(*uint64)(unsafe.Pointer(data)) = v
	default:
		panic("invalid type of field value")
	}

}

func (f *Field) field() C.struct_rte_acl_field {
	ret := C.struct_rte_acl_field{}
	setValue(&ret.value, f.Value)
	setValue(&ret.mask_range, f.MaskRange)
	return ret
}

type RuleData struct {
	CategoryMask uint32
	Priority     int32
	Userdata     uint32
}

type Param struct {
	Name       string
	SocketID   int
	RuleSize   uint32
	MaxRuleNum uint32
}

type Rule struct {
	Data   RuleData
	Fields []Field
}

type Context C.struct_rte_acl_ctx

func Create(p *Param) (*Context, error) {
	params := C.struct_rte_acl_param{}
	params.name = C.CString(p.Name)
	defer C.free(unsafe.Pointer(params.name))

	params.max_rule_num = C.uint32_t(p.MaxRuleNum)
	params.rule_size = C.uint32_t(p.RuleSize)
	params.socket_id = C.int(p.SocketID)

	if ctx := (*Context)(C.rte_acl_create(&params)); ctx != nil {
		return ctx, nil
	}

	return nil, common.RteErrno()
}

func (ctx *Context) Free() {
	C.rte_acl_free((*C.struct_rte_acl_ctx)(ctx))
}

func (ctx *Context) Reset() {
	C.rte_acl_reset((*C.struct_rte_acl_ctx)(ctx))
}

func (ctx *Context) Dump() {
	C.rte_acl_dump((*C.struct_rte_acl_ctx)(ctx))
}

func ListDump() {
	C.rte_acl_list_dump()
}

func (ctx *Context) ResetRules() {
	C.rte_acl_reset_rules((*C.struct_rte_acl_ctx)(ctx))
}

func FindExisting(name string) (*Context, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if ctx := (*Context)(C.rte_acl_find_existing(cname)); ctx != nil {
		return ctx, nil
	}

	return nil, common.RteErrno()
}

// data must be an array of buffers NOT allocated in Go.
func (ctx *Context) Classify(data []unsafe.Pointer, categories uint32, results []uint32) error {
	return common.IntErr(int64(C.rte_acl_classify(
		(*C.struct_rte_acl_ctx)(ctx),
		(**C.uint8_t)(unsafe.Pointer(&data[0])),
		(*C.uint32_t)(&results[0]),
		C.uint32_t(len(data)),
		C.uint32_t(categories))))
}

func (ctx *Context) Build(cfg *Config) error {
	c := C.struct_rte_acl_config{}
	c.num_categories = C.uint32_t(cfg.Categories)
	c.max_size = C.size_t(cfg.MaxSize)

	// copy field defs
	c.num_fields = C.uint32_t(len(cfg.Defs))

	if int(c.num_fields) > len(c.defs) {
		return fmt.Errorf("too many field definitions: %d > %d", len(cfg.Defs), len(c.defs))
	}

	for i := range cfg.Defs {
		def := &cfg.Defs[i]
		c.defs[i] = C.struct_rte_acl_field_def{
			_type:       C.uint8_t(def.Type),
			size:        C.uint8_t(def.Size),
			field_index: C.uint8_t(def.FieldIndex),
			input_index: C.uint8_t(def.InputIndex),
			offset:      C.uint32_t(def.Offset),
		}
	}

	return common.IntErr(int64(C.rte_acl_build((*C.struct_rte_acl_ctx)(ctx), &c)))
}

func (ctx *Context) AddRules(input []Rule) error {
	fieldsNum := len(input[0].Fields)
	ruleSize := RuleSize(fieldsNum)

	rules := make([]byte, int(ruleSize)*len(input))

	for i := range input {
		rule := (*C.struct_sample_rule)(unsafe.Add(unsafe.Pointer(&rules[0]), i*int(ruleSize)))

		ruleData := &rule.data
		ruleFields := unsafe.Slice(&rule.field[0], fieldsNum)

		ruleData.category_mask = C.uint32_t(input[i].Data.CategoryMask)
		ruleData.priority = C.int32_t(input[i].Data.Priority)
		ruleData.userdata = C.uint32_t(input[i].Data.Userdata)

		for j := range ruleFields {
			ruleFields[j] = input[i].Fields[j].field()
		}
	}

	return common.IntErr(int64(C.rte_acl_add_rules(
		(*C.struct_rte_acl_ctx)(ctx),
		(*C.struct_rte_acl_rule)(unsafe.Pointer(&rules[0])),
		C.uint32_t(len(input)),
	)))
}
