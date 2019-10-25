#ifndef _SAMPLE_FLOW_
#define _SAMPLE_FLOW_

#include <rte_errno.h>
#include <rte_flow.h>

static int eth_vlan_ip4_udp(uint16_t port_id, struct rte_flow_action_rss *rss,
			    struct rte_flow **pf, struct rte_flow_error *err)
{
	/* Declaring structs being used. 8< */
	struct rte_flow_attr attr;
	struct rte_flow_item pattern[20];
	struct rte_flow_action action[20];
	struct rte_flow *flow = NULL;
	/* >8 End of declaring structs being used. */
	int res;

	memset(pattern, 0, sizeof(pattern));
	memset(action, 0, sizeof(action));

	/* Set the rule attribute, only ingress packets will be checked. 8< */
	memset(&attr, 0, sizeof(struct rte_flow_attr));
	attr.ingress = 1;
	/* >8 End of setting the rule attribute. */

	/*
	 * create the action sequence.
	 * one action only,  move packet to queue
	 */
	action[0].type = RTE_FLOW_ACTION_TYPE_RSS;
	action[0].conf = rss;
	action[1].type = RTE_FLOW_ACTION_TYPE_END;

	/*
	 * set the first level of the pattern (ETH).
	 * since in this example we just want to get the
	 * ipv4 we set this level to allow all.
	 */

	/* IPv4 we set this level to allow all. 8< */
	pattern[0].type = RTE_FLOW_ITEM_TYPE_ETH;
	pattern[1].type = RTE_FLOW_ITEM_TYPE_VLAN;
	pattern[2].type = RTE_FLOW_ITEM_TYPE_IPV4;
	pattern[3].type = RTE_FLOW_ITEM_TYPE_UDP;
	pattern[4].type = RTE_FLOW_ITEM_TYPE_END;
	/* >8 End of setting the first level of the pattern. */

	/* Validate the rule and create it. 8< */
	res = rte_flow_validate(port_id, &attr, pattern, action, err);
	if (!res && pf != NULL) {
		if ((*pf = rte_flow_create(port_id, &attr, pattern, action, err)) == NULL) {
			res = rte_errno;
		}
	}
	/* >8 End of validation the rule and create it. */

	return res;
}

#endif				/* _SAMPLE_FLOW_ */
