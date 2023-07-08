/*
Package acl wraps DPDK ACL classified library.
*/
package acl

/*
#include <rte_acl.h>

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

// Type of the field.
const (
	FieldTypeMask    uint8 = C.RTE_ACL_FIELD_TYPE_MASK
	FieldTypeBitmask uint8 = C.RTE_ACL_FIELD_TYPE_BITMASK
	FieldTypeRange   uint8 = C.RTE_ACL_FIELD_TYPE_RANGE
)

// RuleSize returns amount of memory occupied by a rule with numDefs
// fields.
func RuleSize(numDefs int) uint32 {
	return uint32(C.get_rule_size(C.int(numDefs)))
}

// FieldDef is an ACL Field definition.
//
// Each field in the ACL rule has an associate definition.  It defines
// the type of field, its size, its offset in the input buffer, the
// field index, and the input index.  For performance reasons, the
// inner loop of the search function is unrolled to process four input
// bytes at a time. This requires the input to be grouped into sets of
// 4 consecutive bytes. The loop processes the first input byte as
// part of the setup and then subsequent bytes must be in groups of 4
// consecutive bytes.
type FieldDef struct {
	Type, Size uint8
	FieldIndex uint8
	InputIndex uint8
	Offset     uint32
}

// Config is an ACL build configuration. Defines the fields of an ACL
// trie and number of categories to build with.
type Config struct {
	Categories uint32
	Defs       []FieldDef
	MaxSize    int
}

// Field defines the value of a field for a rule.
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

// RuleData contains miscellaneous data for ACL rule.
type RuleData struct {
	CategoryMask uint32
	Priority     int32
	Userdata     uint32
}

// Param contains parameters used when creating the ACL context.
type Param struct {
	Name       string
	SocketID   int
	RuleSize   uint32
	MaxRuleNum uint32
}

// Rule contains a rule to store in ACL content.
type Rule struct {
	Data   RuleData
	Fields []Field
}

// Context is an ACL context containing rules and optimized built trie
// to search.
type Context C.struct_rte_acl_ctx

// Create new Context using specified params. Returns new instance of
// Context and an error.
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

// Free de-allocates all memory used by ACL context.
func (ctx *Context) Free() {
	C.rte_acl_free((*C.struct_rte_acl_ctx)(ctx))
}

// Reset deletes all rules from the ACL context and destroy all
// internal run-time structures. This function is not multi-thread
// safe.
func (ctx *Context) Reset() {
	C.rte_acl_reset((*C.struct_rte_acl_ctx)(ctx))
}

// Dump an ACL context structure to the console.
func (ctx *Context) Dump() {
	C.rte_acl_dump((*C.struct_rte_acl_ctx)(ctx))
}

// ListDump dumps all ACL context structures to the console.
func ListDump() {
	C.rte_acl_list_dump()
}

// ResetRules delete all rules from the ACL context. This function is
// not multi-thread safe. Note that internal run-time structures are
// not affected.
func (ctx *Context) ResetRules() {
	C.rte_acl_reset_rules((*C.struct_rte_acl_ctx)(ctx))
}

// FindExisting finds an existing ACL context object and return a
// pointer to it.
func FindExisting(name string) (*Context, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	if ctx := (*Context)(C.rte_acl_find_existing(cname)); ctx != nil {
		return ctx, nil
	}

	return nil, common.RteErrno()
}

// Classify performs search for a matching ACL rule for each input
// data buffer. Each input data buffer can have up to *categories*
// matches.  That implies that results array should be big enough to
// hold (categories * len(data)) elements.  Also categories parameter
// should be either one or multiple of RTE_ACL_RESULTS_MULTIPLIER and
// can't be bigger than RTE_ACL_MAX_CATEGORIES.  If more than one rule
// is applicable for given input buffer and given category, then rule
// with highest priority will be returned as a match.  Note, that it
// is a caller's responsibility to ensure that input parameters are
// valid and point to correct memory locations.
//
// data must be an array of buffers NOT allocated in Go.
func (ctx *Context) Classify(data []unsafe.Pointer, categories uint32, results []uint32) error {
	return common.IntErr(int64(C.rte_acl_classify(
		(*C.struct_rte_acl_ctx)(ctx),
		(**C.uint8_t)(unsafe.Pointer(&data[0])),
		(*C.uint32_t)(&results[0]),
		C.uint32_t(len(data)),
		C.uint32_t(categories))))
}

// Build analyzes set of rules and build required internal run-time
// structures. This function is not multi-thread safe.
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

// AddRules adds rules to an existing ACL context. This function is
// not multi-thread safe.
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
