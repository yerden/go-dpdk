package mbuf

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_mbuf.h>
*/
import "C"

import (
	"errors"
	// "syscall"
	"reflect"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/mempool"
)

var TooLargeData = errors.New("data size can't be larger then priv_size")

// Mbuf contains a packet.
type Mbuf C.struct_rte_mbuf

func mbuf(m *Mbuf) *C.struct_rte_mbuf {
	return (*C.struct_rte_mbuf)(unsafe.Pointer(m))
}

func mbufs(ms []*Mbuf) **C.struct_rte_mbuf {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&ms))
	return (**C.struct_rte_mbuf)(unsafe.Pointer(sh.Data))
}

func mp(m *mempool.Mempool) *C.struct_rte_mempool {
	return (*C.struct_rte_mempool)(unsafe.Pointer(m))
}

// PktMbufFree returns this mbuf into its originating mempool along
// with all its segments.
func PktMbufFree(m *Mbuf) {
	C.rte_pktmbuf_free(mbuf(m))
}

// RawFree returns this mbuf into its originating mempool.
func RawFree(m *Mbuf) {
	C.rte_mbuf_raw_free(mbuf(m))
}

// PktMbufClone clones the mbuf using supplied mempool as the buffer
// source.
func PktMbufClone(m *Mbuf, p *mempool.Mempool) {
	C.rte_pktmbuf_clone(mbuf(m), mp(p))
}

// PktMbufAlloc allocate an uninitialized mbuf from mempool p.
// Note that NULL may be returned if allocation failed.
func PktMbufAlloc(p *mempool.Mempool) *Mbuf {
	m := C.rte_pktmbuf_alloc(mp(p))
	return (*Mbuf)(m)
}

// PktMbufAllocBulk allocate a bulk of mbufs.
func PktMbufAllocBulk(p *mempool.Mempool, ms []*Mbuf) error {
	return common.Err(C.rte_pktmbuf_alloc_bulk(mp(p), mbufs(ms), C.uint(len(ms))))
}

// PktMbufPrivSize get the application private size of mbufs
// stored in a pktmbuf_pool. The private size of mbuf is a zone
// located between the rte_mbuf structure and the data buffer
// where an application can store data associated to a packet.
func PktMbufPrivSize(p *mempool.Mempool) int {
	return (int)(C.rte_pktmbuf_priv_size(mp(p)))
}

// GetPrivData return data stored in private data area
// embedded in the given mbuf. Note that no check is made
// to ensure that a private data area actually exists in the supplied mbuf.
func GetPrivData(m *Mbuf) []byte {
	a := C.rte_mbuf_to_priv(mbuf(m))

	var priv []byte
	cmb := (*C.struct_rte_mbuf)(m)
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&priv))
	sh.Data = uintptr(a)
	if sh.Data != 0 {
		sh.Len = int(cmb.priv_size)
		sh.Cap = sh.Len
	}
	return priv
}

// PutToPriv append the given data to a private data area.
// Note that the data size cannot be larger than size.
// Note that no check is made to ensure that a private data area
// actually exists in the supplied mbuf.
func PutToPriv(m *Mbuf, data []byte) error {
	cmb := (*C.struct_rte_mbuf)(m)
	if len(data) > int(cmb.priv_size) {
		return TooLargeData
	}

	a := C.rte_mbuf_to_priv(mbuf(m))

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	copy((*(*[]byte)(unsafe.Pointer(sh)))[:], data)
	return nil
}

// PktMbufAppend append the given data to an mbuf.
func PktMbufAppend(m *Mbuf, data []byte) {
	a := C.rte_pktmbuf_append(mbuf(m), C.uint16_t(len(data)))

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	copy((*(*[]byte)(unsafe.Pointer(sh)))[:], data)
}

// PktMbufReset reset the fields of a packet mbuf to their default values.
func PktMbufReset(m *Mbuf) {
	C.rte_pktmbuf_reset(mbuf(m))
}

// GetPool return a pool from which mbuf was allocated.
func (m *Mbuf) GetPool() *mempool.Mempool {
	rteMbuf := mbuf(m)
	memp := rteMbuf.pool
	return (*mempool.Mempool)(unsafe.Pointer(memp))
}

// Data returns contained packet.
func (m *Mbuf) Data() []byte {
	var d []byte
	buf := mbuf(m)
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	sh.Data = uintptr(buf.buf_addr) + uintptr(buf.data_off)
	sh.Len = int(buf.data_len)
	sh.Cap = int(buf.data_len)
	return d
}
