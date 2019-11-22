package mempool

/*
#include <rte_config.h>
#include <rte_mempool.h>
*/
import "C"

import (
	// "reflect"
	// "unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Cache is a structure that stores a per-core object cache. This can
// be used by non-EAL threads to enable caching when they interact
// with a mempool.
type Cache C.struct_rte_mempool_cache

// CreateCache creates a user-owned mempool cache.
//
// You may specify OptCacheSize and OptSocket options only.
func CreateCache(opts ...Option) (*Cache, error) {
	conf := &mpConf{socket: C.SOCKET_ID_ANY}
	for _, o := range opts {
		o.f(conf)
	}

	mpc := C.rte_mempool_cache_create(conf.cacheSize, conf.socket)
	if mpc == nil {
		return nil, common.Errno(nil)
	}
	return (*Cache)(mpc), nil
}

// Flush a user-owned mempool cache to the specified mempool.
func (mpc *Cache) Flush(mp *Mempool) {
	C.rte_mempool_cache_flush((*C.struct_rte_mempool_cache)(mpc), (*C.struct_rte_mempool)(mp))
}

// Free a user-owned mempool cache.
func (mpc *Cache) Free() {
	C.rte_mempool_cache_free((*C.struct_rte_mempool_cache)(mpc))
}
