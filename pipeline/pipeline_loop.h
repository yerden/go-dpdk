#ifndef _PIPELINE_LOOP_H_
#define _PIPELINE_LOOP_H_

#include <stdint.h>
#include <rte_pipeline.h>

struct lcore_arg {
	struct rte_pipeline *p;

	// iterations before calling rte_pipeline_flush()
	// should be power of 2.
	uint32_t flush;

	// if not zero then stop the loop.
	volatile uint8_t *stop;
};

static int run_pipeline_loop(void *arg)
{
	struct lcore_arg *ctx = arg;
	uint32_t n;

	for (n = 0; !*ctx->stop; n++) {
		rte_pipeline_run(ctx->p);

		if (ctx->flush && (n & (ctx->flush - 1)) == 0)
			rte_pipeline_flush(ctx->p);
	}
}

#endif				/* _PIPELINE_LOOP_H_ */
