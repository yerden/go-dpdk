/*
Package ring wraps RTE ring library.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package ring

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_ring.h>
#include <rte_memory.h>
#include <rte_malloc.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Ring is a fixed-size queue, implemented as a table of pointers.
// Head and tail pointers are modified atomically, allowing concurrent
// access to it. It has the following features:
//
// * FIFO (First In First Out)
//
// * Maximum size is fixed; the pointers are stored in a table.
//
// * Lockless implementation.
//
// * Multi- or single-consumer dequeue.
//
// * Multi- or single-producer enqueue.
//
// * Bulk dequeue.
//
// * Bulk enqueue.
// Note: the ring implementation is not preemptible. Refer to
// Programmer's guide/Environment Abstraction Layer/Multiple
// pthread/Known Issues/rte_ring for more information.
type Ring C.struct_rte_ring

type ringConf struct {
	cname  *C.char
	socket C.int
	flags  C.uint
}

// Option alters ring behaviour.
type Option struct {
	f func(*ringConf)
}

const (
	// SingleConsumer specifies that default dequeue operation will
	// exhibit 'single-consumer' behaviour.
	SingleConsumer uint = C.RING_F_SC_DEQ
	// SingleProducer specifies that default enqueue operation will
	// exhibit 'single-producer' behaviour.
	SingleProducer = C.RING_F_SP_ENQ
	// ExactSize specifies how to handle ring size during Create/Init.
	// Ring is to hold exactly requested number of entries. Without
	// this flag set, the ring size requested must be a power of 2,
	// and the usable space will be that size - 1. With the flag, the
	// requested size will be rounded up to the next power of two, but
	// the usable space will be exactly that requested. Worst case, if
	// a power-of-2 size is requested, half the ring space will be
	// wasted.
	ExactSize = C.RING_F_EXACT_SZ
)

// Shortcuts for ring creation flags.
var (
	OptSC = OptFlag(SingleConsumer)
	OptSP = OptFlag(SingleProducer)
	OptES = OptFlag(ExactSize)
)

// OptSocket specifies the socket id where the memzone would be
// created in Create.
func OptSocket(socket uint) Option {
	return Option{func(rc *ringConf) {
		rc.socket = C.int(socket)
	}}
}

// OptFlag add one of permitted flags for the ring creation.
func OptFlag(flag uint) Option {
	return Option{func(rc *ringConf) {
		rc.flags |= C.uint(flag)
	}}
}

func cGoString(s string) *C.char {
	a := append([]byte(s), 0)
	return (*C.char)(unsafe.Pointer(&a[0]))
}

func makeOpts(name string, opts []Option) *ringConf {
	rc := &ringConf{socket: C.SOCKET_ID_ANY, cname: cGoString(name)}
	for i := range opts {
		opts[i].f(rc)
	}
	return rc
}

// Create creates new ring named name in memory.
//
// This function uses rte_memzone_reserve() to allocate memory. Then
// it calls rte_ring_init() to initialize an empty ring.
//
// The new ring size is set to count, which must be a power of two.
// Water marking is disabled by default. The real usable ring size is
// count-1 instead of count to differentiate a free ring from an empty
// ring.
//
// The ring is added in RTE_TAILQ_RING list.
func Create(name string, count uint, opts ...Option) (*Ring, error) {
	rc := makeOpts(name, opts)
	r := (*Ring)(C.rte_ring_create(rc.cname, C.uint(count), rc.socket, rc.flags))
	if r == nil {
		return nil, common.Errno(nil)
	}
	return r, nil
}

// Init initializes a ring structure in memory. The size of the memory
// area must be large enough to store the ring structure and the
// object table. It is advised to use GetMemSize() to get the
// appropriate size.
//
// The ring size is set to count, which must be a power of two. Water
// marking is disabled by default. The real usable ring size is
// count-1 instead of count to differentiate a free ring from an empty
// ring.
//
// The ring is not added in RTE_TAILQ_RING global list. Indeed, the
// memory given by the caller may not be shareable among dpdk
// processes.
func (r *Ring) Init(name string, count uint, opts ...Option) error {
	rc := makeOpts(name, opts)
	return common.Errno(C.rte_ring_init((*C.struct_rte_ring)(r), rc.cname,
		C.uint(count), rc.flags))
}

// New allocates and initializes Ring in Go memory. It allocates a
// slice of bytes with enough length to hold a Ring with requested
// parameters. Then slice is casted to Ring and initialized with
// Init.
//
// Please note that OptSocket is irrelevant in this case and is unused
// if specified.
func New(name string, count uint, opts ...Option) (*Ring, error) {
	size, err := GetMemSize(count)
	if err != nil {
		return nil, err
	}

	p := make([]byte, size)
	r := (*Ring)(unsafe.Pointer(&p[0]))
	return r, r.Init(name, count, opts...)
}

// GetMemSize calculates the memory size needed for a ring.
//
// This function returns the number of bytes needed for a ring, given
// the number of elements in it. This value is the sum of the size of
// the structure rte_ring and the size of the memory needed by the
// objects pointers. The value is aligned to a cache line size.
//
// count should be power of 2. If that is not the case, EINVAL error
// will be returned.
func GetMemSize(count uint) (int, error) {
	sz := C.rte_ring_get_memsize(C.uint(count))
	return common.IntOrErr(int(sz))
}

// Free deallocates all memory used by the ring.
func (r *Ring) Free() {
	C.rte_ring_free((*C.struct_rte_ring)(r))
}

// Size returns the size of the ring.
//
// NOTE: this is not the same as the usable space in the ring. To
// query that use Cap().
func (r *Ring) Size() uint {
	return uint(C.rte_ring_get_size((*C.struct_rte_ring)(r)))
}

// Cap returns the number of elements which can be stored in the ring.
func (r *Ring) Cap() uint {
	return uint(C.rte_ring_get_capacity((*C.struct_rte_ring)(r)))
}

// FreeCount returns the number of free entries in a ring.
func (r *Ring) FreeCount() uint {
	return uint(C.rte_ring_free_count((*C.struct_rte_ring)(r)))
}

// IsFull tests if the ring is full.
func (r *Ring) IsFull() bool {
	return C.rte_ring_full((*C.struct_rte_ring)(r)) != 0
}

// IsEmpty tests if the ring is empty.
func (r *Ring) IsEmpty() bool {
	return C.rte_ring_empty((*C.struct_rte_ring)(r)) != 0
}

// Lookup searches a ring from its name in RTE_TAILQ_RING, i.e. among
// those created with Create.
func Lookup(name string) (*Ring, bool) {
	rc := makeOpts(name, nil)
	r := (*Ring)(C.rte_ring_lookup(rc.cname))
	return r, r != nil
}
