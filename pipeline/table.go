package pipeline

/*
#include <rte_pipeline.h>

enum {
	TABLE_ENTRY_PORT_ID = offsetof(struct rte_pipeline_table_entry, port_id),
	TABLE_ENTRY_TABLE_ID = offsetof(struct rte_pipeline_table_entry, table_id),
	TABLE_ENTRY_ACTION_DATA = offsetof(struct rte_pipeline_table_entry, action_data[0]),
};

*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/table"
)

// Table represents created table in a pipeline.
type Table C.uint32_t

// TableActionHitFunc is executed on the hit in table.
type TableActionHitFunc C.rte_pipeline_table_action_handler_hit

// TableActionMissFunc is executed on the miss in table.
type TableActionMissFunc C.rte_pipeline_table_action_handler_miss

// TableEntry specifies the entry in table.
type TableEntry C.struct_rte_pipeline_table_entry

// TableParams configures table creation in the pipeline.
type TableParams struct {
	Params table.Params

	// OnHit and OnMiss action handlers must be implemented in C.
	OnHit  TableActionHitFunc
	OnMiss TableActionMissFunc

	// Argument for OnHit/OnMiss action handlers.
	ActionArg unsafe.Pointer

	// ActionDataSize specifies size of action data in an entry.
	ActionDataSize uintptr
}

// Action specifies what to do with packets on table output.
type Action uint32

// Reserved actions.
const (
	// Drop packet.
	ActionDrop Action = C.RTE_PIPELINE_ACTION_DROP

	// Send to port specified in TableEntry.
	ActionPort Action = C.RTE_PIPELINE_ACTION_PORT

	// Send to port specified in packet metadata.
	ActionPortMeta Action = C.RTE_PIPELINE_ACTION_PORT_META

	// Send to table specified in TableEntry.
	ActionTable Action = C.RTE_PIPELINE_ACTION_TABLE
)

// SetAction sets configured action in the entry.
func (entry *TableEntry) SetAction(code Action) {
	entry.action = uint32(code)
}

// GetAction returns configured action in the entry.
func (entry *TableEntry) GetAction() (code uint32) {
	return entry.action
}

// SetPortID sets configured output port ID (meta-data for "Send
// packet to output port" action).
func (entry *TableEntry) SetPortID(port PortOut) {
	p := (*PortOut)(unsafe.Add(unsafe.Pointer(entry), C.TABLE_ENTRY_PORT_ID))
	*p = port
}

// GetPortID returns configured output port ID (meta-data for "Send
// packet to output port" action).
func (entry *TableEntry) GetPortID() (port PortOut) {
	p := (*PortOut)(unsafe.Add(unsafe.Pointer(entry), C.TABLE_ENTRY_PORT_ID))
	return *p
}

// SetTableID sets configured table in the entry for "Send packet to
// table" action.
func (entry *TableEntry) SetTableID(table Table) {
	p := (*Table)(unsafe.Add(unsafe.Pointer(entry), C.TABLE_ENTRY_TABLE_ID))
	*p = table
}

// GetTableID returns configured table in the entry for "Send packet
// to table" action.
func (entry *TableEntry) GetTableID() Table {
	p := (*Table)(unsafe.Add(unsafe.Pointer(entry), C.TABLE_ENTRY_TABLE_ID))
	return *p
}

// GetActionData returns pointer to action user data.
func (entry *TableEntry) GetActionData() unsafe.Pointer {
	return unsafe.Add(unsafe.Pointer(entry), C.TABLE_ENTRY_ACTION_DATA)
}

// TableCreate creates new table in the pipeline.
func (pl *Pipeline) TableCreate(p *TableParams) (id Table, err error) {
	params := &C.struct_rte_pipeline_table_params{}
	params.ops = (*C.struct_rte_table_ops)(unsafe.Pointer(p.Params.Ops()))

	{
		x, dtor := p.Params.Transform(alloc)
		defer dtor(x)
		params.arg_create = x
	}

	params.f_action_hit = (C.rte_pipeline_table_action_handler_hit)(p.OnHit)
	params.f_action_miss = (C.rte_pipeline_table_action_handler_miss)(p.OnMiss)
	params.arg_ah = p.ActionArg

	params.action_data_size = C.uint(p.ActionDataSize)

	rc := C.rte_pipeline_table_create(
		(*C.struct_rte_pipeline)(pl),
		params,
		(*C.uint32_t)(&id))

	return id, common.IntErr(int64(rc))
}

// NewTableEntry creates new table entry template allocated in Go. Use
// actionDataSize as specified during table creation.
func NewTableEntry(actionDataSize uintptr) *TableEntry {
	b := make([]byte, unsafe.Sizeof(TableEntry{})+actionDataSize)
	return (*TableEntry)(unsafe.Pointer(&b[0]))
}

// TableDefaultEntryAdd adds default entry to table in the pipeline.
// Returns pointer to actual TableEntry instance in the pipeline's
// table.
func (pl *Pipeline) TableDefaultEntryAdd(table Table, entry *TableEntry) (*TableEntry, error) {
	var p *C.struct_rte_pipeline_table_entry
	rc := C.rte_pipeline_table_default_entry_add((*C.struct_rte_pipeline)(pl), C.uint32_t(table), (*C.struct_rte_pipeline_table_entry)(entry), &p)
	return (*TableEntry)(p), common.IntErr(int64(rc))
}

// TableDefaultEntryDelete deletes default entry from table in the pipeline.
// entry then is filled with removed entry.
func (pl *Pipeline) TableDefaultEntryDelete(table Table, entry *TableEntry) error {
	rc := C.rte_pipeline_table_default_entry_delete((*C.struct_rte_pipeline)(pl), C.uint32_t(table), (*C.struct_rte_pipeline_table_entry)(entry))
	return common.IntErr(int64(rc))
}

// TableEntryAdd adds key/entry into table. On successful invocation,
// ptrEntry points to found entry. Returns true of entry was present
// in the map prior to invocation.
//
// pkey is table-specific pointer to a struct containing a key. Please
// refer to table's API for key specification.
func (pl *Pipeline) TableEntryAdd(table Table, pkey unsafe.Pointer, entry *TableEntry, ptrEntry **TableEntry) (bool, error) {
	var found C.int

	rc := C.rte_pipeline_table_entry_add(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(table),
		pkey,
		(*C.struct_rte_pipeline_table_entry)(entry),
		&found,
		(**C.struct_rte_pipeline_table_entry)(unsafe.Pointer(ptrEntry)))
	return found != 0, common.IntErr(int64(rc))
}

// TableEntryDelete removes key from table. entry is filled with
// removed entry.
//
// pkey is table-specific pointer to a struct containing a key. Please
// refer to table's API for key specification.
func (pl *Pipeline) TableEntryDelete(table Table, pkey unsafe.Pointer, entry *TableEntry) (bool, error) {
	var found C.int

	rc := C.rte_pipeline_table_entry_delete(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(table),
		pkey,
		&found,
		(*C.struct_rte_pipeline_table_entry)(entry))
	return found != 0, common.IntErr(int64(rc))
}
