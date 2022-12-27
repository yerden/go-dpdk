package table

/*
#include <rte_table.h>
*/
import "C"
import (
	"unsafe"
)

// Stats is a table statistics.
type Stats struct {
	PacketsIn         uint64
	PacketsLookupMiss uint64
}

var _ uintptr = unsafe.Sizeof(Stats{}) - unsafe.Sizeof(C.struct_rte_table_stats{})
var _ uintptr = unsafe.Sizeof(C.struct_rte_table_stats{}) - unsafe.Sizeof(Stats{})
