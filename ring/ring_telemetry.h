#ifndef _RING_TELEMETRY_H_
#define _RING_TELEMETRY_H_

#include <rte_common.h>
#include <rte_telemetry.h>
#include <rte_ring.h>

#include <rte_version.h>
#if RTE_VERSION < RTE_VERSION_NUM(23, 3, 0, 0)
#define tel_data_add_dict_uint rte_tel_data_add_dict_u64
#else
#define tel_data_add_dict_uint rte_tel_data_add_dict_uint
#endif

static inline const char *
trim_prefix(const char *s, const char *pre)
{
	int len = strlen(pre);
	return s + len * !strncmp(s, pre, len);
}

static void
memzone_reap_rings(
		const struct rte_memzone *mz,
		void *arg)
{
	struct rte_tel_data *d = (struct rte_tel_data *)(arg);
	const char *name = trim_prefix(mz->name, RTE_RING_MZ_PREFIX);
	if (rte_ring_lookup(name))
		rte_tel_data_add_array_string(d, name);
}

static int
ring_list(
		__rte_unused const char *cmd,
		__rte_unused const char *params,
		struct rte_tel_data *d)
{
	rte_tel_data_start_array(d, RTE_TEL_STRING_VAL);
	rte_memzone_walk(memzone_reap_rings, d);
	return 0;
}

static inline int
ring_info(
		__rte_unused const char *cmd,
		const char *name,
		struct rte_tel_data *d)
{
	if (!name|| strlen(name) == 0)
		return -EINVAL;

	struct rte_ring *r = rte_ring_lookup(name);
	if (r == NULL)
		return -ENOENT;

	rte_tel_data_start_dict(d);
	rte_tel_data_add_dict_string(d, "ring_name", name);
	rte_tel_data_add_dict_string(d, "ring_mz_name", r->memzone->name);
	rte_tel_data_add_dict_int(d, "ring_is_sc", rte_ring_is_cons_single(r));
	rte_tel_data_add_dict_int(d, "ring_is_sp", rte_ring_is_prod_single(r));
	rte_tel_data_add_dict_int(d, "ring_prod_sync_type", rte_ring_get_prod_sync_type(r));
	rte_tel_data_add_dict_int(d, "ring_cons_sync_type", rte_ring_get_cons_sync_type(r));
	tel_data_add_dict_uint(d, "ring_size", rte_ring_get_size(r));
	tel_data_add_dict_uint(d, "ring_count", rte_ring_count(r));
	tel_data_add_dict_uint(d, "ring_capacity", rte_ring_get_capacity(r));

	return 0;
}

telemetry_cb ring_list_cb = ring_list;
telemetry_cb ring_info_cb = ring_info;

#endif /* _RING_TELEMETRY_H_ */
