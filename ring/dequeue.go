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

// Dequeue dequeues single object from Ring.
func (r *Ring) Dequeue() (uintptr, bool) {
	objs := []uintptr{0}
	n, _ := r.DequeueBulk(objs)
	return objs[0], n != 0
}

// ScDequeue dequeues single object from Ring.
func (r *Ring) ScDequeue() (uintptr, bool) {
	objs := []uintptr{0}
	n, _ := r.ScDequeueBulk(objs)
	return objs[0], n != 0
}

// McDequeue dequeues single object from Ring.
func (r *Ring) McDequeue() (uintptr, bool) {
	objs := []uintptr{0}
	n, _ := r.McDequeueBulk(objs)
	return objs[0], n != 0
}

// McDequeueBulk dequeues objects into given slice of pointers.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of remaining ring entries in the ring after the enqueue
// operation has finished.
func (r *Ring) McDequeueBulk(obj []uintptr) (n, avail uint32) {
	return uint32(C.mc_dequeue_bulk(args(r, obj, &avail))), avail
}

// ScDequeueBulk dequeues objects into given slice of pointers.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of remaining ring entries in the ring after the enqueue
// operation has finished.
func (r *Ring) ScDequeueBulk(obj []uintptr) (n, avail uint32) {
	return uint32(C.sc_dequeue_bulk(args(r, obj, &avail))), avail
}

// DequeueBulk dequeues objects into given slice of pointers.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of remaining ring entries in the ring after the enqueue
// operation has finished.
func (r *Ring) DequeueBulk(obj []uintptr) (n, avail uint32) {
	return uint32(C.dequeue_bulk(args(r, obj, &avail))), avail
}

// McDequeueBurst dequeues objects into given slice of pointers.
// Returns number of dequeued objects and amount of remaining ring
// entries in the ring after the enqueue operation has finished.
// after the enqueue operation has finished.
func (r *Ring) McDequeueBurst(obj []uintptr) (n, avail uint32) {
	return uint32(C.mc_dequeue_burst(args(r, obj, &avail))), avail
}

// ScDequeueBurst dequeues objects into given slice of pointers.
// Returns number of dequeued objects and amount of remaining ring
// entries in the ring after the enqueue operation has finished.
// after the enqueue operation has finished.
func (r *Ring) ScDequeueBurst(obj []uintptr) (n, avail uint32) {
	return uint32(C.sc_dequeue_burst(args(r, obj, &avail))), avail
}

// DequeueBurst dequeues objects into given slice of pointers.
// Returns number of dequeued objects and amount of remaining ring
// entries in the ring after the enqueue operation has finished.
// after the enqueue operation has finished.
func (r *Ring) DequeueBurst(obj []uintptr) (n, avail uint32) {
	return uint32(C.dequeue_burst(args(r, obj, &avail))), avail
}
