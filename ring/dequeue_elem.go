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

// DequeueElem dequeues an object from given Ring.
func (r *Ring) DequeueElem(obj unsafe.Pointer, esize uintptr) bool {
	n, _ := r.DequeueBulkElem(obj, esize, 1)
	return n != 0
}

// ScDequeueElem dequeues an object from given Ring.
func (r *Ring) ScDequeueElem(obj unsafe.Pointer, esize uintptr) bool {
	n, _ := r.ScDequeueBulkElem(obj, esize, 1)
	return n != 0
}

// McDequeueElem dequeues an object from given Ring.
func (r *Ring) McDequeueElem(obj unsafe.Pointer, esize uintptr) bool {
	n, _ := r.McDequeueBulkElem(obj, esize, 1)
	return n != 0
}

// McDequeueBulkElem dequeues given objects to slice from Ring.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of space in the ring after the dequeue operation has
// finished.
func (r *Ring) McDequeueBulkElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.mc_dequeue_bulk_elem(argsElem(r, objtable, esize, n)))
}

// ScDequeueBulkElem dequeues given objects to slice from Ring.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of space in the ring after the dequeue operation has
// finished.
func (r *Ring) ScDequeueBulkElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.sc_dequeue_bulk_elem(argsElem(r, objtable, esize, n)))
}

// DequeueBulkElem dequeues given objects to slice from Ring.
// Returns number of dequeued objects (either 0 or len(obj)) and
// amount of space in the ring after the dequeue operation has
// finished.
func (r *Ring) DequeueBulkElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.dequeue_bulk_elem(argsElem(r, objtable, esize, n)))
}

// McDequeueBurstElem dequeues given objects to slice from Ring.
// Returns number of dequeued objects and amount of space in the ring
// after the dequeue operation has finished.
func (r *Ring) McDequeueBurstElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.mc_dequeue_burst_elem(argsElem(r, objtable, esize, n)))
}

// ScDequeueBurstElem dequeues given objects to slice from Ring.
// Returns number of dequeued objects and amount of space in the ring
// after the dequeue operation has finished.
func (r *Ring) ScDequeueBurstElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.sc_dequeue_burst_elem(argsElem(r, objtable, esize, n)))
}

// DequeueBurstElem dequeues given objects to slice from Ring. Returns
// number of dequeued objects and amount of space in the ring after
// the dequeue operation has finished.
func (r *Ring) DequeueBurstElem(objtable unsafe.Pointer, esize uintptr, n int) (processed, free uint32) {
	return ret(C.dequeue_burst_elem(argsElem(r, objtable, esize, n)))
}
