/*
Package flow encaps RTE generic flow API.

This interface provides the ability to program packet matching and
associated actions in hardware through flow rules.
*/
package flow

/*
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"

// Flow item type constants.
const (
	/**
	 * [META]
	 *
	 * End marker for item lists. Prevents further processing of items,
	 * thereby ending the pattern.
	 *
	 * No associated specification structure.
	 */
	ItemTypeEnd ItemType = C.RTE_FLOW_ITEM_TYPE_END

	/**
	 * [META]
	 *
	 * Used as a placeholder for convenience. It is ignored and simply
	 * discarded by PMDs.
	 *
	 * No associated specification structure.
	 */
	ItemTypeVoid ItemType = C.RTE_FLOW_ITEM_TYPE_VOID

	/**
	 * [META]
	 *
	 * Inverted matching, i.e. process packets that do not match the
	 * pattern.
	 *
	 * No associated specification structure.
	 */
	ItemTypeInvert ItemType = C.RTE_FLOW_ITEM_TYPE_INVERT

	/**
	 * Matches any protocol in place of the current layer, a single ANY
	 * may also stand for several protocol layers.
	 *
	 * See struct rte_flow_item_any.
	 */
	ItemTypeAny ItemType = C.RTE_FLOW_ITEM_TYPE_ANY

	/**
	 * [META]
	 *
	 * Matches traffic originating from (ingress) or going to (egress)
	 * the physical function of the current device.
	 *
	 * No associated specification structure.
	 */
	ItemTypePf ItemType = C.RTE_FLOW_ITEM_TYPE_PF

	/**
	 * [META]
	 *
	 * Matches traffic originating from (ingress) or going to (egress) a
	 * given virtual function of the current device.
	 *
	 * See struct rte_flow_item_vf.
	 */
	ItemTypeVf ItemType = C.RTE_FLOW_ITEM_TYPE_VF

	/**
	 * [META]
	 *
	 * Matches traffic originating from (ingress) or going to (egress) a
	 * physical port of the underlying device.
	 *
	 * See struct rte_flow_item_phy_port.
	 */
	ItemTypePhyPort ItemType = C.RTE_FLOW_ITEM_TYPE_PHY_PORT

	/**
	 * [META]
	 *
	 * Matches traffic originating from (ingress) or going to (egress) a
	 * given DPDK port ID.
	 *
	 * See struct rte_flow_item_port_id.
	 */
	ItemTypePortID ItemType = C.RTE_FLOW_ITEM_TYPE_PORT_ID

	/**
	 * Matches a byte string of a given length at a given offset.
	 *
	 * See struct rte_flow_item_raw.
	 */
	ItemTypeRaw ItemType = C.RTE_FLOW_ITEM_TYPE_RAW

	/**
	 * Matches an Ethernet header.
	 *
	 * See struct rte_flow_item_eth.
	 */
	ItemTypeEth ItemType = C.RTE_FLOW_ITEM_TYPE_ETH

	/**
	 * Matches an 802.1Q/ad VLAN tag.
	 *
	 * See struct rte_flow_item_vlan.
	 */
	ItemTypeVlan ItemType = C.RTE_FLOW_ITEM_TYPE_VLAN

	/**
	 * Matches an IPv4 header.
	 *
	 * See struct rte_flow_item_ipv4.
	 */
	ItemTypeIpv4 ItemType = C.RTE_FLOW_ITEM_TYPE_IPV4

	/**
	 * Matches an IPv6 header.
	 *
	 * See struct rte_flow_item_ipv6.
	 */
	ItemTypeIpv6 ItemType = C.RTE_FLOW_ITEM_TYPE_IPV6

	/**
	 * Matches an ICMP header.
	 *
	 * See struct rte_flow_item_icmp.
	 */
	ItemTypeICMP ItemType = C.RTE_FLOW_ITEM_TYPE_ICMP

	/**
	 * Matches a UDP header.
	 *
	 * See struct rte_flow_item_udp.
	 */
	ItemTypeUDP ItemType = C.RTE_FLOW_ITEM_TYPE_UDP

	/**
	 * Matches a TCP header.
	 *
	 * See struct rte_flow_item_tcp.
	 */
	ItemTypeTCP ItemType = C.RTE_FLOW_ITEM_TYPE_TCP

	/**
	 * Matches a SCTP header.
	 *
	 * See struct rte_flow_item_sctp.
	 */
	ItemTypeSCTP ItemType = C.RTE_FLOW_ITEM_TYPE_SCTP

	/**
	 * Matches a VXLAN header.
	 *
	 * See struct rte_flow_item_vxlan.
	 */
	ItemTypeVxlan ItemType = C.RTE_FLOW_ITEM_TYPE_VXLAN

	/**
	 * Matches a E_TAG header.
	 *
	 * See struct rte_flow_item_e_tag.
	 */
	ItemTypeETag ItemType = C.RTE_FLOW_ITEM_TYPE_E_TAG

	/**
	 * Matches a NVGRE header.
	 *
	 * See struct rte_flow_item_nvgre.
	 */
	ItemTypeNvgre ItemType = C.RTE_FLOW_ITEM_TYPE_NVGRE

	/**
	 * Matches a MPLS header.
	 *
	 * See struct rte_flow_item_mpls.
	 */
	ItemTypeMpls ItemType = C.RTE_FLOW_ITEM_TYPE_MPLS

	/**
	 * Matches a GRE header.
	 *
	 * See struct rte_flow_item_gre.
	 */
	ItemTypeGre ItemType = C.RTE_FLOW_ITEM_TYPE_GRE

	/**
	 * [META]
	 *
	 * Fuzzy pattern match, expect faster than default.
	 *
	 * This is for device that support fuzzy matching option.
	 * Usually a fuzzy matching is fast but the cost is accuracy.
	 *
	 * See struct rte_flow_item_fuzzy.
	 */
	ItemTypeFuzzy ItemType = C.RTE_FLOW_ITEM_TYPE_FUZZY

	/**
	 * Matches a GTP header.
	 *
	 * Configure flow for GTP packets.
	 *
	 * See struct rte_flow_item_gtp.
	 */
	ItemTypeGtp ItemType = C.RTE_FLOW_ITEM_TYPE_GTP

	/**
	 * Matches a GTP header.
	 *
	 * Configure flow for GTP-C packets.
	 *
	 * See struct rte_flow_item_gtp.
	 */
	ItemTypeGtpc ItemType = C.RTE_FLOW_ITEM_TYPE_GTPC

	/**
	 * Matches a GTP header.
	 *
	 * Configure flow for GTP-U packets.
	 *
	 * See struct rte_flow_item_gtp.
	 */
	ItemTypeGtpu ItemType = C.RTE_FLOW_ITEM_TYPE_GTPU

	/**
	 * Matches a ESP header.
	 *
	 * See struct rte_flow_item_esp.
	 */
	ItemTypeEsp ItemType = C.RTE_FLOW_ITEM_TYPE_ESP

	/**
	 * Matches a GENEVE header.
	 *
	 * See struct rte_flow_item_geneve.
	 */
	ItemTypeGeneve ItemType = C.RTE_FLOW_ITEM_TYPE_GENEVE

	/**
	 * Matches a VXLAN-GPE header.
	 *
	 * See struct rte_flow_item_vxlan_gpe.
	 */
	ItemTypeVxlanGpe ItemType = C.RTE_FLOW_ITEM_TYPE_VXLAN_GPE

	/**
	 * Matches an ARP header for Ethernet/IPv4.
	 *
	 * See struct rte_flow_item_arp_eth_ipv4.
	 */
	ItemTypeArpEthIpv4 ItemType = C.RTE_FLOW_ITEM_TYPE_ARP_ETH_IPV4

	/**
	 * Matches the presence of any IPv6 extension header.
	 *
	 * See struct rte_flow_item_ipv6_ext.
	 */
	ItemTypeIpv6Ext ItemType = C.RTE_FLOW_ITEM_TYPE_IPV6_EXT

	/**
	 * Matches any ICMPv6 header.
	 *
	 * See struct rte_flow_item_icmp6.
	 */
	ItemTypeICMP6 ItemType = C.RTE_FLOW_ITEM_TYPE_ICMP6

	/**
	 * Matches an ICMPv6 neighbor discovery solicitation.
	 *
	 * See struct rte_flow_item_icmp6_nd_ns.
	 */
	ItemTypeICMP6NdNs ItemType = C.RTE_FLOW_ITEM_TYPE_ICMP6_ND_NS

	/**
	 * Matches an ICMPv6 neighbor discovery advertisement.
	 *
	 * See struct rte_flow_item_icmp6_nd_na.
	 */
	ItemTypeICMP6NdNa ItemType = C.RTE_FLOW_ITEM_TYPE_ICMP6_ND_NA

	/**
	 * Matches the presence of any ICMPv6 neighbor discovery option.
	 *
	 * See struct rte_flow_item_icmp6_nd_opt.
	 */
	ItemTypeICMP6NdOpt ItemType = C.RTE_FLOW_ITEM_TYPE_ICMP6_ND_OPT

	/**
	 * Matches an ICMPv6 neighbor discovery source Ethernet link-layer
	 * address option.
	 *
	 * See struct rte_flow_item_icmp6_nd_opt_sla_eth.
	 */
	ItemTypeICMP6NdOptSLAEth ItemType = C.RTE_FLOW_ITEM_TYPE_ICMP6_ND_OPT_SLA_ETH

	/**
	 * Matches an ICMPv6 neighbor discovery target Ethernet link-layer
	 * address option.
	 *
	 * See struct rte_flow_item_icmp6_nd_opt_tla_eth.
	 */
	ItemTypeICMP6NdOptTlaEth ItemType = C.RTE_FLOW_ITEM_TYPE_ICMP6_ND_OPT_TLA_ETH

	/**
	 * Matches specified mark field.
	 *
	 * See struct rte_flow_item_mark.
	 */
	ItemTypeMark ItemType = C.RTE_FLOW_ITEM_TYPE_MARK

	/**
	 * [META]
	 *
	 * Matches a metadata value specified in mbuf metadata field.
	 * See struct rte_flow_item_meta.
	 */
	ItemTypeMeta ItemType = C.RTE_FLOW_ITEM_TYPE_META
)

// Action type constants.
const (
	/**
	 * End marker for action lists. Prevents further processing of
	 * actions, thereby ending the list.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeEnd ActionType = C.RTE_FLOW_ACTION_TYPE_END

	/**
	 * Used as a placeholder for convenience. It is ignored and simply
	 * discarded by PMDs.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeVoid ActionType = C.RTE_FLOW_ACTION_TYPE_VOID

	/**
	 * Leaves traffic up for additional processing by subsequent flow
	 * rules; makes a flow rule non-terminating.
	 *
	 * No associated configuration structure.
	 */
	ActionTypePassthru ActionType = C.RTE_FLOW_ACTION_TYPE_PASSTHRU

	/**
	 * RTE_FLOW_ACTION_TYPE_JUMP
	 *
	 * Redirects packets to a group on the current device.
	 *
	 * See struct rte_flow_action_jump.
	 */
	ActionTypeJump ActionType = C.RTE_FLOW_ACTION_TYPE_JUMP

	/**
	 * Attaches an integer value to packets and sets PKT_RX_FDIR and
	 * PKT_RX_FDIR_ID mbuf flags.
	 *
	 * See struct rte_flow_action_mark.
	 */
	ActionTypeMark ActionType = C.RTE_FLOW_ACTION_TYPE_MARK

	/**
	 * Flags packets. Similar to MARK without a specific value; only
	 * sets the PKT_RX_FDIR mbuf flag.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeFlag ActionType = C.RTE_FLOW_ACTION_TYPE_FLAG

	/**
	 * Assigns packets to a given queue index.
	 *
	 * See struct rte_flow_action_queue.
	 */
	ActionTypeQueue ActionType = C.RTE_FLOW_ACTION_TYPE_QUEUE

	/**
	 * Drops packets.
	 *
	 * PASSTHRU overrides this action if both are specified.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeDrop ActionType = C.RTE_FLOW_ACTION_TYPE_DROP

	/**
	 * Enables counters for this flow rule.
	 *
	 * These counters can be retrieved and reset through rte_flow_query(),
	 * see struct rte_flow_query_count.
	 *
	 * See struct rte_flow_action_count.
	 */
	ActionTypeCount ActionType = C.RTE_FLOW_ACTION_TYPE_COUNT

	/**
	 * Similar to QUEUE, except RSS is additionally performed on packets
	 * to spread them among several queues according to the provided
	 * parameters.
	 *
	 * See struct rte_flow_action_rss.
	 */
	ActionTypeRss ActionType = C.RTE_FLOW_ACTION_TYPE_RSS

	/**
	 * Directs matching traffic to the physical function (PF) of the
	 * current device.
	 *
	 * No associated configuration structure.
	 */
	ActionTypePf ActionType = C.RTE_FLOW_ACTION_TYPE_PF

	/**
	 * Directs matching traffic to a given virtual function of the
	 * current device.
	 *
	 * See struct rte_flow_action_vf.
	 */
	ActionTypeVf ActionType = C.RTE_FLOW_ACTION_TYPE_VF

	/**
	 * Directs packets to a given physical port index of the underlying
	 * device.
	 *
	 * See struct rte_flow_action_phy_port.
	 */
	ActionTypePhyPort ActionType = C.RTE_FLOW_ACTION_TYPE_PHY_PORT

	/**
	 * Directs matching traffic to a given DPDK port ID.
	 *
	 * See struct rte_flow_action_port_id.
	 */
	ActionTypePortID ActionType = C.RTE_FLOW_ACTION_TYPE_PORT_ID

	/**
	 * Traffic metering and policing (MTR).
	 *
	 * See struct rte_flow_action_meter.
	 * See file rte_mtr.h for MTR object configuration.
	 */
	ActionTypeMeter ActionType = C.RTE_FLOW_ACTION_TYPE_METER

	/**
	 * Redirects packets to security engine of current device for security
	 * processing as specified by security session.
	 *
	 * See struct rte_flow_action_security.
	 */
	ActionTypeSecurity ActionType = C.RTE_FLOW_ACTION_TYPE_SECURITY

	/**
	 * Implements OFPAT_SET_MPLS_TTL ("MPLS TTL") as defined by the
	 * OpenFlow Switch Specification.
	 *
	 * See struct rte_flow_action_of_set_mpls_ttl.
	 */
	ActionTypeOfSetMplsTTL ActionType = C.RTE_FLOW_ACTION_TYPE_OF_SET_MPLS_TTL

	/**
	 * Implements OFPAT_DEC_MPLS_TTL ("decrement MPLS TTL") as defined
	 * by the OpenFlow Switch Specification.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeOfDecMplsTTL ActionType = C.RTE_FLOW_ACTION_TYPE_OF_DEC_MPLS_TTL

	/**
	 * Implements OFPAT_SET_NW_TTL ("IP TTL") as defined by the OpenFlow
	 * Switch Specification.
	 *
	 * See struct rte_flow_action_of_set_nw_ttl.
	 */
	ActionTypeOfSetNwTTL ActionType = C.RTE_FLOW_ACTION_TYPE_OF_SET_NW_TTL

	/**
	 * Implements OFPAT_DEC_NW_TTL ("decrement IP TTL") as defined by
	 * the OpenFlow Switch Specification.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeOfDecNwTTL ActionType = C.RTE_FLOW_ACTION_TYPE_OF_DEC_NW_TTL

	/**
	 * Implements OFPAT_COPY_TTL_OUT ("copy TTL "outwards" -- from
	 * next-to-outermost to outermost") as defined by the OpenFlow
	 * Switch Specification.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeOfCopyTTLOut ActionType = C.RTE_FLOW_ACTION_TYPE_OF_COPY_TTL_OUT

	/**
	 * Implements OFPAT_COPY_TTL_IN ("copy TTL "inwards" -- from
	 * outermost to next-to-outermost") as defined by the OpenFlow
	 * Switch Specification.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeOfCopyTTLIn ActionType = C.RTE_FLOW_ACTION_TYPE_OF_COPY_TTL_IN

	/**
	 * Implements OFPAT_POP_VLAN ("pop the outer VLAN tag") as defined
	 * by the OpenFlow Switch Specification.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeOfPopVlan ActionType = C.RTE_FLOW_ACTION_TYPE_OF_POP_VLAN

	/**
	 * Implements OFPAT_PUSH_VLAN ("push a new VLAN tag") as defined by
	 * the OpenFlow Switch Specification.
	 *
	 * See struct rte_flow_action_of_push_vlan.
	 */
	ActionTypeOfPushVlan ActionType = C.RTE_FLOW_ACTION_TYPE_OF_PUSH_VLAN

	/**
	 * Implements OFPAT_SET_VLAN_VID ("set the 802.1q VLAN id") as
	 * defined by the OpenFlow Switch Specification.
	 *
	 * See struct rte_flow_action_of_set_vlan_vid.
	 */
	ActionTypeOfSetVlanVid ActionType = C.RTE_FLOW_ACTION_TYPE_OF_SET_VLAN_VID

	/**
	 * Implements OFPAT_SET_LAN_PCP ("set the 802.1q priority") as
	 * defined by the OpenFlow Switch Specification.
	 *
	 * See struct rte_flow_action_of_set_vlan_pcp.
	 */
	ActionTypeOfSetVlanPcp ActionType = C.RTE_FLOW_ACTION_TYPE_OF_SET_VLAN_PCP

	/**
	 * Implements OFPAT_POP_MPLS ("pop the outer MPLS tag") as defined
	 * by the OpenFlow Switch Specification.
	 *
	 * See struct rte_flow_action_of_pop_mpls.
	 */
	ActionTypeOfPopMpls ActionType = C.RTE_FLOW_ACTION_TYPE_OF_POP_MPLS

	/**
	 * Implements OFPAT_PUSH_MPLS ("push a new MPLS tag") as defined by
	 * the OpenFlow Switch Specification.
	 *
	 * See struct rte_flow_action_of_push_mpls.
	 */
	ActionTypeOfPushMpls ActionType = C.RTE_FLOW_ACTION_TYPE_OF_PUSH_MPLS

	/**
	 * Encapsulate flow in VXLAN tunnel as defined in
	 * rte_flow_action_vxlan_encap action structure.
	 *
	 * See struct rte_flow_action_vxlan_encap.
	 */
	ActionTypeVxlanEncap ActionType = C.RTE_FLOW_ACTION_TYPE_VXLAN_ENCAP

	/**
	 * Decapsulate outer most VXLAN tunnel from matched flow.
	 *
	 * If flow pattern does not define a valid VXLAN tunnel (as specified by
	 * RFC7348) then the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION
	 * error.
	 */
	ActionTypeVxlanDecap ActionType = C.RTE_FLOW_ACTION_TYPE_VXLAN_DECAP

	/**
	 * Encapsulate flow in NVGRE tunnel defined in the
	 * rte_flow_action_nvgre_encap action structure.
	 *
	 * See struct rte_flow_action_nvgre_encap.
	 */
	ActionTypeNvgreEncap ActionType = C.RTE_FLOW_ACTION_TYPE_NVGRE_ENCAP

	/**
	 * Decapsulate outer most NVGRE tunnel from matched flow.
	 *
	 * If flow pattern does not define a valid NVGRE tunnel (as specified by
	 * RFC7637) then the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION
	 * error.
	 */
	ActionTypeNvgreDecap ActionType = C.RTE_FLOW_ACTION_TYPE_NVGRE_DECAP

	/**
	 * Add outer header whose template is provided in its data buffer
	 *
	 * See struct rte_flow_action_raw_encap.
	 */
	ActionTypeRawEncap ActionType = C.RTE_FLOW_ACTION_TYPE_RAW_ENCAP

	/**
	 * Remove outer header whose template is provided in its data buffer.
	 *
	 * See struct rte_flow_action_raw_decap
	 */
	ActionTypeRawDecap ActionType = C.RTE_FLOW_ACTION_TYPE_RAW_DECAP

	/**
	 * Modify IPv4 source address in the outermost IPv4 header.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_IPV4,
	 * then the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_ipv4.
	 */
	ActionTypeSetIpv4Src ActionType = C.RTE_FLOW_ACTION_TYPE_SET_IPV4_SRC

	/**
	 * Modify IPv4 destination address in the outermost IPv4 header.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_IPV4,
	 * then the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_ipv4.
	 */
	ActionTypeSetIpv4Dst ActionType = C.RTE_FLOW_ACTION_TYPE_SET_IPV4_DST

	/**
	 * Modify IPv6 source address in the outermost IPv6 header.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_IPV6,
	 * then the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_ipv6.
	 */
	ActionTypeSetIpv6Src ActionType = C.RTE_FLOW_ACTION_TYPE_SET_IPV6_SRC

	/**
	 * Modify IPv6 destination address in the outermost IPv6 header.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_IPV6,
	 * then the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_ipv6.
	 */
	ActionTypeSetIpv6Dst ActionType = C.RTE_FLOW_ACTION_TYPE_SET_IPV6_DST

	/**
	 * Modify source port number in the outermost TCP/UDP header.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_TCP
	 * or RTE_FLOW_ITEM_TYPE_UDP, then the PMD should return a
	 * RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_tp.
	 */
	ActionTypeSetTpSrc ActionType = C.RTE_FLOW_ACTION_TYPE_SET_TP_SRC

	/**
	 * Modify destination port number in the outermost TCP/UDP header.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_TCP
	 * or RTE_FLOW_ITEM_TYPE_UDP, then the PMD should return a
	 * RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_tp.
	 */
	ActionTypeSetTpDst ActionType = C.RTE_FLOW_ACTION_TYPE_SET_TP_DST

	/**
	 * Swap the source and destination MAC addresses in the outermost
	 * Ethernet header.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_ETH,
	 * then the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * No associated configuration structure.
	 */
	ActionTypeMacSwap ActionType = C.RTE_FLOW_ACTION_TYPE_MAC_SWAP

	/**
	 * Decrease TTL value directly
	 *
	 * No associated configuration structure.
	 */
	ActionTypeDecTTL ActionType = C.RTE_FLOW_ACTION_TYPE_DEC_TTL

	/**
	 * Set TTL value
	 *
	 * See struct rte_flow_action_set_ttl
	 */
	ActionTypeSetTTL ActionType = C.RTE_FLOW_ACTION_TYPE_SET_TTL

	/**
	 * Set source MAC address from matched flow.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_ETH,
	 * the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_mac.
	 */
	ActionTypeSetMacSrc ActionType = C.RTE_FLOW_ACTION_TYPE_SET_MAC_SRC

	/**
	 * Set destination MAC address from matched flow.
	 *
	 * If flow pattern does not define a valid RTE_FLOW_ITEM_TYPE_ETH,
	 * the PMD should return a RTE_FLOW_ERROR_TYPE_ACTION error.
	 *
	 * See struct rte_flow_action_set_mac.
	 */
	ActionTypeSetMacDst ActionType = C.RTE_FLOW_ACTION_TYPE_SET_MAC_DST
)

// HashFunction represents hash functions for RSS.
type HashFunction uint32

// RSS hash function identifiers.
const (
	HashFunctionDefault           HashFunction = C.RTE_ETH_HASH_FUNCTION_DEFAULT            // Default
	HashFunctionToeplitz          HashFunction = C.RTE_ETH_HASH_FUNCTION_TOEPLITZ           // Toeplitz.
	HashFunctionSimpleXor         HashFunction = C.RTE_ETH_HASH_FUNCTION_SIMPLE_XOR         // Simple XOR.
	HashFunctionSymmetricToeplitz HashFunction = C.RTE_ETH_HASH_FUNCTION_SYMMETRIC_TOEPLITZ // symmetrics toeplitz
	HashFunctionMax               HashFunction = C.RTE_ETH_HASH_FUNCTION_MAX
)

// ConvOpType represents conversion operation type.
type ConvOpType uint32

// Conversion operation types.
const (
	/**
	 * No operation to perform.
	 *
	 * rte_flow_conv() simply returns 0.
	 */
	ConvOpNone ConvOpType = C.RTE_FLOW_CONV_OP_NONE

	/**
	 * Convert attributes structure.
	 *
	 * This is a basic copy of an attributes structure.
	 *
	 * - @p src type:
	 *   @code const struct rte_flow_attr * @endcode
	 * - @p dst type:
	 *   @code struct rte_flow_attr * @endcode
	 */
	ConvOpAttr ConvOpType = C.RTE_FLOW_CONV_OP_ATTR

	/**
	 * Convert a single item.
	 *
	 * Duplicates @p spec, @p last and @p mask but not outside objects.
	 *
	 * - @p src type:
	 *   @code const struct rte_flow_item * @endcode
	 * - @p dst type:
	 *   @code struct rte_flow_item * @endcode
	 */
	ConvOpItem ConvOpType = C.RTE_FLOW_CONV_OP_ITEM

	/**
	 * Convert a single action.
	 *
	 * Duplicates @p conf but not outside objects.
	 *
	 * - @p src type:
	 *   @code const struct rte_flow_action * @endcode
	 * - @p dst type:
	 *   @code struct rte_flow_action * @endcode
	 */
	ConvOpAction ConvOpType = C.RTE_FLOW_CONV_OP_ACTION

	/**
	 * Convert an entire pattern.
	 *
	 * Duplicates all pattern items at once with the same constraints as
	 * RTE_FLOW_CONV_OP_ITEM.
	 *
	 * - @p src type:
	 *   @code const struct rte_flow_item * @endcode
	 * - @p dst type:
	 *   @code struct rte_flow_item * @endcode
	 */
	ConvOpPattern ConvOpType = C.RTE_FLOW_CONV_OP_PATTERN

	/**
	 * Convert a list of actions.
	 *
	 * Duplicates the entire list of actions at once with the same
	 * constraints as RTE_FLOW_CONV_OP_ACTION.
	 *
	 * - @p src type:
	 *   @code const struct rte_flow_action * @endcode
	 * - @p dst type:
	 *   @code struct rte_flow_action * @endcode
	 */
	ConvOpActions ConvOpType = C.RTE_FLOW_CONV_OP_ACTIONS

	/**
	 * Convert a complete flow rule description.
	 *
	 * Comprises attributes, pattern and actions together at once with
	 * the usual constraints.
	 *
	 * - @p src type:
	 *   @code const struct rte_flow_conv_rule * @endcode
	 * - @p dst type:
	 *   @code struct rte_flow_conv_rule * @endcode
	 */
	ConvOpRule ConvOpType = C.RTE_FLOW_CONV_OP_RULE

	/**
	 * Convert item type to its name string.
	 *
	 * Writes a NUL-terminated string to @p dst. Like snprintf(), the
	 * returned value excludes the terminator which is always written
	 * nonetheless.
	 *
	 * - @p src type:
	 *   @code (const void *)enum rte_flow_item_type @endcode
	 * - @p dst type:
	 *   @code char * @endcode
	 **/
	ConvOpItemName ConvOpType = C.RTE_FLOW_CONV_OP_ITEM_NAME

	/**
	 * Convert action type to its name string.
	 *
	 * Writes a NUL-terminated string to @p dst. Like snprintf(), the
	 * returned value excludes the terminator which is always written
	 * nonetheless.
	 *
	 * - @p src type:
	 *   @code (const void *)enum rte_flow_action_type @endcode
	 * - @p dst type:
	 *   @code char * @endcode
	 **/
	ConvOpActionName ConvOpType = C.RTE_FLOW_CONV_OP_ACTION_NAME

	/**
	 * Convert item type to pointer to item name.
	 *
	 * Retrieves item name pointer from its type. The string itself is
	 * not copied; instead, a unique pointer to an internal static
	 * constant storage is written to @p dst.
	 *
	 * - @p src type:
	 *   @code (const void *)enum rte_flow_item_type @endcode
	 * - @p dst type:
	 *   @code const char ** @endcode
	 */
	ConvOpItemNamePtr ConvOpType = C.RTE_FLOW_CONV_OP_ITEM_NAME_PTR

	/**
	 * Convert action type to pointer to action name.
	 *
	 * Retrieves action name pointer from its type. The string itself is
	 * not copied; instead, a unique pointer to an internal static
	 * constant storage is written to @p dst.
	 *
	 * - @p src type:
	 *   @code (const void *)enum rte_flow_action_type @endcode
	 * - @p dst type:
	 *   @code const char ** @endcode
	 */
	ConvOpActionNamePtr ConvOpType = C.RTE_FLOW_CONV_OP_ACTION_NAME_PTR
)
