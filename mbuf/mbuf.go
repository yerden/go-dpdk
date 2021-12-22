package mbuf

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_mbuf.h>

char *reset_and_append(struct rte_mbuf *mbuf, void *ptr, size_t len)
{
	rte_pktmbuf_reset(mbuf);
	char *data = rte_pktmbuf_append(mbuf, len);
	rte_memcpy(data, ptr, len);
	return data;
}

struct rte_mbuf *alloc_reset_and_append(struct rte_mempool *mp, void *ptr, size_t len)
{
	struct rte_mbuf *mbuf;
	mbuf = rte_pktmbuf_alloc(mp);
	rte_pktmbuf_reset(mbuf);
	char *data = rte_pktmbuf_append(mbuf, len);
	if (data == NULL)
		return NULL;
	rte_memcpy(data, ptr, len);

	return mbuf;
}


static inline void *get_go_struct(struct rte_mbuf *mbuf, void *ptr, size_t len)
{
	char *data_offset;
	data_offset = rte_pktmbuf_mtod(mbuf, char *);
	rte_memcpy(ptr, data_offset, len);
}
*/
import "C"

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/mempool"
)

var TooLargeData = errors.New("data size can't be larger then priv_size")
var NullData = errors.New("NULL response returned")

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

// GetPrivSize return a size of private data area.
func (m *Mbuf) GetPrivSize() uint16 {
	rteMbuf := mbuf(m)
	s := rteMbuf.priv_size
	return uint16(s)
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

// GetPrivData returns the content of the mbufs private area.
// The private area has a certain length,
// which is set when creating the mbufpool, do not try to increase it.
// Feel free to edit the contents.
func (m *Mbuf) GetPrivData() *common.CStruct {
	rteMbuf := mbuf(m)
	p := unsafe.Add(unsafe.Pointer(m), unsafe.Sizeof(*m))
	return &common.CStruct{Ptr: p, Len: int(rteMbuf.priv_size)}
}

// ResetAndAppend reset the fields of a mbuf to their default values
// and append the given data to an mbuf. Error may be returned
// if there is not enough tailroom space in the last segment of mbuf.
// Len is the amount of data to append (in bytes).
func (m *Mbuf) ResetAndAppend(data *common.CStruct) error {
	ptr := C.reset_and_append(mbuf(m), data.Ptr, C.size_t(data.Len))
	if ptr == nil {
		return NullData
	}
	return nil
}

// AllocResetAndAppend allocate an uninitialized mbuf from mempool p,
// reset the fields of a mbuf to their default values
// and append the given data to an mbuf.
// Note that NULL may be returned if allocation failed or if
// there is not enough tailroom space in the last segment of mbuf.
// p is the mempool from which the mbuf is allocated.
// Len is the amount of data to append (in bytes).
func AllocResetAndAppend(p *mempool.Mempool, data *common.CStruct) *Mbuf {
	m := C.alloc_reset_and_append(mp(p), data.Ptr, C.size_t(data.Len))
	return (*Mbuf)(unsafe.Pointer(m))
}

// GetGoStruct represent bytes from data_room of given mbuf in go representation.
// Accepts a pointer to the go type or structure
// to which the data should be written, and the amount of data.
func (m *Mbuf) GetGoStruct(data *common.CStruct) {
	C.get_go_struct(mbuf(m), data.Ptr, C.size_t(data.Len))
}
