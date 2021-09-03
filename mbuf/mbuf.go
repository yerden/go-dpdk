package mbuf

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_mbuf.h>

char *reset_and_append(struct rte_mbuf *mbuf, void *ptr, size_t len)
{
	rte_pktmbuf_reset(mbuf);
	char *data = rte_pktmbuf_append(mbuf, len);
	if (data == NULL)
		return NULL;
	rte_memcpy(data, ptr, len);
	return data;
}
struct rte_mbuf *alloc_reset_and_append(struct rte_mempool *mp, void *ptr, size_t len)
{
	struct rte_mbuf *mbuf;
	mbuf = rte_pktmbuf_alloc(mp);
	if (mbuf == NULL)
		return NULL;
	rte_pktmbuf_reset(mbuf);
	char *data = rte_pktmbuf_append(mbuf, len);
	if (data == NULL)
		return NULL;
	rte_memcpy(data, ptr, len);
	return mbuf;
}

enum {
	MBUF_RSS_OFF = offsetof(struct rte_mbuf, hash.rss),
};

*/
import "C"

import (
	"errors"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/mempool"
)

// ErrNullData is returned if NULL is returned by Cgo call.
var ErrNullData = errors.New("NULL response returned")

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
func (m *Mbuf) PktMbufFree() {
	C.rte_pktmbuf_free(mbuf(m))
}

// RawFree returns this mbuf into its originating mempool.
func (m *Mbuf) RawFree() {
	C.rte_mbuf_raw_free(mbuf(m))
}

// PktMbufClone clones the mbuf using supplied mempool as the buffer
// source. NOTE: NULL may return if allocation fails.
func (m *Mbuf) PktMbufClone(p *mempool.Mempool) {
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
	e := C.rte_pktmbuf_alloc_bulk(mp(p), mbufs(ms), C.uint(len(ms)))
	return syscall.Errno(e)
}

// PktMbufPrivSize get the application private size of mbufs
// stored in a pktmbuf_pool. The private size of mbuf is a zone
// located between the rte_mbuf structure and the data buffer
// where an application can store data associated to a packet.
func PktMbufPrivSize(p *mempool.Mempool) int {
	return (int)(C.rte_pktmbuf_priv_size(mp(p)))
}

// PktMbufAppend append the given data to an mbuf.
// Error may be returned if there is not enough tailroom
// space in the last segment of mbuf.
func (m *Mbuf) PktMbufAppend(data []byte) error {
	ptr := C.rte_pktmbuf_append(mbuf(m), C.uint16_t(len(data)))
	if ptr == nil {
		return ErrNullData
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&ptr))
	copy((*(*[]byte)(unsafe.Pointer(sh)))[:], data)
	return nil
}

// PktMbufReset reset the fields of a packet mbuf to their default values.
func (m *Mbuf) PktMbufReset() {
	C.rte_pktmbuf_reset(mbuf(m))
}

// Mempool return a pool from which mbuf was allocated.
func (m *Mbuf) Mempool() *mempool.Mempool {
	rteMbuf := mbuf(m)
	memp := rteMbuf.pool
	return (*mempool.Mempool)(unsafe.Pointer(memp))
}

// PrivSize return a size of private data area.
func (m *Mbuf) PrivSize() uint16 {
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

// PrivData sets ptr to point to mbuf's private area. Private data
// length is set to the priv_size field of an mbuf itself. Although
// this length may be 0 the private area may still be usable as
// HeadRoomSize is not 0.
//
// Feel free to edit the contents. A pointer to the headroom
// will be returned if the length of the private zone is 0.
func (m *Mbuf) PrivData(ptr *common.CStruct) {
	rteMbuf := mbuf(m)
	ptr.Init(unsafe.Add(unsafe.Pointer(m), unsafe.Sizeof(*m)), int(rteMbuf.priv_size))
}

// ResetAndAppend reset the fields of a mbuf to their default values
// and append the given data to an mbuf. Error may be returned
// if there is not enough tailroom space in the last segment of mbuf.
// Len is the amount of data to append (in bytes).
func (m *Mbuf) ResetAndAppend(data *common.CStruct) error {
	ptr := C.reset_and_append(mbuf(m), data.Ptr, C.size_t(data.Len))
	if ptr == nil {
		return ErrNullData
	}
	return nil
}

// AllocResetAndAppend allocates an uninitialized mbuf from mempool p,
// resets the fields of the mbuf to default values and appends the
// given data to an mbuf.
//
// Note that NULL may be returned if allocation failed or if there is
// not enough tailroom space in the last segment of mbuf.  p is the
// mempool from which the mbuf is allocated. Data is C array
// representation of data to add.
func AllocResetAndAppend(p *mempool.Mempool, data *common.CStruct) *Mbuf {
	m := C.alloc_reset_and_append(mp(p), data.Ptr, C.size_t(data.Len))
	return (*Mbuf)(unsafe.Pointer(m))
}

// HeadRoomSize returns the value of the data_off field,
// which must be equal to the size of the headroom in concrete mbuf.
func (m *Mbuf) HeadRoomSize() uint16 {
	rteMbuf := mbuf(m)
	return uint16(rteMbuf.data_off)
}

// TailRoomSize returns available length that can be appended to mbuf.
func (m *Mbuf) TailRoomSize() uint16 {
	rteMbuf := mbuf(m)
	return uint16(rteMbuf.buf_len - rteMbuf.data_off - rteMbuf.data_len)
}

// BufLen represents DataRoomSize that was initialized in
// mempool.CreateMbufPool.
//
// NOTE: Max available data length that mbuf can hold is BufLen -
// HeadRoomSize.
func (m *Mbuf) BufLen() uint16 {
	rteMbuf := mbuf(m)
	return uint16(rteMbuf.buf_len)
}

// PktMbufHeadRoomSize represents RTE_PKTMBUF_HEADROOM size in
// concrete mbuf.
//
// NOTE: This implies Cgo call and is used for testing purposes only.
// Use HeadRoomSize instead.
func (m *Mbuf) PktMbufHeadRoomSize() uint16 {
	return uint16(C.rte_pktmbuf_headroom(mbuf(m)))
}

// PktMbufTailRoomSize represents RTE_PKTMBUF_TAILROOM which is
// available length that can be appended to mbuf.
//
// NOTE: This implies Cgo call and is used for testing purposes only.
// Use TailRoomSize instead.
func (m *Mbuf) PktMbufTailRoomSize() uint16 {
	return uint16(C.rte_pktmbuf_tailroom(mbuf(m)))
}

// HashRss returns hash.rss field of an mbuf.
func (m *Mbuf) HashRss() uint32 {
	p := unsafe.Pointer(m)
	p = unsafe.Pointer(uintptr(p) + C.MBUF_RSS_OFF)
	return *(*uint32)(p)
}
