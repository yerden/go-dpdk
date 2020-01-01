package port

/*
#include <rte_config.h>
#include <rte_port.h>

void *go_rd_create(void *ops_table, void *params, int socket_id)
{
	struct rte_port_in_ops *ops = ops_table;
	return ops->f_create(params, socket_id);
}

int go_rd_rx(void *ops_table, void *port, struct rte_mbuf **pkts, uint32_t n_pkts)
{
	struct rte_port_in_ops *ops = ops_table;
	return ops->f_rx(port, pkts, n_pkts);
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

int go_wr_tx(void *ops_table, void *port, struct rte_mbuf *pkt)
{
	struct rte_port_out_ops *ops = ops_table;
	return ops->f_tx(port, pkt);
}

int go_wr_tx_bulk(void *ops_table, void *port, struct rte_mbuf **pkts, uint64_t pkts_mask)
{
	struct rte_port_out_ops *ops = ops_table;
	return ops->f_tx_bulk(port, pkts, pkts_mask);
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

	"github.com/yerden/go-dpdk/mbuf"
)

// InOps describes input port interface defining the input port
// operation.
type InOps C.struct_rte_port_in_ops

// Create creates new port instance.
func (ops *InOps) Create(socket int, arg *InArg) *In {
	return (*In)(C.go_rd_create(unsafe.Pointer(ops), unsafe.Pointer(arg), C.int(socket)))
}

// Free releases all memory allocated when creating port instance.
func (ops *InOps) Free(p *In) error {
	return err(C.go_rd_free(unsafe.Pointer(ops), unsafe.Pointer(p)))
}

// Rx receives packets into specified array.
func (ops *InOps) Rx(p *In, pkts []*mbuf.Mbuf) int {
	return int(C.go_rd_rx(unsafe.Pointer(ops), unsafe.Pointer(p),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint32_t(len(pkts))))
}

// OutOps describes output port interface defining the output port
// operation.
type OutOps C.struct_rte_port_out_ops

// Create creates new port instance.
func (ops *OutOps) Create(socket int, arg *OutArg) *Out {
	return (*Out)(C.go_wr_create(unsafe.Pointer(ops), unsafe.Pointer(arg), C.int(socket)))
}

// Free releases all memory allocated when creating port instance.
func (ops *OutOps) Free(p *Out) error {
	return err(C.go_wr_free(unsafe.Pointer(ops), unsafe.Pointer(p)))
}

// Tx submits given packet via port instance.
func (ops *OutOps) Tx(p *Out, pkt *mbuf.Mbuf) int {
	return int(C.go_wr_tx(unsafe.Pointer(ops), unsafe.Pointer(p),
		(*C.struct_rte_mbuf)(unsafe.Pointer(pkt))))
}

// TxBulk submits given packets via port instance according to specified mask.
// if n-th bit of a mask is set, n-th mbuf from pkts is considered valid.
func (ops *OutOps) TxBulk(p *Out, pkts []*mbuf.Mbuf, mask uint64) int {
	return int(C.go_wr_tx_bulk(unsafe.Pointer(ops), unsafe.Pointer(p),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint64_t(mask)))
}
