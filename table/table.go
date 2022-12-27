/*
Package table implements RTE Table.

This tool is part of the DPDK Packet Framework tool suite and provides a
standard interface to implement different types of lookup tables for data plane
processing.

Virtually any search algorithm that can uniquely associate data to a lookup key
can be fitted under this lookup table abstraction. For the flow table use-case,
the lookup key is an n-tuple of packet fields that uniquely identifies a
traffic flow, while data represents actions and action meta-data associated
with the same traffic flow.
*/
package table

/*
#cgo pkg-config: libdpdk
#include <rte_table.h>

static void *
go_table_create(struct rte_table_ops *ops, void *params, int socket_id, uint32_t entry_size)
{
	return ops->f_create(params, socket_id, entry_size);
}

static int
go_table_free(struct rte_table_ops *ops, void *table)
{
	return ops->f_free(table);
}

*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Ops is the function table implementing table.
type Ops C.struct_rte_table_ops

// Table is the instance of a lookup table.
type Table [0]byte

// Params is the parameters for creating a table.
type Params interface {
	common.Transformer
	Ops() *Ops
}

// Create a table manually. May return nil pointer in case of an
// error.
func Create(socket int, p Params, entrySize uint32) *Table {
	ops := (*C.struct_rte_table_ops)(p.Ops())
	arg, dtor := p.Transform(alloc)
	defer dtor(arg)
	return (*Table)(C.go_table_create(ops, arg, C.int(socket), C.uint32_t(entrySize)))
}

// Free deletes a table.
func (t *Table) Free(tableOps *Ops) error {
	ops := (*C.struct_rte_table_ops)(tableOps)
	return common.IntErr(int64(C.go_table_free(ops, unsafe.Pointer(t))))
}
