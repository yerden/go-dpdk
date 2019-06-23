package ring

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_ring.h>
#include <rte_errno.h>
#include <rte_memory.h>

#include "ring.h"

static int go_rte_errno() {
	return rte_errno;
}

struct someptr {
	void *p;
};
*/
import "C"

import (
	"unsafe"
)

var (
	// uintptr is just the size of void *
	_ [C.sizeof_struct_someptr]struct{} = [unsafe.Sizeof(uintptr(0))]struct{}{}
)

// Enqueue enqueues an object into given Ring.
func (r *Ring) Enqueue(obj uintptr) bool {
	n, _ := r.EnqueueBulk([]uintptr{obj})
	return n != 0
}

// SpEnqueue enqueues an object into given Ring.
func (r *Ring) SpEnqueue(obj uintptr) bool {
	n, _ := r.SpEnqueueBulk([]uintptr{obj})
	return n != 0
}

// MpEnqueue enqueues an object into given Ring.
func (r *Ring) MpEnqueue(obj uintptr) bool {
	n, _ := r.MpEnqueueBulk([]uintptr{obj})
	return n != 0
}

func args(r *Ring, obj []uintptr, ptr *uint32) (*C.struct_rte_ring,
	*C.uintptr_t, C.uint, *C.uint) {
	return (*C.struct_rte_ring)(r), (*C.uintptr_t)(unsafe.Pointer(&obj[0])),
		C.uint(len(obj)), (*C.uint)(ptr)
}

// MpEnqueueBulk enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) MpEnqueueBulk(obj []uintptr) (n, free uint32) {
	return uint32(C.mp_enqueue_bulk(args(r, obj, &free))), free
}

// SpEnqueueBulk enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) SpEnqueueBulk(obj []uintptr) (n, free uint32) {
	return uint32(C.sp_enqueue_bulk(args(r, obj, &free))), free
}

// EnqueueBulk enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) EnqueueBulk(obj []uintptr) (n, free uint32) {
	return uint32(C.enqueue_bulk(args(r, obj, &free))), free
}

// MpEnqueueBurst enqueues given objects from slice into Ring.
// Returns number of enqueued objects and amount of space in the ring
// after the enqueue operation has finished.
func (r *Ring) MpEnqueueBurst(obj []uintptr) (n, free uint32) {
	return uint32(C.mp_enqueue_burst(args(r, obj, &free))), free
}

// SpEnqueueBurst enqueues given objects from slice into Ring.
// Returns number of enqueued objects and amount of space in the ring
// after the enqueue operation has finished.
func (r *Ring) SpEnqueueBurst(obj []uintptr) (n, free uint32) {
	return uint32(C.sp_enqueue_burst(args(r, obj, &free))), free
}

// EnqueueBurst enqueues given objects from slice into Ring. Returns
// number of enqueued objects and amount of space in the ring after
// the enqueue operation has finished.
func (r *Ring) EnqueueBurst(obj []uintptr) (n, free uint32) {
	return uint32(C.enqueue_burst(args(r, obj, &free))), free
}
