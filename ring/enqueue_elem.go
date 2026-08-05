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

// EnqueueElem enqueues an object into given Ring.
func (r *Ring) EnqueueElem(obj unsafe.Pointer, esize uintptr) bool {
	n, _ := r.EnqueueBulkElem(obj, esize, 1)
	return n != 0
}

// SpEnqueueElem enqueues an object into given Ring.
func (r *Ring) SpEnqueueElem(obj unsafe.Pointer, esize uintptr) bool {
	n, _ := r.SpEnqueueBulkElem(obj, esize, 1)
	return n != 0
}

// MpEnqueueElem enqueues an object into given Ring.
func (r *Ring) MpEnqueueElem(obj unsafe.Pointer, esize uintptr) bool {
	n, _ := r.MpEnqueueBulkElem(obj, esize, 1)
	return n != 0
}

// MpEnqueueBulkElem enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) MpEnqueueBulkElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.mp_enqueue_bulk_elem(argsElem(r, objtable, esize, n)))
}

// SpEnqueueBulkElem enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) SpEnqueueBulkElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.sp_enqueue_bulk_elem(argsElem(r, objtable, esize, n)))
}

// EnqueueBulkElem enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) EnqueueBulkElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.enqueue_bulk_elem(argsElem(r, objtable, esize, n)))
}

// MpEnqueueBurstElem enqueues given objects from slice into Ring.
// Returns number of enqueued objects and amount of space in the ring
// after the enqueue operation has finished.
func (r *Ring) MpEnqueueBurstElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.mp_enqueue_burst_elem(argsElem(r, objtable, esize, n)))
}

// SpEnqueueBurstElem enqueues given objects from slice into Ring.
// Returns number of enqueued objects and amount of space in the ring
// after the enqueue operation has finished.
func (r *Ring) SpEnqueueBurstElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.sp_enqueue_burst_elem(argsElem(r, objtable, esize, n)))
}

// EnqueueBurstElem enqueues given objects from slice into Ring. Returns
// number of enqueued objects and amount of space in the ring after
// the enqueue operation has finished.
func (r *Ring) EnqueueBurstElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.enqueue_burst_elem(argsElem(r, objtable, esize, n)))
}
