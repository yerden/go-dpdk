package ethdev

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_ethdev.h>
*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/yerden/go-dpdk/mbuf"
)

// TxBuffer is a structure used to buffer packets for future TX Used by APIs
// rte_eth_tx_buffer and rte_eth_tx_buffer_flush.
type TxBuffer C.struct_rte_eth_dev_tx_buffer

func fracCeilRound(x int, y int) int {
	n := x / y
	if y*n != x {
		n++
	}
	return n
}

// NewTxBuffer creates new TxBuffer with cnt mbufs.
func NewTxBuffer(cnt int) *TxBuffer {
	amount := TxBufferSize(cnt)
	ptrSize := unsafe.Sizeof(&mbuf.Mbuf{})
	data := make([]*mbuf.Mbuf, fracCeilRound(int(amount), int(ptrSize)))
	buf := (*TxBuffer)(unsafe.Pointer(&data[0]))
	buf.Init(cnt)
	return buf
}

// TxBufferSize is the occupied memory for TxBuffer of length mbufCnt.
func TxBufferSize(mbufCnt int) uintptr {
	return unsafe.Sizeof(TxBuffer{}) + uintptr(mbufCnt)*unsafe.Sizeof(&mbuf.Mbuf{})
}

// RxBurst receives packets for port pid and queue qid. Returns number of
// packets retrieved into pkts.
func (pid Port) RxBurst(qid uint16, pkts []*mbuf.Mbuf) uint16 {
	return uint16(C.rte_eth_rx_burst(C.uint16_t(pid), C.uint16_t(qid),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint16_t(len(pkts))))
}

// TxBurst sends packets over port pid and queue qid. Returns number of
// packets sent from pkts.
func (pid Port) TxBurst(qid uint16, pkts []*mbuf.Mbuf) uint16 {
	return uint16(C.rte_eth_tx_burst(C.uint16_t(pid), C.uint16_t(qid),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint16_t(len(pkts))))
}

// TxBufferFlush Send any packets queued up for transmission on a port
// and HW queue.
//
// This causes an explicit flush of packets previously buffered via the
// rte_eth_tx_buffer() function. It returns the number of packets
// successfully sent to the NIC, and calls the error callback for any
// unsent packets. Unless explicitly set up otherwise, the default
// callback simply frees the unsent packets back to the owning mempool.
func (pid Port) TxBufferFlush(qid uint16, buf *TxBuffer) uint16 {
	return uint16(C.rte_eth_tx_buffer_flush(C.uint16_t(pid), C.uint16_t(qid),
		(*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf))))
}

// TxBuffer buffers a single packet for future transmission on a port
// and queue.
//
// This function takes a single mbuf/packet and buffers it for later
// transmission on the particular port and queue specified. Once the
// buffer is full of packets, an attempt will be made to transmit all
// the buffered packets. In case of error, where not all packets can
// be transmitted, a callback is called with the unsent packets as a
// parameter. If no callback is explicitly set up, the unsent packets
// are just freed back to the owning mempool. The function returns the
// number of packets actually sent i.e. 0 if no buffer flush occurred,
// otherwise the number of packets successfully flushed
func (pid Port) TxBuffer(qid uint16, buf *TxBuffer, m *mbuf.Mbuf) uint16 {
	return uint16(C.rte_eth_tx_buffer(C.uint16_t(pid), C.uint16_t(qid),
		(*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf)),
		(*C.struct_rte_mbuf)(unsafe.Pointer(m))))
}

// Mbufs returns a slice of packets contained in TxBuffer.
func (buf *TxBuffer) Mbufs() []*mbuf.Mbuf {
	var d []*mbuf.Mbuf
	b := (*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	sh.Data = uintptr(unsafe.Pointer(b)) + unsafe.Sizeof(*b)
	sh.Len = int(b.length)
	sh.Cap = int(b.size)
	return d
}

// Init initializes pre-allocated TxBuffer which must have enough
// memory to contain cnt mbufs.
func (buf *TxBuffer) Init(cnt int) {
	b := (*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf))
	C.rte_eth_tx_buffer_init(b, C.uint16_t(cnt))
}
