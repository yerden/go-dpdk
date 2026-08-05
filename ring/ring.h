#ifndef _RING_H_
#define _RING_H_

struct compound_int {
	unsigned int n;
	unsigned int rc;
};

typedef void * ptr_t;

#define GO_RING_FUNC_ELEM(func)                            \
static struct compound_int func ## _elem(                  \
    struct rte_ring *r, unsigned int esize,                \
    uintptr_t objs, unsigned int n) {                      \
  struct compound_int out;                                 \
  void *obj_table = (typeof(obj_table))objs;               \
  out.rc = rte_ring_ ## func ## _elem(r, obj_table, esize, \
    n, &out.n); \
  return out;                                              \
}

#define GO_RING_OP(func)                                   \
	GO_RING_FUNC_ELEM(func)

// wrap dequeue API
GO_RING_OP(mc_dequeue_burst)
GO_RING_OP(mc_dequeue_bulk)
GO_RING_OP(sc_dequeue_burst)
GO_RING_OP(sc_dequeue_bulk)
GO_RING_OP(dequeue_burst)
GO_RING_OP(dequeue_bulk)

// wrap enqueue API
GO_RING_OP(mp_enqueue_burst)
GO_RING_OP(mp_enqueue_bulk)
GO_RING_OP(sp_enqueue_burst)
GO_RING_OP(sp_enqueue_bulk)
GO_RING_OP(enqueue_burst)
GO_RING_OP(enqueue_bulk)

#endif /* _RING_H_ */

