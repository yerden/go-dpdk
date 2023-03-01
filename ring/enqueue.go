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

// Enqueue enqueues an object into given Ring.
func (r *Ring) Enqueue(obj unsafe.Pointer) bool {
	n, _ := r.EnqueueBulk([]unsafe.Pointer{obj})
	return n != 0
}

// SpEnqueue enqueues an object into given Ring.
func (r *Ring) SpEnqueue(obj unsafe.Pointer) bool {
	n, _ := r.SpEnqueueBulk([]unsafe.Pointer{obj})
	return n != 0
}

// MpEnqueue enqueues an object into given Ring.
func (r *Ring) MpEnqueue(obj unsafe.Pointer) bool {
	n, _ := r.MpEnqueueBulk([]unsafe.Pointer{obj})
	return n != 0
}

// MpEnqueueBulk enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) MpEnqueueBulk(obj []unsafe.Pointer) (n, free uint32) {
	return ret(C.mp_enqueue_bulk_elem(argsSliceElem(r, obj)))
}

// SpEnqueueBulk enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) SpEnqueueBulk(obj []unsafe.Pointer) (n, free uint32) {
	return ret(C.sp_enqueue_bulk_elem(argsSliceElem(r, obj)))
}

// EnqueueBulk enqueues given objects from slice into Ring.
// Returns number of enqueued objects (either 0 or len(obj)) and
// amount of space in the ring after the enqueue operation has
// finished.
func (r *Ring) EnqueueBulk(obj []unsafe.Pointer) (n, free uint32) {
	return ret(C.enqueue_bulk_elem(argsSliceElem(r, obj)))
}

// MpEnqueueBurst enqueues given objects from slice into Ring.
// Returns number of enqueued objects and amount of space in the ring
// after the enqueue operation has finished.
func (r *Ring) MpEnqueueBurst(obj []unsafe.Pointer) (n, free uint32) {
	return ret(C.mp_enqueue_burst_elem(argsSliceElem(r, obj)))
}

// SpEnqueueBurst enqueues given objects from slice into Ring.
// Returns number of enqueued objects and amount of space in the ring
// after the enqueue operation has finished.
func (r *Ring) SpEnqueueBurst(obj []unsafe.Pointer) (n, free uint32) {
	return ret(C.sp_enqueue_burst_elem(argsSliceElem(r, obj)))
}

// EnqueueBurst enqueues given objects from slice into Ring. Returns
// number of enqueued objects and amount of space in the ring after
// the enqueue operation has finished.
func (r *Ring) EnqueueBurst(obj []unsafe.Pointer) (n, free uint32) {
	return ret(C.enqueue_burst_elem(argsSliceElem(r, obj)))
}
