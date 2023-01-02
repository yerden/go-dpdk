#ifndef _PIPELINE_LOOP_H_
#define _PIPELINE_LOOP_H_

#include <stdint.h>
#include <rte_pipeline.h>


/**
 * Test if the pipeline should be stopped.
 *
 * @param params
 *   Handle to private data.
 *
 * @param p
 *   Pipeline instance.
 *
 * @return
 *   Return 0 if the pipeline resumes. Return >0 if the pipeline
 *   stops.
 */
typedef int (*pipeline_op_ctrl)(
		void *params,
		struct rte_pipeline *p);

// Pipeline loop operations.
struct pipeline_ops {
	pipeline_op_ctrl f_ctrl;
};

struct lcore_arg {
	struct rte_pipeline *p;
	struct pipeline_ops ops;
	void *ops_arg;
	uint32_t flush;
};

static int run_pipeline_loop(void *arg)
{
	struct lcore_arg *ctx = arg;
	uint32_t n;
	int rc;
	struct pipeline_ops *ops = &ctx->ops;
	struct rte_pipeline *p = ctx->p;

	for (n = 0; ; n++) {
		rte_pipeline_run(p);

		if (!ctx->flush || ((n & (ctx->flush - 1)) != 0))
			continue;

		rte_pipeline_flush(p);

		if ((rc = ops->f_ctrl(ctx->ops_arg, p)) > 0) {
			/* pipeline is signalled to stop */
			return 0;
		}

		if (rc < 0) {
			/* pipeline control error */
			break;
		}
	}

	return rc;
}

#endif				/* _PIPELINE_LOOP_H_ */
