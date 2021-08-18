#ifndef _RX_BUFFER_H_
#define _RX_BUFFER_H_

#include <stdlib.h>
#include <rte_mbuf.h>
#include <rte_malloc.h>
#include <rte_ethdev.h>

#define RX_BUFFER_SIZE(sz) (sizeof(struct rx_buffer) + ((sz)-1) * sizeof(struct rte_mbuf *))

struct rx_buffer {
	// port and queue id, must be set manually.
	uint16_t pid, qid;

	// number of mbufs allocated and occupied length.
	uint16_t size, length;

	// cursor position.
	uint16_t n;

	uint16_t padding[3];

	// must preallocate one pointer for Go compiler to be happy
	// and able to take pointer of this array.
	struct rte_mbuf *pkts[1];
};

static int new_rx_buffer(int socket, uint16_t count, struct rx_buffer **pbuf)
{
	size_t size = RX_BUFFER_SIZE(count);
	struct rx_buffer *buf = (typeof(buf)) rte_zmalloc_socket("rx_buffer", size, 0, socket);
	if (buf == NULL) {
		return -ENOMEM;
	}

	buf->size = count;
	*pbuf = buf;
	return 0;
}

static uint16_t recharge_rx_buffer(struct rx_buffer *buf)
{
	uint16_t n;

	for (n = 0; n < buf->length; n++) {
		rte_pktmbuf_free(buf->pkts[n]);
	}

	buf->length = rte_eth_rx_burst(buf->pid, buf->qid, buf->pkts, buf->size);
	buf->n = 0;
	return buf->length;
}

#endif				/* _RX_BUFFER_H_ */
