package util

/*
#cgo pkg-config: libdpdk

#define MBUF_ARRAY_USER_SIZE 64

#include <stdlib.h>
#include <rte_config.h>
#include <rte_ethdev.h>
#include "mbuf_array.h"

*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mbuf"
)

func geterr(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

// MbufArray is a wrapper of port and queue id.
type MbufArray C.struct_mbuf_array

// NewMbufArray allocates new MbufArray from huge pages memory for
// socket NUMA node and containing up to size mbufs. If EAL failed to
// allocate memory it will panic.
func NewMbufArray(socket int, size uint16) *MbufArray {
	var p *C.struct_mbuf_array
	if e := geterr(C.new_mbuf_array(C.int(socket), C.ushort(size), &p)); e != nil {
		panic(e)
	}
	return (*MbufArray)(p)
}

// Opaque returns slice of bytes pointing to user-defined data in MbufArray.
func (buf *MbufArray) Opaque() (opaque []byte) {
	return (*[unsafe.Sizeof(buf.opaque)]byte)(unsafe.Pointer(&opaque))[:]
}

// NewEthdevMbufArray allocates new MbufArray from huge pages memory
// for specified queue id, socket NUMA node and containing up to size
// mbufs. If EAL failed to allocate memory it will panic.
func NewEthdevMbufArray(pid ethdev.Port, qid uint16, socket int, size uint16) *MbufArray {
	p := NewMbufArray(socket, size)
	opaque := (*C.struct_ethdev_data)(unsafe.Pointer(&p.opaque))
	opaque.pid = C.ushort(pid)
	opaque.qid = C.ushort(qid)
	return p
}

// Free releases acquired huge pages memory.
func (buf *MbufArray) Free() {
	C.rte_free(unsafe.Pointer(buf))
}

func (buf *MbufArray) cursor() *mbuf.Mbuf {
	p := (*[1 << 31]*mbuf.Mbuf)(unsafe.Pointer(&buf.pkts[0]))
	return p[buf.n]
}

// Buffer returns all mbufs in buf.
func (buf *MbufArray) Buffer() (ret []*mbuf.Mbuf) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	sh.Len = int(buf.size)
	sh.Cap = int(buf.size)
	sh.Data = uintptr(unsafe.Pointer(&buf.pkts[0]))
	return
}

// Mbufs returns all mbufs retrieved by the ethdev API.
func (buf *MbufArray) Mbufs() (ret []*mbuf.Mbuf) {
	return buf.Buffer()[:buf.length]
}

// Recharge releases previously retrieved packets and retrieve new
// ones. Returns number of retrived packets.
func (buf *MbufArray) Recharge() int {
	return int(C.mbuf_array_ethdev_reload((*C.struct_mbuf_array)(buf)))
}
