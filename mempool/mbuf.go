package mempool

/*
#include <stdint.h>

#include <rte_errno.h>
#include <rte_config.h>
#include <rte_mbuf.h>

static int errget() {
	return rte_errno;
}

static int is_priv_size_aligned(uint16_t priv_size) {
	return RTE_ALIGN(priv_size, RTE_MBUF_PRIV_ALIGN) == priv_size;
}
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

func (mp *Mempool) MbufInitCb() MempoolObjCb {
	return mp.MakeObjCb(C.rte_pktmbuf_init, nil)
}

func CreateMbufPool(name string, n uint32, dataRoomSize uint16, opts ...Option) (*Mempool, error) {
	conf := &mpConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	// check alignment
	if C.is_priv_size_aligned(C.uint16_t(conf.privDataSize)) == 0 {
		return nil, common.Errno(C.EINVAL)
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
		return nil, errno(C.errget())
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
	mp.ObjIter(mp.MakeObjCb(C.rte_pktmbuf_init, nil))
	return mp, nil
}
