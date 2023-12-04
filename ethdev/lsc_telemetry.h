#ifndef _LSC_TELEMETRY_H_
#define _LSC_TELEMETRY_H_

#include <rte_common.h>
#include <rte_telemetry.h>

#include <rte_ethdev.h>
#include <rte_atomic.h>

#include <rte_version.h>
#if RTE_VERSION < RTE_VERSION_NUM(23, 3, 0, 0)
#define tel_data_add_dict_uint rte_tel_data_add_dict_u64
#else
#define tel_data_add_dict_uint rte_tel_data_add_dict_uint
#endif

struct lsc_counters {
	struct {
		int enabled;
		rte_atomic64_t counter;
	} ports[RTE_MAX_ETHPORTS];
};

extern struct lsc_counters global_lsc_counters;

static int
lsc_counters_callback(
		__rte_unused uint16_t port_id,
		enum rte_eth_event_type event,
		void *cb_arg,
		__rte_unused void *ret_param)
{
	rte_atomic64_t *counter = cb_arg;
	rte_atomic64_add(counter, event == RTE_ETH_EVENT_INTR_LSC);
	return 0;
}

static int
lsc_counters_callback_register(uint16_t port_id)
{
	struct lsc_counters *c = &global_lsc_counters;
	if (!rte_eth_dev_is_valid_port(port_id))
		return -EINVAL;

	c->ports[port_id].enabled = 1;
	rte_atomic64_t *counter = &c->ports[port_id].counter;
	return rte_eth_dev_callback_register(port_id, RTE_ETH_EVENT_INTR_LSC,
			lsc_counters_callback, counter);
}

static int
lsc_counters_callback_unregister(uint16_t port_id)
{
	struct lsc_counters *c = &global_lsc_counters;
	if (!rte_eth_dev_is_valid_port(port_id))
		return -EINVAL;

	c->ports[port_id].enabled = 0;
	rte_atomic64_t *counter = &c->ports[port_id].counter;
	return rte_eth_dev_callback_unregister(port_id, RTE_ETH_EVENT_INTR_LSC,
			lsc_counters_callback, counter);
}

static int
ethdev_lsc(
		__rte_unused const char *cmd,
		const char *params,
		struct rte_tel_data *d)
{
	struct lsc_counters *c = &global_lsc_counters;
	int port_id;
	char *end_param;

	if (params == NULL || strlen(params) == 0)
		return -1;

	port_id = strtoul(params, &end_param, 0);
	if (*end_param != '\0')
		return -1;

	if (!rte_eth_dev_is_valid_port(port_id))
		return -ENOENT;

	rte_tel_data_start_dict(d);
	rte_tel_data_add_dict_int(d, "enabled", c->ports[port_id].enabled);
	tel_data_add_dict_uint(d, "lsc_counter", rte_atomic64_read(&c->ports[port_id].counter));
	return 0;
}

static inline void
lsc_register_telemetry_cmd(const char *name)
{
	rte_telemetry_register_cmd(name, ethdev_lsc,
		"Returns LSC counter for specified port id");
}

#endif /* _LSC_TELEMETRY_H_ */
