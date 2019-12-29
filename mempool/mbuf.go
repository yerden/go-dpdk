package mempool

/*
#include <stdint.h>

#include <rte_config.h>
#include <rte_mbuf.h>

static int is_priv_size_aligned(uint16_t priv_size) {
	return RTE_ALIGN(priv_size, RTE_MBUF_PRIV_ALIGN) == priv_size;
}
*/
import "C"

import (
	"unsafe"
)

// CreateMbufPool creates mempool of mbufs. See CreateEmpty options
// for a list of options. Only differences are described below.
//
// dataRoomSize specifies the maximum size of data buffer in each
// mbuf, including RTE_PKTMBUF_HEADROOM.
//
// OptPrivateDataSize semantics is different here. It specifies the
// size of private application data between the rte_mbuf structure and
// the data buffer.  This value must be aligned to
// RTE_MBUF_PRIV_ALIGN.
//
// The created mempool is already populated and its objects are
// initialized with rte_pktmbuf_init.
func CreateMbufPool(name string, n uint32, dataRoomSize uint16, opts ...Option) (*Mempool, error) {
	conf := &mpConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cOps *C.char
	if conf.opsName != nil {
		cOps = C.CString(*conf.opsName)
		defer C.free(unsafe.Pointer(cOps))
	}
	mp := (*Mempool)(C.rte_pktmbuf_pool_create_by_ops(cname, C.uint(n),
		conf.cacheSize, C.ushort(conf.privDataSize), C.uint16_t(dataRoomSize),
		conf.socket, cOps))

	if mp == nil {
		return nil, err()
	}

	return mp, nil
}
