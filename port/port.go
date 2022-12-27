package port

/*
#include <rte_config.h>
#include <rte_port.h>

static void *
go_in_create(struct rte_port_in_ops *ops, void *arg, int socket)
{
	return ops->f_create(arg, socket);
}

static int
go_in_free(struct rte_port_in_ops *ops, void *port)
{
	return ops->f_free(port);
}

static void *
go_out_create(struct rte_port_out_ops *ops, void *arg, int socket)
{
	return ops->f_create(arg, socket);
}

static int
go_out_free(struct rte_port_out_ops *ops, void *port)
{
	return ops->f_free(port);
}

*/
import "C"
import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

type (
	// InOps is the function table which implements input port.
	InOps C.struct_rte_port_in_ops

	// In is the instantiated input port.
	In [0]byte
)

type (
	// OutOps is the function table which implements output port.
	OutOps C.struct_rte_port_out_ops

	// Out is the instantiated output port.
	Out [0]byte
)

// InParams describes the input port interface.
type InParams interface {
	// Returns allocated opaque structure along with its destructor.
	// Since InParams describes Go implementation of the port
	// configuration this member allocates its C counterpart as stated
	// in DPDK rte_port.
	common.Transformer

	// Returns pointer to associated rte_port_in_ops.
	InOps() *InOps
}

// OutParams describes configuration and behaviour of output port.
type OutParams interface {
	// Returns allocated opaque structure argument along with its
	// destructor. It is used with ops function table.
	common.Transformer

	// Returns pointer to associated rte_port_out_ops.
	OutOps() *OutOps
}

var alloc = &common.StdAlloc{}

// CreateIn creates input port for specified socket and configuration.
//
// It may return nil in case of an error.
func CreateIn(socket int, params InParams) *In {
	ops := (*C.struct_rte_port_in_ops)(params.InOps())

	arg, dtor := params.Transform(alloc)
	defer dtor(arg)

	return (*In)(C.go_in_create(ops, arg, C.int(socket)))
}

// Free destroys created input port.
func (port *In) Free(inOps *InOps) error {
	ops := (*C.struct_rte_port_in_ops)(inOps)
	return common.IntErr(int64(C.go_in_free(ops, unsafe.Pointer(port))))
}

// CreateOut creates output port for specified socket and
// configuration.
//
// It may return nil in case of an error.
func CreateOut(socket int, params OutParams) *Out {
	ops := (*C.struct_rte_port_out_ops)(params.OutOps())

	arg, dtor := params.Transform(alloc)
	defer dtor(arg)

	return (*Out)(C.go_out_create(ops, arg, C.int(socket)))
}

// Free destroys created output port.
func (port *Out) Free(outOps *OutOps) error {
	ops := (*C.struct_rte_port_out_ops)(outOps)
	return common.IntErr(int64(C.go_out_free(ops, unsafe.Pointer(port))))
}
