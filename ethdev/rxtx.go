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

type TxBuffer C.struct_rte_eth_dev_tx_buffer

func fracCeilRound(x int, y int) int {
	n := x / y
	if y*n != x {
		n++
	}
	return n
}

func NewTxBuffer(cnt int) *TxBuffer {
	amount := TxBufferSize(cnt)
	ptrSize := unsafe.Sizeof(&mbuf.Mbuf{})
	data := make([]*mbuf.Mbuf, fracCeilRound(int(amount), int(ptrSize)))
	buf := (*TxBuffer)(unsafe.Pointer(&data[0]))
	buf.Init(cnt)
	return buf
}

func TxBufferSize(mbufCnt int) uintptr {
	return unsafe.Sizeof(TxBuffer{}) + uintptr(mbufCnt)*unsafe.Sizeof(&mbuf.Mbuf{})
}

func (pid Port) RxBurst(qid uint16, pkts []*mbuf.Mbuf) uint16 {
	return uint16(C.rte_eth_rx_burst(C.uint16_t(pid), C.uint16_t(qid),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint16_t(len(pkts))))
}

func (pid Port) TxBurst(qid uint16, pkts []*mbuf.Mbuf) uint16 {
	return uint16(C.rte_eth_tx_burst(C.uint16_t(pid), C.uint16_t(qid),
		(**C.struct_rte_mbuf)(unsafe.Pointer(&pkts[0])), C.uint16_t(len(pkts))))
}

func (pid Port) TxBufferFlush(qid uint16, buf *TxBuffer) uint16 {
	return uint16(C.rte_eth_tx_buffer_flush(C.uint16_t(pid), C.uint16_t(qid),
		(*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf))))
}

func (pid Port) TxBuffer(qid uint16, buf *TxBuffer, m *mbuf.Mbuf) uint16 {
	return uint16(C.rte_eth_tx_buffer(C.uint16_t(pid), C.uint16_t(qid),
		(*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf)),
		(*C.struct_rte_mbuf)(unsafe.Pointer(m))))
}

func (buf *TxBuffer) Mbufs() []*mbuf.Mbuf {
	var d []*mbuf.Mbuf
	b := (*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf))
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	sh.Data = uintptr(unsafe.Pointer(b)) + unsafe.Sizeof(*b)
	sh.Len = int(b.length)
	sh.Cap = int(b.size)
	return d
}

func (buf *TxBuffer) Init(cnt int) {
	b := (*C.struct_rte_eth_dev_tx_buffer)(unsafe.Pointer(buf))
	C.rte_eth_tx_buffer_init(b, C.uint16_t(cnt))
}
