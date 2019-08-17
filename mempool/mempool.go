/*
Package mempool wraps RTE mempool library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package mempool

/*
#include <rte_config.h>
#include <rte_mempool.h>

extern void goObjectCb(struct rte_mempool *mp, void *opaque, void *obj, unsigned idx);
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Mempool represents RTE mempool.
type Mempool C.struct_rte_mempool

// mempool configuration
type mpConf struct {
	cacheSize    C.uint
	privDataSize C.uint
	socket       C.int
	flags        C.uint

	// ops
	opsName       *string
	opsPoolConfig unsafe.Pointer
}

// OptOpsName specifies mempool's name. If option is used in Create
// this name is used when creating memzone. This name maybe then used
// to lookup the mempool.

// OptOpsName sets the ops of a mempool. Currently implemented in DPDK
// are: 'ring_mp_mc', 'ring_sp_mc', 'ring_mp_sc', 'ring_sp_sc',
// 'stack', 'lf_stack'.
func OptOpsName(name string) Option {
	return Option{func(conf *mpConf) {
		conf.opsName = &name
	}}
}

// OptOpsPoolConfig specifies opaque argument for mempool ops.
// Currently it's not in use by any implemented ops.
func OptOpsPoolConfig(p unsafe.Pointer) Option {
	return Option{func(conf *mpConf) {
		conf.opsPoolConfig = p
	}}
}

// Option is used to configure mempool at creation time.
type Option struct {
	f func(*mpConf)
}

// OptCacheSize specifies cache size. If zero, the rte_mempool library
// will try to limit the accesses to the common lockless pool, by
// maintaining a per-lcore object cache. This argument must be lower
// or equal to CONFIG_RTE_MEMPOOL_CACHE_MAX_SIZE and n / 1.5 where n
// is number of elements. It is advised to choose cache_size to have
// "n modulo cache_size == 0": if this is not the case, some elements
// will always stay in the pool and will never be used. The access to
// the per-lcore table is of course faster than the
// multi-producer/consumer pool. The cache can be disabled if the
// cache_size argument is set to 0; it can be useful to avoid losing
// objects in cache.
func OptCacheSize(size uint32) Option {
	return Option{func(conf *mpConf) {
		conf.cacheSize = C.uint(size)
	}}
}

// OptPrivateDataSize specifies size of the private data appended
// after the mempool structure. This is useful for storing some
// private data after the mempool structure, as is done for
// rte_mbuf_pool for example.
func OptPrivateDataSize(size uint32) Option {
	return Option{func(conf *mpConf) {
		conf.privDataSize = C.uint(size)
	}}
}

// OptSocket specifies socket identifier in the case of NUMA. The
// value can be SOCKET_ID_ANY if there is no NUMA constraint for the
// reserved zone.
func OptSocket(socket int) Option {
	return Option{func(conf *mpConf) {
		conf.socket = C.int(socket)
	}}
}

// OptFlag specifies various flags to use when creating mempool.
func OptFlag(flag uint) Option {
	return Option{func(conf *mpConf) {
		conf.flags |= C.uint(flag)
	}}
}

// Various non-parameterized options for mempool creation.
const (
	// By default, objects addresses are spread between channels in
	// RAM: the pool allocator will add padding between objects
	// depending on the hardware configuration. See Memory alignment
	// constraints for details. If this flag is set, the allocator
	// will just align them to a cache line.
	NoSpread uint = C.MEMPOOL_F_NO_SPREAD
	// By default, the returned objects are cache-aligned. This flag
	// removes this constraint, and no padding will be present between
	// objects. This flag implies NoSpread.
	NoCacheAlign = C.MEMPOOL_F_NO_CACHE_ALIGN
	// If this flag is set, the default behavior when using
	// rte_mempool_put() or rte_mempool_put_bulk() is
	// "single-producer". Otherwise, it is "multi-producers".
	SPPut = C.MEMPOOL_F_SP_PUT
	// If this flag is set, the default behavior when using
	// rte_mempool_get() or rte_mempool_get_bulk() is
	// "single-consumer". Otherwise, it is "multi-consumers".
	SCGet = C.MEMPOOL_F_SC_GET
	// If set, allocated objects won't necessarily be contiguous in IO
	// memory.
	NoPhysContig = C.MEMPOOL_F_NO_PHYS_CONTIG
)

// Option shortcuts.
var (
	OptNoSpread     = OptFlag(NoSpread)
	OptNoCacheAlign = OptFlag(NoCacheAlign)
	OptSPPut        = OptFlag(SPPut)
	OptSCGet        = OptFlag(SCGet)
	OptNoPhysContig = OptFlag(NoPhysContig)
)

func cGoString(s string) *C.char {
	a := append([]byte(s), 0)
	return (*C.char)(unsafe.Pointer(&a[0]))
}

// CreateEmpty creates new empty mempool. The mempool is allocated and
// initialized, but it is not populated: no memory is allocated for
// the mempool elements. The user has to call PopulateDefault() or
// other API to add memory chunks to the pool. Once populated, the
// user may also want to initialize each object with ObjIter/ObjIterC.
func CreateEmpty(name string, n, eltsize uint32, opts ...Option) (*Mempool, error) {
	conf := &mpConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	cname := cGoString(name)
	mp := (*Mempool)(C.rte_mempool_create_empty(cname, C.uint(n), C.uint(eltsize),
		conf.cacheSize, conf.privDataSize, conf.socket, conf.flags))

	if mp == nil {
		return nil, common.Errno(nil)
	}

	if conf.opsName != nil {
		err := mp.SetOpsByName(*conf.opsName, conf.opsPoolConfig)
		if err != nil {
			mp.Free()
			return nil, err
		}
	}

	return (*Mempool)(mp), nil
}

// SetOpsByName sets the ops of a mempool.  This can only be done on a
// mempool that is not populated, i.e. just after a call to
// CreateEmpty().
func (mp *Mempool) SetOpsByName(name string, poolConfig unsafe.Pointer) error {
	cName := cGoString(name)
	cmp := (*C.struct_rte_mempool)(mp)
	return common.Errno(C.rte_mempool_set_ops_byname(cmp, cName, poolConfig))
}

// PopulateDefault adds memory for objects in the pool at init. This
// is the default function used by rte_mempool_create() to populate
// the mempool. It adds memory allocated using rte_memzone_reserve().
func (mp *Mempool) PopulateDefault() (int, error) {
	rc := C.rte_mempool_populate_default((*C.struct_rte_mempool)(mp))
	return common.IntOrErr(rc)
}

// Free the mempool. Unlink the mempool from global list, free the
// memory chunks, and all memory referenced by the mempool. The
// objects must not be used by other cores as they will be freed.
func (mp *Mempool) Free() {
	C.rte_mempool_free((*C.struct_rte_mempool)(mp))
}

// ObjectFunc is an object action for mempool iteration.
type ObjectFunc func(unsafe.Pointer)

var (
	mpCb = common.NewRegistryArray()
)

//export goObjectCb
func goObjectCb(mp *C.struct_rte_mempool, opaque, obj unsafe.Pointer, obj_idx C.uint) {
	cb := *(*common.ObjectID)(opaque)
	fn := mpCb.Read(cb).(ObjectFunc)
	fn(obj)
}

// ObjIter calls a function for each mempool element. Iterate across
// all objects attached to a rte_mempool and call the callback
// function on it.
func (mp *Mempool) ObjIter(fn ObjectFunc) uint32 {
	cb := mpCb.Create(fn)
	defer mpCb.Delete(cb)

	objCb := (*C.rte_mempool_obj_cb_t)(C.goObjectCb)
	cmp := (*C.struct_rte_mempool)(mp)
	return uint32(C.rte_mempool_obj_iter(cmp, objCb, unsafe.Pointer(&cb)))
}

// ObjIterC calls a function for each mempool element. Iterate across
// all objects attached to a rte_mempool and call the callback
// function on it.
//
// Callback function should conform to rte_mempool_obj_cb_t type.
func (mp *Mempool) ObjIterC(fn, opaque unsafe.Pointer) uint32 {
	objCb := (*C.rte_mempool_obj_cb_t)(fn)
	cmp := (*C.struct_rte_mempool)(mp)
	return uint32(C.rte_mempool_obj_iter(cmp, objCb, opaque))
}
