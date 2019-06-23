/*
Package port wraps RTE port library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package port

/*
#include <rte_config.h>
#include <rte_port.h>
*/
import "C"

import (
	"unsafe"
)

// ReaderOps describes input port interface defining the input port
// operation.
type ReaderOps C.struct_rte_port_in_ops

// WriterOps describes output port interface defining the output port
// operation.
type WriterOps C.struct_rte_port_out_ops

// Reader implements reader port capability which allows to read
// packets from it.
type Reader interface {
	// ReaderOps returns pointer to statically allocated call table.
	ReaderOps() *ReaderOps
	// NewArg allocates an opaque argument which is required by
	// ReaderOps.
	NewArg() unsafe.Pointer
}

// Writer implements writer port capability which allows to write
// packets to it.
type Writer interface {
	// WriterOps returns pointer to statically allocated call table.
	WriterOps() *WriterOps
	// NewArg allocates an opaque argument which is required by
	// WriterOps.
	NewArg() unsafe.Pointer
}
