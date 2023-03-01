package ring

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_ring.h>
#include <rte_errno.h>
#include <rte_memory.h>

#include "ring.h"
*/
import "C"

import (
	"unsafe"
)

const (
	ptrSize = unsafe.Sizeof(unsafe.Pointer(nil))
)

// Dequeue dequeues single object from Ring.
func (r *Ring) Dequeue() (unsafe.Pointer, bool) {
	var p unsafe.Pointer
	return p, r.DequeueElem(unsafe.Pointer(&p), ptrSize)
}

// ScDequeue dequeues single object from Ring.
func (r *Ring) ScDequeue() (unsafe.Pointer, bool) {
	var p unsafe.Pointer
	return p, r.ScDequeueElem(unsafe.Pointer(&p), ptrSize)
}

// McDequeue dequeues single object from Ring.
func (r *Ring) McDequeue() (unsafe.Pointer, bool) {
	var p unsafe.Pointer
	return p, r.McDequeueElem(unsafe.Pointer(&p), ptrSize)
}

// McDequeueBulk dequeues objects into given slice of pointers.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of remaining ring entries in the ring after the enqueue
// operation has finished.
func (r *Ring) McDequeueBulk(obj []unsafe.Pointer) (n, avail uint32) {
	return ret(C.mc_dequeue_bulk_elem(argsSliceElem(r, obj)))
}

// ScDequeueBulk dequeues objects into given slice of pointers.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of remaining ring entries in the ring after the enqueue
// operation has finished.
func (r *Ring) ScDequeueBulk(obj []unsafe.Pointer) (n, avail uint32) {
	return ret(C.sc_dequeue_bulk_elem(argsSliceElem(r, obj)))
}

// DequeueBulk dequeues objects into given slice of pointers.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of remaining ring entries in the ring after the enqueue
// operation has finished.
func (r *Ring) DequeueBulk(obj []unsafe.Pointer) (n, avail uint32) {
	return ret(C.dequeue_bulk_elem(argsSliceElem(r, obj)))
}

// McDequeueBurst dequeues objects into given slice of pointers.
// Returns number of dequeued objects and amount of remaining ring
// entries in the ring after the enqueue operation has finished.
// after the enqueue operation has finished.
func (r *Ring) McDequeueBurst(obj []unsafe.Pointer) (n, avail uint32) {
	return ret(C.mc_dequeue_burst_elem(argsSliceElem(r, obj)))
}

// ScDequeueBurst dequeues objects into given slice of pointers.
// Returns number of dequeued objects and amount of remaining ring
// entries in the ring after the enqueue operation has finished.
// after the enqueue operation has finished.
func (r *Ring) ScDequeueBurst(obj []unsafe.Pointer) (n, avail uint32) {
	return ret(C.sc_dequeue_burst_elem(argsSliceElem(r, obj)))
}

// DequeueBurst dequeues objects into given slice of pointers.
// Returns number of dequeued objects and amount of remaining ring
// entries in the ring after the enqueue operation has finished.
// after the enqueue operation has finished.
func (r *Ring) DequeueBurst(obj []unsafe.Pointer) (n, avail uint32) {
	return ret(C.dequeue_burst_elem(argsSliceElem(r, obj)))
}
