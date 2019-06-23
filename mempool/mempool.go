/*
Package mempool wraps RTE mempool library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package mempool

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_errno.h>
#include <rte_memory.h>
#include <rte_mempool.h>

extern void goMempoolObjCbFunc(struct rte_mempool *mp, void *opaque, void *obj, unsigned idx);

static uint32_t mp_obj_iter(
	struct rte_mempool *mp,
	rte_mempool_obj_cb_t *obj_cb,
	uintptr_t obj_cb_arg)
{
	return rte_mempool_obj_iter(mp, obj_cb, (void *)obj_cb_arg);
}

static void mp_obj_cb(
	uintptr_t fnptr,
	struct rte_mempool *mp,
	uintptr_t opaque, uintptr_t obj,
	unsigned idx)
{
	rte_mempool_obj_cb_t *fn = (typeof(fn)) fnptr;
	fn(mp, (void *)opaque, (void *)obj, idx);
}

static int errget() {
	return rte_errno;
}
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

func cptr(i interface{}) C.uintptr_t {
	return C.uintptr_t(common.Uintptr(i))
}

func errno(n C.int) error {
	return common.Errno(int(n))
}

type Mempool C.struct_rte_mempool

type mpConf struct {
	cacheSize    C.uint
	privDataSize C.uint
	socket       C.int
	flags        C.uint

	// ops
	opsName       *string
	opsPoolConfig unsafe.Pointer
}

func OptOpsName(name string) Option {
	return Option{func(conf *mpConf) {
		conf.opsName = &name
	}}
}

func OptOpsPoolConfig(p unsafe.Pointer) Option {
	return Option{func(conf *mpConf) {
		conf.opsPoolConfig = p
	}}
}

type Option struct {
	f func(*mpConf)
}

func OptCacheSize(size uint32) Option {
	return Option{func(conf *mpConf) {
		conf.cacheSize = C.uint(size)
	}}
}

func OptPrivateDataSize(size uint32) Option {
	return Option{func(conf *mpConf) {
		conf.privDataSize = C.uint(size)
	}}
}

func OptSocket(socket int) Option {
	return Option{func(conf *mpConf) {
		conf.socket = C.int(socket)
	}}
}

func optFlag(flag C.uint) Option {
	return Option{func(conf *mpConf) {
		conf.flags |= flag
	}}
}

var (
	OptNoSpread     = optFlag(C.MEMPOOL_F_NO_SPREAD)
	OptNoCacheAlign = optFlag(C.MEMPOOL_F_NO_CACHE_ALIGN)
	OptSPPut        = optFlag(C.MEMPOOL_F_SP_PUT)
	OptSCGet        = optFlag(C.MEMPOOL_F_SC_GET)
	OptNoPhysContig = optFlag(C.MEMPOOL_F_NO_PHYS_CONTIG)
)

func CreateEmpty(name string, n, eltsize uint32, opts ...Option) (*Mempool, error) {
	conf := &mpConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	mp := (*Mempool)(C.rte_mempool_create_empty(cname, C.uint(n), C.uint(eltsize),
		conf.cacheSize, conf.privDataSize, conf.socket, conf.flags))

	if mp == nil {
		return nil, errno(C.errget())
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

func (mp *Mempool) SetOpsByName(name string, poolConfig unsafe.Pointer) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return errno(C.rte_mempool_set_ops_byname((*C.struct_rte_mempool)(mp), cName, poolConfig))
}

func (mp *Mempool) PopulateDefault() (int, error) {
	rc := C.rte_mempool_populate_default((*C.struct_rte_mempool)(mp))
	return common.IntOrErr(int(rc))
}

func (mp *Mempool) Free() {
	C.rte_mempool_free((*C.struct_rte_mempool)(mp))
}

type MempoolObjCb func(unsafe.Pointer)

type mpObjIterCtx struct {
	fn MempoolObjCb
}

//export goMempoolObjCbFunc
func goMempoolObjCbFunc(mp *C.struct_rte_mempool, opaque, obj unsafe.Pointer, obj_idx C.uint) {
	ctx := (*mpObjIterCtx)(opaque)
	ctx.fn(obj)
}

func (mp *Mempool) ObjIter(fn MempoolObjCb) uint32 {
	// avoid GC
	ctx := &mpObjIterCtx{fn}
	defer runtime.KeepAlive(ctx)

	objCb := (*C.rte_mempool_obj_cb_t)(C.goMempoolObjCbFunc)
	return uint32(C.mp_obj_iter((*C.struct_rte_mempool)(mp), objCb, cptr(ctx)))
}

func (mp *Mempool) MakeObjCb(fn, opaque unsafe.Pointer) MempoolObjCb {
	idx := uint32(0)
	return func(obj unsafe.Pointer) {
		C.mp_obj_cb(cptr(fn), (*C.struct_rte_mempool)(mp),
			cptr(opaque), cptr(obj), C.uint(idx))
		idx++
	}
}
