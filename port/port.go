package port

/*
#include <rte_config.h>
#include <rte_port.h>

struct rx_port {
	struct rte_port_in_ops *ops;
	void *port;
};

struct tx_port {
	struct rte_port_out_ops *ops;
	void *port;
};

void *go_rd_create(struct rte_port_in_ops *ops, void *params, int socket_id)
{
	return ops->f_create(params, socket_id);
}

int go_rd_rx(struct rx_port *rx, struct rte_mbuf **pkts, uint32_t n_pkts)
{
	struct rte_port_in_ops *ops  = rx->ops;
	return ops->f_rx(rx->port, pkts, n_pkts);
}

int go_rd_free(struct rx_port *rx)
{
	struct rte_port_in_ops *ops  = rx->ops;
	return ops->f_free(rx->port);
}

void *go_wr_create(struct rte_port_out_ops *ops, void *params, int socket_id)
{
	return ops->f_create(params, socket_id);
}

int go_wr_tx(struct tx_port *tx, struct rte_mbuf *pkt)
{
	struct rte_port_out_ops *ops  = tx->ops;
	return ops->f_tx(tx->port, pkt);
}

int go_wr_tx_bulk(struct tx_port *tx, struct rte_mbuf **pkts, uint64_t pkts_mask)
{
	struct rte_port_out_ops *ops  = tx->ops;
	return ops->f_tx_bulk(tx->port, pkts, pkts_mask);
}

int go_wr_free(struct tx_port *tx)
{
	struct rte_port_out_ops *ops  = tx->ops;
	return ops->f_free(tx->port);
}
*/
import "C"

import (
	"errors"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/mbuf"
)

var (
	// ErrPortCreate signifies error while attempting to call
	// ops->f_create.
	ErrPortCreate = errors.New("port not created")
)

// Rx describes input port interface defining the input port
// operation.
type Rx C.struct_rx_port

// Tx describes output port interface defining the output port
// operation.
type Tx C.struct_tx_port

// RxFactory is used to create new Rx port.
type RxFactory interface {
	// CreateRx creates new port instance.
	CreateRx(socket int) (*Rx, error)
}

// TxFactory is used to create new Tx port.
type TxFactory interface {
	// CreateTx creates new port instance.
	CreateTx(socket int) (*Tx, error)
}

func err(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

// doCreate is a helper to implement RxFactory.
func (rx *Rx) doCreate(socket int, arg unsafe.Pointer) error {
	rx.port = C.go_rd_create(rx.ops, arg, C.int(socket))
	if rx.port == nil {
		return ErrPortCreate
	}

	return nil
}

// Rx receives packets into specified array. pkts must have length > 0.
func (rx *Rx) Rx(pkts []*mbuf.Mbuf) int {
	return int(C.go_rd_rx((*C.struct_rx_port)(rx),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint32_t(len(pkts))))
}

// Free releases all memory allocated when creating port instance.
func (rx *Rx) Free() error {
	return err(C.go_rd_free((*C.struct_rx_port)(rx)))
}

// doCreate is a helper to implement TxFactory.
func (tx *Tx) doCreate(socket int, arg unsafe.Pointer) error {
	tx.port = C.go_wr_create(tx.ops, arg, C.int(socket))
	if tx.port == nil {
		return ErrPortCreate
	}

	return nil
}

// Tx submits given packet via port instance.
func (tx *Tx) Tx(pkt *mbuf.Mbuf) int {
	return int(C.go_wr_tx((*C.struct_tx_port)(tx),
		(*C.struct_rte_mbuf)(unsafe.Pointer(pkt))))
}

// TxBulk submits given packets via port instance according to specified mask.
// if n-th bit of a mask is set, n-th mbuf from pkts is considered valid.
func (tx *Tx) TxBulk(pkts []*mbuf.Mbuf, mask uint64) int {
	return int(C.go_wr_tx_bulk((*C.struct_tx_port)(tx),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint64_t(mask)))
}

// Free releases all memory allocated when creating port instance.
func (tx *Tx) Free() error {
	return err(C.go_wr_free((*C.struct_tx_port)(tx)))
}
