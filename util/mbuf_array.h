#ifndef _RX_BUFFER_H_
#define _RX_BUFFER_H_

#include <stdlib.h>
#include <rte_mbuf.h>
#include <rte_malloc.h>
#include <rte_ethdev.h>

#define RX_BUFFER_SIZE(sz) (sizeof(struct mbuf_array) + ((sz)-1) * sizeof(struct rte_mbuf *))

#ifndef MBUF_ARRAY_USER_SIZE
#define MBUF_ARRAY_USER_SIZE 64
#endif

struct mbuf_array {
	// number of mbufs allocated, occupied length and position
	uint16_t size, length, n;

	uint16_t padding[1];

	// custom data
	uint8_t opaque[MBUF_ARRAY_USER_SIZE];

	// must preallocate one pointer for Go compiler to be happy
	// and able to take pointer of this array.
	struct rte_mbuf *pkts[1];
};

struct ethdev_data {
	uint16_t pid, qid;
};

static int new_mbuf_array(int socket, uint16_t count, struct mbuf_array **pbuf)
{
	size_t size = RX_BUFFER_SIZE(count);
	struct mbuf_array *buf = (typeof(buf)) rte_zmalloc_socket("mbuf_array", size, 0, socket);
	if (buf == NULL) {
		return -ENOMEM;
	}

	buf->size = count;
	*pbuf = buf;
	return 0;
}

static uint16_t mbuf_array_ethdev_reload(struct mbuf_array *buf)
{
	uint16_t n;
	struct ethdev_data *data = (typeof(data)) buf->opaque;

	for (n = 0; n < buf->length; n++) {
		rte_pktmbuf_free(buf->pkts[n]);
	}

	buf->length = rte_eth_rx_burst(data->pid, data->qid, buf->pkts, buf->size);
	buf->n = 0;
	return buf->length;
}

#endif				/* _RX_BUFFER_H_ */
