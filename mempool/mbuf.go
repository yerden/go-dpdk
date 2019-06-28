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
	"syscall"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// CreateMbufPool creates mempool of mbufs. See Create for a list of
// options. Only differencies are described below.
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

	// check alignment
	if C.is_priv_size_aligned(C.uint16_t(conf.privDataSize)) == 0 {
		return nil, syscall.EINVAL
	}

	// calculate element size
	eltSize := C.uint(dataRoomSize) + C.sizeof_struct_rte_mbuf + conf.privDataSize
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	mbp_priv := &C.struct_rte_pktmbuf_pool_private{
		mbuf_data_room_size: C.uint16_t(dataRoomSize),
		mbuf_priv_size:      C.uint16_t(conf.privDataSize),
	}

	// create mempool
	mp := (*Mempool)(C.rte_mempool_create_empty(cname, C.uint(n), eltSize,
		conf.cacheSize, C.sizeof_struct_rte_pktmbuf_pool_private, conf.socket, 0))
	if mp == nil {
		return nil, common.Errno(nil)
	}

	if conf.opsName == nil {
		defaultOps := "ring_mp_mc"
		conf.opsName = &defaultOps
	}

	// set ops
	if err := mp.SetOpsByName(*conf.opsName, nil); err != nil {
		mp.Free()
		return nil, err
	}

	// init and populate pool
	C.rte_pktmbuf_pool_init((*C.struct_rte_mempool)(mp), unsafe.Pointer(mbp_priv))

	if _, err := mp.PopulateDefault(); err != nil {
		mp.Free()
		return nil, err
	}

	// initialize objects
	mp.ObjIterC(C.rte_pktmbuf_init, nil)
	return mp, nil
}
