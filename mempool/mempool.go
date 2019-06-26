/*
Package mempool wraps RTE mempool library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package mempool

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_errno.h>
#include <rte_mempool.h>

extern void goMempoolObjCbFunc(struct rte_mempool *mp, void *opaque, void *obj, unsigned idx);

static int errget() {
	return rte_errno;
}
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

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

// Free destroys mempool.
func (mp *Mempool) Free() {
	C.rte_mempool_free((*C.struct_rte_mempool)(mp))
}

// MempoolObjCb is an object action for mempool iteration.
type MempoolObjCb func(unsafe.Pointer)

var (
	mpCb = common.NewRegistryArray()
)

//export goMempoolObjCbFunc
func goMempoolObjCbFunc(mp *C.struct_rte_mempool, opaque, obj unsafe.Pointer, obj_idx C.uint) {
	cb := *(*common.ObjectID)(opaque)
	fn := mpCb.Read(cb).(MempoolObjCb)
	fn(obj)
}

// ObjIter iterates mempool objects.
func (mp *Mempool) ObjIter(fn MempoolObjCb) uint32 {
	cb := mpCb.Create(fn)
	defer mpCb.Delete(cb)

	objCb := (*C.rte_mempool_obj_cb_t)(C.goMempoolObjCbFunc)
	cmp := (*C.struct_rte_mempool)(mp)
	return uint32(C.rte_mempool_obj_iter(cmp, objCb, unsafe.Pointer(&cb)))
}

// ObjIterC iterates mempool objects. fn is a pointer to
// rte_mempool_obj_cb_t, opaque is its opaque argument.
func (mp *Mempool) ObjIterC(fn, opaque unsafe.Pointer) uint32 {
	objCb := (*C.rte_mempool_obj_cb_t)(fn)
	cmp := (*C.struct_rte_mempool)(mp)
	return uint32(C.rte_mempool_obj_iter(cmp, objCb, opaque))
}
