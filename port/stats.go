package port

/*
#include <rte_port.h>
*/
import "C"
import (
	"unsafe"
)

// InStats is an input port statistics.
type InStats struct {
	PacketsIn   uint64
	PacketsDrop uint64
}

var _ uintptr = unsafe.Sizeof(InStats{}) - unsafe.Sizeof(C.struct_rte_port_in_stats{})
var _ uintptr = unsafe.Sizeof(C.struct_rte_port_in_stats{}) - unsafe.Sizeof(InStats{})

// OutStats is an output port statistics.
type OutStats struct {
	PacketsIn   uint64
	PacketsDrop uint64
}

var _ uintptr = unsafe.Sizeof(OutStats{}) - unsafe.Sizeof(C.struct_rte_port_out_stats{})
var _ uintptr = unsafe.Sizeof(C.struct_rte_port_out_stats{}) - unsafe.Sizeof(OutStats{})
