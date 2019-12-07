package mbuf

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_mbuf.h>
*/
import "C"

import (
	// "syscall"
	"reflect"
	"unsafe"

	"github.com/yerden/go-dpdk/mempool"
)

// Mbuf contains a packet.
type Mbuf C.struct_rte_mbuf

func mbuf(m *Mbuf) *C.struct_rte_mbuf {
	return (*C.struct_rte_mbuf)(unsafe.Pointer(m))
}

func mp(m *mempool.Mempool) *C.struct_rte_mempool {
	return (*C.struct_rte_mempool)(unsafe.Pointer(m))
}

func PktMbufFree(m *Mbuf) {
	C.rte_pktmbuf_free(mbuf(m))
}

func RawFree(m *Mbuf) {
	C.rte_mbuf_raw_free(mbuf(m))
}

func PktMbufClone(m *Mbuf, p *mempool.Mempool) {
	C.rte_pktmbuf_clone(mbuf(m), mp(p))
}

func (m *Mbuf) Data() []byte {
	var d []byte
	buf := mbuf(m)
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	sh.Data = uintptr(buf.buf_addr) + uintptr(buf.data_off)
	sh.Len = int(buf.data_len)
	sh.Cap = int(buf.data_len)
	return d
}
