package ring

/*
#include <rte_config.h>
#include <rte_ring.h>
#include <rte_memory.h>
#include <rte_malloc.h>

enum {
	OFF_RING_CONS = offsetof(struct rte_ring, cons),
	OFF_RING_PROD = offsetof(struct rte_ring, prod),
};
*/
import "C"
import (
	"sync/atomic"
	"unsafe"
)

type ringBehavior uint32

const (
	queueFixed    ringBehavior = C.RTE_RING_QUEUE_FIXED
	queueVariable ringBehavior = C.RTE_RING_QUEUE_VARIABLE
)

type syncType uint32

const (
	syncMT  syncType = C.RTE_RING_SYNC_MT
	syncST  syncType = C.RTE_RING_SYNC_ST
	syncRTS syncType = C.RTE_RING_SYNC_MT_RTS
	syncHTS syncType = C.RTE_RING_SYNC_MT_HTS
)

type ringHeadTail C.struct_rte_ring_headtail
type htsHeadtail C.struct_rte_ring_hts_headtail
type rtsHeadtail C.struct_rte_ring_rts_headtail

type headTail struct {
	head uint32
	tail uint32
}

func (r *Ring) cons() unsafe.Pointer {
	p := unsafe.Pointer(r)
	return unsafe.Add(p, C.OFF_RING_CONS)
}

func (r *Ring) prod() unsafe.Pointer {
	p := unsafe.Pointer(r)
	return unsafe.Add(p, C.OFF_RING_PROD)
}

func (r *Ring) ringHeadTailProd() *ringHeadTail {
	return (*ringHeadTail)(r.prod())
}

func (r *Ring) ringHeadTailCons() *ringHeadTail {
	return (*ringHeadTail)(r.cons())
}

func (r *Ring) htsHeadtailProd() *htsHeadtail {
	return (*htsHeadtail)(r.prod())
}

func (r *Ring) htsHeadtailCons() *htsHeadtail {
	return (*htsHeadtail)(r.cons())
}

func (r *Ring) rtsHeadtailProd() *rtsHeadtail {
	return (*rtsHeadtail)(r.prod())
}

func (r *Ring) rtsHeadtailCons() *rtsHeadtail {
	return (*rtsHeadtail)(r.cons())
}

func (r *Ring) moveConsHead(isSc bool, n uint32, behavior ringBehavior, oldHead, newHead, entries *uint32) uint32 {
	cons := (*headTail)(r.cons())
	prod := (*headTail)(r.prod())
	max := n
	success := false

	for {
		n = max

		// XXX: there is a memory barrier in this place to protect
		// reordering of load/load. In Go we place atomic load for
		// uint32 instead.
		*oldHead = atomic.LoadUint32(&cons.head)

		*entries = prod.tail - *oldHead

		if n > *entries {
			if behavior == queueFixed {
				n = 0
			} else {
				n = *entries
			}
		}

		if n == 0 {
			return 0
		}

		*newHead = *oldHead + n

		if isSc {
			cons.head = *newHead
			success = true
		} else {
			success = atomic.CompareAndSwapUint32(&cons.head, *oldHead, *newHead)
		}

		if success {
			break
		}
	}

	return n
}

func (r *Ring) dequeueElems(oldHead uint32, objTable unsafe.Pointer, eSize, n uint32) {
	panic("TODO")
}

func (r *Ring) updateTail(ht *headTail, oldHead, nextHead uint32, isSc, enqueue bool) {
	panic("TODO")
}

func (r *Ring) doDequeueElem(objTable unsafe.Pointer, eSize, n uint32, behavior ringBehavior, st syncType, available *uint32) uint32 {
	var consHead, consNext, entries uint32
	cons := (*headTail)(r.cons())
	isSc := st != syncMT

	n = r.moveConsHead(isSc, n, behavior, &consHead, &consNext, &entries)
	if n == 0 {
		goto end
	}

	r.dequeueElems(consHead, objTable, eSize, n)
	r.updateTail(cons, consHead, consNext, isSc, false)

end:
	if available != nil {
		*available = entries - n
	}
	return n
}
