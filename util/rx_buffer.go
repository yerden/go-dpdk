package util

/*
#cgo pkg-config: libdpdk
#include <stdlib.h>

#include <rte_config.h>
#include <rte_ethdev.h>
#include "rx_buffer.h"
*/
import "C"

import (
	"reflect"
	"syscall"
	"unsafe"

	"github.com/google/gopacket"
	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mbuf"
)

func err(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

// RxBuffer implements gopacket.ZeroCopyPacketDataSource interface
// wrapping port and queue id.
type RxBuffer C.struct_rx_buffer

// NewRxBuffer allocates new RxBuffer from huge pages memory for
// specified queue id, socket NUMA node and containing up to size
// mbufs. If EAL failed to allocate memory it will panic.
func NewRxBuffer(pid ethdev.Port, qid uint16, socket int, size uint16) *RxBuffer {
	var p *C.struct_rx_buffer
	e := err(C.new_rx_buffer(C.int(socket), C.ushort(size), &p))
	if e != nil {
		panic(e)
	}

	p.pid = C.ushort(pid)
	p.qid = C.ushort(qid)
	return (*RxBuffer)(p)
}

// Free releases acquired huge pages memory.
func (buf *RxBuffer) Free() {
	C.rte_free(unsafe.Pointer(buf))
}

func (buf *RxBuffer) cursor() *mbuf.Mbuf {
	p := (*[1e6]*mbuf.Mbuf)(unsafe.Pointer(&buf.pkts[0]))
	return p[buf.n]
}

// Mbufs returns all mbufs retrieved by the ethdev API.
func (buf *RxBuffer) Mbufs() (ret []*mbuf.Mbuf) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	sh.Len = int(buf.length)
	sh.Cap = int(buf.size)
	sh.Data = uintptr(unsafe.Pointer(&buf.pkts[0]))
	return
}

// Recharge releases previously retrieved packets and retrieve new
// ones.
func (buf *RxBuffer) Recharge() int {
	return int(C.recharge_rx_buffer((*C.struct_rx_buffer)(buf)))
}

// ZeroCopyReadPacketData implements
// gopacket.ZeroCopyPacketDataSource interface.
//
// If RX queue does not yield any packets, syscall.EAGAIN error is
// returned.
//
// XXX: Please note that timestamp for CaptureInfo is not set by this
// call.
func (buf *RxBuffer) ZeroCopyReadPacketData() (data []byte, ci gopacket.CaptureInfo, err error) {
	if buf.n >= buf.length {
		buf.Recharge()
	}

	if buf.length == 0 {
		err = syscall.EAGAIN
		return
	}

	m := buf.cursor()
	data = m.Data()
	ci.Length = len(data)
	ci.CaptureLength = len(data)
	ci.InterfaceIndex = int(buf.pid)
	buf.n++
	return
}
