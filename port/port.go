/*
Package port wraps RTE port library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package port

/*
#include <rte_config.h>
#include <rte_port.h>

void *go_rd_create(void *ops_table, void *params, int socket_id)
{
	struct rte_port_in_ops *ops = ops_table;
	return ops->f_create(params, socket_id);
}

int go_rd_free(void *ops_table, void *port)
{
	struct rte_port_in_ops *ops = ops_table;
	return ops->f_free(port);
}

void *go_wr_create(void *ops_table, void *params, int socket_id)
{
	struct rte_port_out_ops *ops = ops_table;
	return ops->f_create(params, socket_id);
}

int go_wr_free(void *ops_table, void *port)
{
	struct rte_port_out_ops *ops = ops_table;
	return ops->f_free(port);
}
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// ReaderOps describes input port interface defining the input port
// operation.
type ReaderOps C.struct_rte_port_in_ops

// WriterOps describes output port interface defining the output port
// operation.
type WriterOps C.struct_rte_port_out_ops

// ReaderParams implements reader port capability which allows to read
// packets from it.
type ReaderParams interface {
	// ReaderOps returns pointer to statically allocated call table.
	// and an opaque argument which is required to create port.
	ReaderOps() (ops *ReaderOps, arg unsafe.Pointer)
}

// WriterParams implements writer port capability which allows to
// write packets to it.
type WriterParams interface {
	// WriterOps returns pointer to statically allocated call table.
	// and an opaque argument which is required to create port.
	WriterOps() (ops *WriterOps, arg unsafe.Pointer)
}

// Reader is the instance of reader port.
type Reader struct {
	Ops  *ReaderOps
	Port unsafe.Pointer
}

// Writer is the instance of writer port.
type Writer struct {
	Ops  *WriterOps
	Port unsafe.Pointer
}

func err(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

// XXX: we need to wrap calls which are not performance bottlenecks.

// NewReader creates new Reader. ReaderParams and destination NUMA
// socket must be specified. In case of an error, nil is returned.
func NewReader(p ReaderParams, socket int) *Reader {
	ops, arg := p.ReaderOps()
	port := C.go_rd_create(unsafe.Pointer(ops), arg, C.int(socket))
	if port == nil {
		return nil
	}
	return &Reader{ops, port}
}

// Free destroys once created Reader port.
func (rd *Reader) Free() error {
	return err(C.go_rd_free(unsafe.Pointer(rd.Ops), rd.Port))
}

// NewWriter creates new Writer. ReaderParams and destination NUMA
// socket must be specified. In case of an error, nil is returned.
func NewWriter(p WriterParams, socket int) *Writer {
	ops, arg := p.WriterOps()
	port := C.go_wr_create(unsafe.Pointer(ops), arg, C.int(socket))
	if port == nil {
		return nil
	}
	return &Writer{ops, port}
}

// Free destroys once created Writer port.
func (wr *Writer) Free() error {
	return err(C.go_wr_free(unsafe.Pointer(wr.Ops), wr.Port))
}
