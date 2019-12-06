/*
Package memzone wraps RTE memzone library.

The goal of the memzone allocator is to reserve contiguous portions of
physical memory. These zones are identified by a name.

The memzone descriptors are shared by all partitions and are located
in a known place of physical memory. This zone is accessed using
rte_eal_get_configuration(). The lookup (by name) of a memory zone can
be done in any partition and returns the same physical address.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package memzone

/*
#include <rte_config.h>
#include <rte_memzone.h>

typedef void memzone_cb_t(const struct rte_memzone *, void *);
extern void mzCb(struct rte_memzone *, void *);
static void memzone_walk(void *f, void *arg) {
	rte_memzone_walk((memzone_cb_t *)f, arg);
}
static void *memzone_addr(const struct rte_memzone *mz) {
	return mz->addr;
}
*/
import "C"

import (
	"reflect"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Memzone is a contiguous portion of physical memory.
// These zones are identified by a name.
type Memzone C.struct_rte_memzone

// These parameters are used to request memzones to be taken from
// specifically sized hugepages.
const (
	Page2Mb   uint = C.RTE_MEMZONE_2MB
	Page1Gb        = C.RTE_MEMZONE_1GB
	Page16Mb       = C.RTE_MEMZONE_16MB
	Page16Gb       = C.RTE_MEMZONE_16GB
	Page256Kb      = C.RTE_MEMZONE_256KB
	Page256Mb      = C.RTE_MEMZONE_256MB
	Page512Mb      = C.RTE_MEMZONE_512MB
	Page4Gb        = C.RTE_MEMZONE_4GB

	// Allow alternative page size to be used if the requested page
	// size is unavailable. If this flag is not set, the function will
	// return error on an unavailable size request.
	PageSizeHintOnly = C.RTE_MEMZONE_SIZE_HINT_ONLY

	// Ensure reserved memzone is IOVA-contiguous. This option should
	// be used when allocating memory intended for hardware rings etc.
	PageIovaContig = C.RTE_MEMZONE_IOVA_CONTIG
)

type conf struct {
	cname   *C.char
	size    C.size_t
	socket  C.int
	flags   C.uint
	aligned *C.uint
	bound   *C.uint
}

// Option may be specified during Reserve for a new memzone.
type Option struct {
	f func(*conf)
}

func err(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

// OptSocket specifies the socket id where the memzone would be
// created.
func OptSocket(socket uint) Option {
	return Option{func(rc *conf) {
		rc.socket = C.int(socket)
	}}
}

// OptFlag add one of permitted flags for the memzone creation.
func OptFlag(flag uint) Option {
	return Option{func(rc *conf) {
		rc.flags |= C.uint(flag)
	}}
}

// OptAligned requests alignment for resulting memzone.
// Must be a power of 2.
func OptAligned(align uint) Option {
	return Option{func(rc *conf) {
		rc.aligned = new(C.uint)
		*rc.aligned = C.uint(align)
	}}
}

// OptBounded requests boundary for resulting memzone. Must be a power
// of 2 or zero. Zero value implies no boundary condition.
func OptBounded(bound uint) Option {
	return Option{func(rc *conf) {
		rc.bound = new(C.uint)
		*rc.bound = C.uint(bound)
	}}
}

func cGoString(s string) *C.char {
	a := append([]byte(s), 0)
	return (*C.char)(unsafe.Pointer(&a[0]))
}

func makeOpts(name string, size uintptr, opts []Option) *conf {
	rc := &conf{
		socket: C.SOCKET_ID_ANY,
		size:   C.size_t(size),
		cname:  cGoString(name)}
	for i := range opts {
		opts[i].f(rc)
	}
	return rc
}

func doReserve(rc *conf) *C.struct_rte_memzone {
	switch {
	case rc.bound == nil && rc.aligned == nil:
		return C.rte_memzone_reserve(rc.cname, rc.size, rc.socket,
			rc.flags)
	case rc.bound == nil:
		return C.rte_memzone_reserve_aligned(rc.cname, rc.size, rc.socket,
			rc.flags, *rc.aligned)
	default:
		return C.rte_memzone_reserve_bounded(rc.cname, rc.size, rc.socket,
			rc.flags, *rc.aligned, *rc.bound)
	}
}

// Reserve a portion of physical memory.
//
// This function reserves some memory and returns a pointer to a
// correctly filled memzone descriptor. If the allocation cannot be
// done, return NULL.
//
// Note: Reserving memzones with len set to 0 will only attempt to
// allocate memzones from memory that is already available. It will
// not trigger any new allocations.  : When reserving memzones with
// len set to 0, it is preferable to also set a valid socket_id.
// Setting socket_id to SOCKET_ID_ANY is supported, but will likely
// not yield expected results.  Specifically, the resulting memzone
// may not necessarily be the biggest memzone available, but rather
// biggest memzone available on socket id corresponding to an lcore
// from which reservation was called.
func Reserve(name string, size uintptr, opts ...Option) (*Memzone, error) {
	rc := makeOpts(name, size, opts)
	mz := (*Memzone)(doReserve(rc))
	if mz == nil {
		return nil, err()
	}
	return mz, nil
}

// Lookup searches a memzone from its name.
func Lookup(name string) (*Memzone, error) {
	mz := (*Memzone)(C.rte_memzone_lookup(cGoString(name)))
	if mz == nil {
		return nil, err()
	}
	return mz, nil
}

// Free a memzone. EINVAL (invalid parameter) may be returned.
func (mz *Memzone) Free() error {
	return err(C.rte_memzone_free((*C.struct_rte_memzone)(mz)))
}

var (
	callbacks = common.NewRegistryArray()
)

//export mzCb
func mzCb(mz *C.struct_rte_memzone, arg unsafe.Pointer) {
	cb := *(*common.ObjectID)(arg)
	fn := callbacks.Read(cb).(func(*Memzone))
	fn((*Memzone)(mz))
}

// Walk list of all memzones.
func Walk(fn func(*Memzone)) {
	cb := callbacks.Create(fn)
	C.memzone_walk(C.mzCb, unsafe.Pointer(&cb))
	callbacks.Delete(cb)
}

// Name returns name of the memzone.
func (mz *Memzone) Name() string {
	cmz := (*C.struct_rte_memzone)(mz)
	return C.GoString(&cmz.name[0])
}

// Addr returns start virtual address of the memzone.
func (mz *Memzone) Addr() unsafe.Pointer {
	cmz := (*C.struct_rte_memzone)(mz)
	return C.memzone_addr(cmz)
}

// Len returns length of the memzone.
func (mz *Memzone) Len() uintptr {
	cmz := (*C.struct_rte_memzone)(mz)
	return uintptr(cmz.len)
}

// Bytes returns memzone in a form of slice of bytes.
func (mz *Memzone) Bytes() []byte {
	var b []byte
	cmz := (*C.struct_rte_memzone)(mz)
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh.Data = uintptr(C.memzone_addr(cmz))
	sh.Len = int(cmz.len)
	sh.Cap = sh.Len
	return b
}

// SocketID returns NUMA socket ID of the memzone.
func (mz *Memzone) SocketID() int {
	cmz := (*C.struct_rte_memzone)(mz)
	return int(cmz.socket_id)
}

// HugePageSz returns the page size of underlying memory.
func (mz *Memzone) HugePageSz() uint64 {
	cmz := (*C.struct_rte_memzone)(mz)
	return uint64(cmz.hugepage_sz)
}
