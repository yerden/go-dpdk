#ifndef _RING_H_
#define _RING_H_

#define GO_RING_FUNC(func)                                 \
static int func(struct rte_ring *r,                        \
    uintptr_t *objs, unsigned int n,                       \
    unsigned int *ptr) {                                   \
  void **obj_table = (typeof(obj_table))objs;              \
  return rte_ring_ ## func(r, obj_table, n, ptr);          \
}

// wrap dequeue API
GO_RING_FUNC(mc_dequeue_burst)
GO_RING_FUNC(mc_dequeue_bulk)
GO_RING_FUNC(sc_dequeue_burst)
GO_RING_FUNC(sc_dequeue_bulk)
GO_RING_FUNC(dequeue_burst)
GO_RING_FUNC(dequeue_bulk)

// wrap enqueue API
GO_RING_FUNC(mp_enqueue_burst)
GO_RING_FUNC(mp_enqueue_bulk)
GO_RING_FUNC(sp_enqueue_burst)
GO_RING_FUNC(sp_enqueue_bulk)
GO_RING_FUNC(enqueue_burst)
GO_RING_FUNC(enqueue_bulk)

#endif /* _RING_H_ */

