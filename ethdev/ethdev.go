/*
Package ethdev wraps RTE Ethernet Device API.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package ethdev

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_ring.h>
#include <rte_errno.h>
#include <rte_memory.h>
#include <rte_ethdev.h>

static void set_tx_reject_tagged(struct rte_eth_txmode *txm) {
	txm->hw_vlan_reject_tagged = 1;
}

static void set_tx_reject_untagged(struct rte_eth_txmode *txm) {
	txm->hw_vlan_reject_untagged = 1;
}

static void set_tx_insert_pvid(struct rte_eth_txmode *txm) {
	txm->hw_vlan_insert_pvid = 1;
}
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/mempool"
)

// Various RX offloads flags.
const (
	DevRxOffloadVlanStrip      uint64 = C.DEV_RX_OFFLOAD_VLAN_STRIP
	DevRxOffloadIpv4Cksum             = C.DEV_RX_OFFLOAD_IPV4_CKSUM
	DevRxOffloadUdpCksum              = C.DEV_RX_OFFLOAD_UDP_CKSUM
	DevRxOffloadTcpCksum              = C.DEV_RX_OFFLOAD_TCP_CKSUM
	DevRxOffloadTcpLro                = C.DEV_RX_OFFLOAD_TCP_LRO
	DevRxOffloadQinqStrip             = C.DEV_RX_OFFLOAD_QINQ_STRIP
	DevRxOffloadOuterIpv4Cksum        = C.DEV_RX_OFFLOAD_OUTER_IPV4_CKSUM
	DevRxOffloadMacsecStrip           = C.DEV_RX_OFFLOAD_MACSEC_STRIP
	DevRxOffloadHeaderSplit           = C.DEV_RX_OFFLOAD_HEADER_SPLIT
	DevRxOffloadVlanFilter            = C.DEV_RX_OFFLOAD_VLAN_FILTER
	DevRxOffloadVlanExtend            = C.DEV_RX_OFFLOAD_VLAN_EXTEND
	DevRxOffloadJumboFrame            = C.DEV_RX_OFFLOAD_JUMBO_FRAME
	DevRxOffloadScatter               = C.DEV_RX_OFFLOAD_SCATTER
	DevRxOffloadTimestamp             = C.DEV_RX_OFFLOAD_TIMESTAMP
	DevRxOffloadSecurity              = C.DEV_RX_OFFLOAD_SECURITY
	// DevRxOffloadKeepCrc        = C.DEV_RX_OFFLOAD_KEEP_CRC
	// DevRxOffloadSctpCksum      = C.DEV_RX_OFFLOAD_SCTP_CKSUM
	// DevRxOffloadOuterUdpCksum  = C.DEV_RX_OFFLOAD_OUTER_UDP_CKSUM

	DevRxOffloadChecksum = (DevRxOffloadIpv4Cksum |
		DevRxOffloadUdpCksum |
		DevRxOffloadTcpCksum)
	DevRxOffloadVlan = (DevRxOffloadVlanStrip |
		DevRxOffloadVlanFilter |
		DevRxOffloadVlanExtend)
)

// Various TX offloads flags.
const (
	DevTxOffloadVlanInsert     uint64 = C.DEV_TX_OFFLOAD_VLAN_INSERT
	DevTxOffloadIpv4Cksum             = C.DEV_TX_OFFLOAD_IPV4_CKSUM
	DevTxOffloadUdpCksum              = C.DEV_TX_OFFLOAD_UDP_CKSUM
	DevTxOffloadTcpCksum              = C.DEV_TX_OFFLOAD_TCP_CKSUM
	DevTxOffloadSctpCksum             = C.DEV_TX_OFFLOAD_SCTP_CKSUM
	DevTxOffloadTcpTso                = C.DEV_TX_OFFLOAD_TCP_TSO
	DevTxOffloadUdpTso                = C.DEV_TX_OFFLOAD_UDP_TSO
	DevTxOffloadOuterIpv4Cksum        = C.DEV_TX_OFFLOAD_OUTER_IPV4_CKSUM
	DevTxOffloadQinqInsert            = C.DEV_TX_OFFLOAD_QINQ_INSERT
	DevTxOffloadVxlanTnlTso           = C.DEV_TX_OFFLOAD_VXLAN_TNL_TSO
	DevTxOffloadGreTnlTso             = C.DEV_TX_OFFLOAD_GRE_TNL_TSO
	DevTxOffloadIpipTnlTso            = C.DEV_TX_OFFLOAD_IPIP_TNL_TSO
	DevTxOffloadGeneveTnlTso          = C.DEV_TX_OFFLOAD_GENEVE_TNL_TSO
	DevTxOffloadMacsecInsert          = C.DEV_TX_OFFLOAD_MACSEC_INSERT
	DevTxOffloadMtLockfree            = C.DEV_TX_OFFLOAD_MT_LOCKFREE
	DevTxOffloadMultiSegs             = C.DEV_TX_OFFLOAD_MULTI_SEGS
	DevTxOffloadMbufFastFree          = C.DEV_TX_OFFLOAD_MBUF_FAST_FREE
	DevTxOffloadSecurity              = C.DEV_TX_OFFLOAD_SECURITY
	// DevTxOffloadIpTnlTso       = C.DEV_TX_OFFLOAD_IP_TNL_TSO
	// DevTxOffloadOuterUdpCksum  = C.DEV_TX_OFFLOAD_OUTER_UDP_CKSUM
	// DevTxOffloadMatchMetadata  = C.DEV_TX_OFFLOAD_MATCH_METADATA
)

// Device supported speeds bitmap flags.
const (
	EthLinkSpeedAutoneg uint = C.ETH_LINK_SPEED_AUTONEG /**< Autonegotiate (all speeds) */
	EthLinkSpeedFixed        = C.ETH_LINK_SPEED_FIXED   /**< Disable autoneg (fixed speed) */
	EthLinkSpeed10mHd        = C.ETH_LINK_SPEED_10M_HD  /**<  10 Mbps half-duplex */
	EthLinkSpeed10m          = C.ETH_LINK_SPEED_10M     /**<  10 Mbps full-duplex */
	EthLinkSpeed100mHd       = C.ETH_LINK_SPEED_100M_HD /**< 100 Mbps half-duplex */
	EthLinkSpeed100m         = C.ETH_LINK_SPEED_100M    /**< 100 Mbps full-duplex */
	EthLinkSpeed1g           = C.ETH_LINK_SPEED_1G      /**<   1 Gbps */
	EthLinkSpeed2_5g         = C.ETH_LINK_SPEED_2_5G    /**< 2.5 Gbps */
	EthLinkSpeed5g           = C.ETH_LINK_SPEED_5G      /**<   5 Gbps */
	EthLinkSpeed10g          = C.ETH_LINK_SPEED_10G     /**<  10 Gbps */
	EthLinkSpeed20g          = C.ETH_LINK_SPEED_20G     /**<  20 Gbps */
	EthLinkSpeed25g          = C.ETH_LINK_SPEED_25G     /**<  25 Gbps */
	EthLinkSpeed40g          = C.ETH_LINK_SPEED_40G     /**<  40 Gbps */
	EthLinkSpeed50g          = C.ETH_LINK_SPEED_50G     /**<  50 Gbps */
	EthLinkSpeed56g          = C.ETH_LINK_SPEED_56G     /**<  56 Gbps */
	EthLinkSpeed100g         = C.ETH_LINK_SPEED_100G    /**< 100 Gbps */
)

// A set of values to identify what method is to be used to route
// packets to multiple queues.
const (
	EthMqRxNone       uint = C.ETH_MQ_RX_NONE         /** None of DCB,RSS or VMDQ mode */
	EthMqRxRss             = C.ETH_MQ_RX_RSS          /** For RX side, only RSS is on */
	EthMqRxDcb             = C.ETH_MQ_RX_DCB          /** For RX side,only DCB is on. */
	EthMqRxDcbRss          = C.ETH_MQ_RX_DCB_RSS      /** Both DCB and RSS enable */
	EthMqRxVmdqOnly        = C.ETH_MQ_RX_VMDQ_ONLY    /** Only VMDQ, no RSS nor DCB */
	EthMqRxVmdqRss         = C.ETH_MQ_RX_VMDQ_RSS     /** RSS mode with VMDQ */
	EthMqRxVmdqDcb         = C.ETH_MQ_RX_VMDQ_DCB     /** Use VMDQ+DCB to route traffic to queues */
	EthMqRxVmdqDcbRss      = C.ETH_MQ_RX_VMDQ_DCB_RSS /** Enable both VMDQ and DCB in VMDq */
)

// A set of values to identify what method is to be used to transmit
// packets using multi-TCs.
const (
	EthMqTxNone     uint = C.ETH_MQ_TX_NONE      /**< It is in neither DCB nor VT mode. */
	EthMqTxDcb           = C.ETH_MQ_TX_DCB       /**< For TX side,only DCB is on. */
	EthMqTxVmdqDcb       = C.ETH_MQ_TX_VMDQ_DCB  /**< For TX side,both DCB and VT is on. */
	EthMqTxVmdqOnly      = C.ETH_MQ_TX_VMDQ_ONLY /**< Only VT on, no DCB */
)

// reinterpret C.int error.
// C.int is not C.int in other packages, hence...
func errno(n C.int) error {
	return common.Errno(int(n))
}

// DevOption represents device option which is then used by
// DevConfigure to setup Ethernet device.
type DevOption struct {
	f func(*C.struct_rte_eth_conf)
}

// configuration options for RX queue
type rxqConf struct {
	socket C.int
	mp     *mempool.Mempool
	rx     C.struct_rte_eth_rxconf
}

// configuration options for TX queue
type txqConf struct {
	socket C.int
	tx     C.struct_rte_eth_txconf
}

// RxOption represents an option which is used to setup RX queue on
// Ethernet device.
type RxOption struct {
	f func(*rxqConf)
}

// TxOption represents an option which is used to setup TX queue on
// Ethernet device.
type TxOption struct {
	f func(*txqConf)
}

// EthRxMode is used to configure Ethernet device through
// DevOptEthRxMode option.
type EthRxMode struct {
	// The multi-queue packet distribution mode to be used, e.g. RSS.
	// See EthMqRx* constants.
	MqMode uint
	// Only used if JUMBO_FRAME enabled.
	MaxRxPktLen uint32
	// hdr buf size (header_split enabled).
	SplitHdrSize uint16
	// Per-port Rx offloads to be set using DevRxOffload* flags. Only
	// offloads set on rx_offload_capa field on rte_eth_dev_info
	// structure are allowed to be set.
	Offloads uint64
}

// EthRxMode is used to configure Ethernet device through
// DevOptEthTxMode option.
type EthTxMode struct {
	// TX multi-queues mode.
	MqMode uint
	// Per-port Tx offloads to be set using DevTxOffload*
	// flags. Only offloads set on tx_offload_capa field on
	// rte_eth_dev_info structure are allowed to be set.
	Offloads uint64
	// For i40e specifically.
	Pvid uint16
	// If set, reject sending out tagged pkts.
	HwVlanRejectTagged bool
	// If set, reject sending out untagged pkts
	HwVlanRejectUntagged bool
	// If set, enable port based VLAN insertion
	HwVlanInsertPvid bool
}

// A structure used to configure the Receive Side Scaling (RSS)
// feature of an Ethernet port.  If not nil, the Key points to an
// array holding the RSS key to use for hashing specific header fields
// of received packets.  Otherwise, a default random hash key is used
// by the device driver.
//
// To maintain compatibility the Key should be 40 bytes long.  To be
// compatible, this length will be checked in i40e only. Others assume
// 40 bytes to be used as before.
//
// The Hf field indicates the different types of IPv4/IPv6 packets to
// which the RSS hashing must be applied.  Supplying an *rss_hf* equal
// to zero disables the RSS feature.
type EthRssConf struct {
	/**< If not NULL, 40-byte hash key. */
	Key []byte
	/**< Hash functions to apply. */
	Hf uint64
}

// A structure used to configure the ring threshold registers of an RX/TX queue
// for an Ethernet port.
type EthThresh struct {
	PThresh uint8 /**< Ring prefetch threshold. */
	HThresh uint8 /**< Ring host threshold. */
	WThresh uint8 /**< Ring writeback threshold. */
}

// A structure used to configure an RX ring of an Ethernet port.
type EthRxqConf struct {
	// Drives the freeing of RX descriptors.
	FreeThresh uint16
	// Drop packets if no descriptors are available.
	DropEn uint8
	// Do not start queue with rte_eth_dev_start().
	DeferredStart uint8
	// Per-queue Rx offloads to be set using DevRxOffload* flags.
	// Only offloads set on rx_queue_offload_capa or rx_offload_capa
	// fields on rte_eth_dev_info structure are allowed to be set.
	Offloads uint64
}

// A structure used to configure a TX ring of an Ethernet port.
type EthTxqConf struct {
	// Drives the setting of RS bit on TXDs.
	RsThresh uint16
	// Start freeing TX buffers if there are less free descriptors than this value.
	FreeThresh uint16
	// Do not start queue with rte_eth_dev_start().
	DeferredStart uint8
	// Per-queue Tx offloads to be set using DevTxOffload* flags. Only
	// offloads set on tx_queue_offload_capa or tx_offload_capa fields
	// on rte_eth_dev_info structure are allowed to be set.
	Offloads uint64
}

// Port is the number of the Ethernet device.
type Port uint16

// DevOptLinkSpeeds sets allowed speeds for the device.
// EthLinkSpeedFixed disables link autonegotiation, and a unique speed
// shall be set. Otherwise, the bitmap defines the set of speeds to be
// advertised. If the special value EthLinkSpeedAutoneg is used, all
// speeds supported are advertised.
func DevOptLinkSpeeds(speeds uint) DevOption {
	return DevOption{func(ec *C.struct_rte_eth_conf) {
		ec.link_speeds = C.uint(speeds)
	}}
}

// DevOptEthRxMode specifies port RX configuration.
func DevOptEthRxMode(conf EthRxMode) DevOption {
	return DevOption{func(ec *C.struct_rte_eth_conf) {
		ec.rxmode = C.struct_rte_eth_rxmode{
			mq_mode:        uint32(conf.MqMode),
			max_rx_pkt_len: C.uint(conf.MaxRxPktLen),
			split_hdr_size: C.ushort(conf.SplitHdrSize),
			offloads:       C.ulong(conf.Offloads),
		}
	}}
}

// DevOptEthTxMode specifies port TX configuration.
func DevOptEthTxMode(conf EthTxMode) DevOption {
	return DevOption{func(ec *C.struct_rte_eth_conf) {
		ec.txmode = C.struct_rte_eth_txmode{
			mq_mode:  uint32(conf.MqMode),
			offloads: C.ulong(conf.Offloads),
			pvid:     C.ushort(conf.Pvid),
		}
		if conf.HwVlanRejectTagged {
			C.set_tx_reject_tagged(&ec.txmode)
		}
		if conf.HwVlanRejectUntagged {
			C.set_tx_reject_untagged(&ec.txmode)
		}
		if conf.HwVlanInsertPvid {
			C.set_tx_insert_pvid(&ec.txmode)
		}
	}}
}

// DevOptLoopbackMode specifies loopback operation mode. By default
// the value is 0, meaning the loopback mode is disabled.  Read the
// datasheet of given ethernet controller for details. The possible
// values of this field are defined in implementation of each driver.
func DevOptLoopbackMode(mode uint32) DevOption {
	return DevOption{func(ec *C.struct_rte_eth_conf) {
		ec.lpbk_mode = C.uint(mode)
	}}
}

// DevOptRss specifies RSS configuration.
func DevOptRss(conf EthRssConf) DevOption {
	return DevOption{func(ec *C.struct_rte_eth_conf) {
		rssKey := (*C.uchar)(nil)
		if conf.Key != nil && len(conf.Key) > 0 {
			rssKey = (*C.uchar)(unsafe.Pointer(&conf.Key[0]))
		}
		ec.rx_adv_conf.rss_conf = C.struct_rte_eth_rss_conf{
			rss_key:     rssKey,
			rss_key_len: C.uchar(len(conf.Key)),
			rss_hf:      C.ulong(conf.Hf),
		}
	}}
}

// DevConfigure configures an Ethernet device. This function must be
// invoked first before any other function in the Ethernet API. This
// function can also be re-invoked when a device is in the stopped
// state.
//
// nrxq and ntxq are the numbers of receive and transmit queues to set
// up for the Ethernet device, respectively.
//
// Several DevOpt* options may be specified as well.
func (pid Port) DevConfigure(nrxq, ntxq uint16, opts ...DevOption) error {
	conf := &C.struct_rte_eth_conf{}
	for i := range opts {
		opts[i].f(conf)
	}

	return errno(C.rte_eth_dev_configure(C.ushort(pid), C.ushort(nrxq), C.ushort(nrxq), conf))
}

// RxOptConf specifies the configuration an RX ring of an Ethernet
// port.
func RxOptConf(conf EthRxqConf) RxOption {
	return RxOption{func(q *rxqConf) {
		q.rx = C.struct_rte_eth_rxconf{
			rx_free_thresh:    C.ushort(conf.FreeThresh),
			rx_drop_en:        C.uchar(conf.DropEn),
			rx_deferred_start: C.uchar(conf.DeferredStart),
			offloads:          C.ulong(conf.Offloads),
		}
	}}
}

// RxOptThresh specifies the configuration of RX ring threshold
// registers.
func RxOptThresh(conf EthThresh) RxOption {
	return RxOption{func(q *rxqConf) {
		q.rx.rx_thresh = C.struct_rte_eth_thresh{
			pthresh: C.uchar(conf.PThresh),
			hthresh: C.uchar(conf.HThresh),
			wthresh: C.uchar(conf.WThresh),
		}
	}}
}

func RxOptMempool(mp *mempool.Mempool) RxOption {
	return RxOption{func(q *rxqConf) {
		q.mp = mp
	}}
}

// TxqOptConf specifies the configuration an TX ring of an Ethernet
// port.
func TxqOptConf(conf EthTxqConf) TxOption {
	return TxOption{func(q *txqConf) {
		q.tx = C.struct_rte_eth_txconf{
			tx_rs_thresh:      C.ushort(conf.RsThresh),
			tx_free_thresh:    C.ushort(conf.FreeThresh),
			tx_deferred_start: C.uchar(conf.DeferredStart),
			offloads:          C.ulong(conf.Offloads),
		}
	}}
}

// TxqOptThresh specifies the configuration of TX ring threshold
// registers.
func TxqOptThresh(conf EthThresh) TxOption {
	return TxOption{func(q *txqConf) {
		q.tx.tx_thresh = C.struct_rte_eth_thresh{
			pthresh: C.uchar(conf.PThresh),
			hthresh: C.uchar(conf.HThresh),
			wthresh: C.uchar(conf.WThresh),
		}
	}}
}

// RxOptSocket specifies the NUMA socket id for RX.  The socket
// argument is the socket identifier in case of NUMA.  The value can
// be SOCKET_ID_ANY if there is no NUMA constraint for the DMA memory
// allocated for the receive descriptors of the ring.
func RxOptSocket(socket int) RxOption {
	return RxOption{func(q *rxqConf) {
		q.socket = C.int(socket)
	}}
}

// TxOptSocket specifies the NUMA socket id for TX.  The socket
// argument is the socket identifier in case of NUMA.  The value can
// be SOCKET_ID_ANY if there is no NUMA constraint for the DMA memory
// allocated for the transmit descriptors of the ring.
func TxOptSocket(socket int) TxOption {
	return TxOption{func(q *txqConf) {
		q.socket = C.int(socket)
	}}
}

// RxqSetup allocates and sets up a receive queue for an Ethernet
// device.
//
// RxOptMempool is mandatory option.
//
// The function allocates a contiguous block of memory for nDesc
// receive descriptors from a memory zone associated with *socket_id*
// and initializes each receive descriptor with a network buffer
// allocated from the memory pool *mb_pool*.
//
// qid is the index of the receive queue to set up. The value must be
// in the range [0, nb_rx_queue - 1] previously supplied to
// rte_eth_dev_configure().
//
// nDesc is the number of receive descriptors to allocate for the
// receive ring.
//
// mp is the pointer to the memory pool from which to allocate
// *rte_mbuf* network memory buffers to populate each descriptor of
// the receive ring.

// RxOptSocket specifies the socket identifier in case of NUMA.  The
// value can be *SOCKET_ID_ANY* if there is no NUMA constraint for the
// DMA memory allocated for the receive descriptors of the ring.
//
// RxOptConf specifies the configuration data to be used for the
// receive queue.  The *rx_conf* structure contains an *rx_thresh*
// structure with the values of the Prefetch, Host, and Write-Back
// threshold registers of the receive ring.  In addition it contains
// the hardware offloads features to activate using the
// DEV_RX_OFFLOAD_* flags.  If an offloading set in rx_conf->offloads
// hasn't been set in the input argument eth_conf->rxmode.offloads to
// rte_eth_dev_configure(), it is a new added offloading, it must be
// per-queue type and it is enabled for the queue.  No need to repeat
// any bit in rx_conf->offloads which has already been enabled in
// rte_eth_dev_configure() at port level. An offloading enabled at
// port level can't be disabled at queue level.
//
// Return codes:
//
// - 0: Success, receive queue correctly set up.
//
// - -EIO: if device is removed.
//
// - -EINVAL: The memory pool pointer is null or the size of network
//    buffers which can be allocated from this memory pool does not
//    fit the various buffer sizes allowed by the device controller.
//
// - -ENOMEM: Unable to allocate the receive ring descriptors or to
//    allocate network memory buffers from the memory pool when
//    initializing receive descriptors.
func (pid Port) RxqSetup(qid, nDesc uint16, mp *mempool.Mempool, opts ...RxOption) error {
	conf := &rxqConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	return errno(C.rte_eth_rx_queue_setup(C.ushort(pid), C.ushort(qid),
		C.ushort(nDesc), C.uint(conf.socket), &conf.rx,
		(*C.struct_rte_mempool)(unsafe.Pointer(mp))))
}

// TxqSetup allocates and set up a transmit queue for an Ethernet
// device.
//
// qid is the index of the transmit queue to set up.  The value must
// be in the range [0, nb_tx_queue - 1] previously supplied to
// rte_eth_dev_configure().
//
// nDesc is the number of transmit descriptors to allocate for the
// transmit ring.
//
// TxOptSocket specifies the socket identifier in case of NUMA.
// Its value can be *SOCKET_ID_ANY* if there is no NUMA constraint for
// the DMA memory allocated for the transmit descriptors of the ring.
//
// TxqOptConf specifies configuration data to be used for the transmit
// queue.  NULL value is allowed, in which case default TX
// configuration will be used.
//
// The *tx_conf* structure contains the following data:
//
// - The *tx_thresh* structure with the values of the Prefetch, Host,
// and Write-Back threshold registers of the transmit ring.  When
// setting Write-Back threshold to the value greater then zero,
// *tx_rs_thresh* value should be explicitly set to one.
//
// - The *tx_free_thresh* value indicates the [minimum] number of
// network buffers that must be pending in the transmit ring to
// trigger their [implicit] freeing by the driver transmit function.
//
// - The *tx_rs_thresh* value indicates the [minimum] number of
// transmit descriptors that must be pending in the transmit ring
// before setting the RS bit on a descriptor by the driver transmit
// function.  The *tx_rs_thresh* value should be less or equal then
// *tx_free_thresh* value, and both of them should be less then
// *nb_tx_desc* - 3.
//
// - The *offloads* member contains Tx offloads to be enabled.  If an
// offloading set in tx_conf->offloads hasn't been set in the input
// argument eth_conf->txmode.offloads to rte_eth_dev_configure(), it
// is a new added offloading, it must be per-queue type and it is
// enabled for the queue.  No need to repeat any bit in
// tx_conf->offloads which has already been enabled in
// rte_eth_dev_configure() at port level. An offloading enabled at
// port level can't be disabled at queue level.
//
// Note that setting *tx_free_thresh* or *tx_rs_thresh* value to 0
// forces the transmit function to use default values.
// Return codes:
//
// - 0: Success, the transmit queue is correctly set up.
//
// - -ENOMEM: Unable to allocate the transmit ring descriptors.
func (pid Port) TxqSetup(qid, nDesc uint16, opts ...TxOption) error {
	conf := &txqConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	return errno(C.rte_eth_tx_queue_setup(C.ushort(pid), C.ushort(qid),
		C.ushort(nDesc), C.uint(conf.socket), &conf.tx))
}

// Reset a Ethernet device and keep its port id.
//
// When a port has to be reset passively, the DPDK application can
// invoke this function. For example when a PF is reset, all its VFs
// should also be reset. Normally a DPDK application can invoke this
// function when RTE_ETH_EVENT_INTR_RESET event is detected, but can
// also use it to start a port reset in other circumstances.
//
// When this function is called, it first stops the port and then
// calls the PMD specific dev_uninit( ) and dev_init( ) to return the
// port to initial state, in which no Tx and Rx queues are setup, as
// if the port has been reset and not started. The port keeps the port
// id it had before the function call.
//
// After calling rte_eth_dev_reset( ), the application should use
// rte_eth_dev_configure( ), rte_eth_rx_queue_setup( ),
// rte_eth_tx_queue_setup( ), and rte_eth_dev_start( ) to reconfigure
// the device as appropriate.
//
// Note: To avoid unexpected behavior, the application should stop
// calling Tx and Rx functions before calling rte_eth_dev_reset( ).
// For thread safety, all these controlling functions should be called
// from the same thread.
//
// Return codes:
//
//   - (0) if successful.
//
//   - (-EINVAL) if port identifier is invalid.
//
//   - (-ENOTSUP) if hardware doesn't support this function.
//
//   - (-EPERM) if not ran from the primary process.
//
//   - (-EIO) if re-initialisation failed or device is removed.
//
//   - (-ENOMEM) if the reset failed due to OOM.
//
//   - (-EAGAIN) if the reset temporarily failed and should be retried later.
func (pid Port) Reset() error {
	return errno(C.rte_eth_dev_reset(C.ushort(pid)))
}

// Start an Ethernet device.
//
// The device start step is the last one and consists of setting the
// configured offload features and in starting the transmit and the
// receive units of the device.
//
// Device RTE_ETH_DEV_NOLIVE_MAC_ADDR flag causes MAC address to be
// set before PMD port start callback function is invoked.
//
// On success, all basic functions exported by the Ethernet API (link
// status, receive/transmit, and so on) can be invoked.
//
// Return codes:
//
// - 0: Success, Ethernet device started.
//
// - <0: Error code of the driver device start function.
func (pid Port) Start() error {
	return errno(C.rte_eth_dev_start(C.ushort(pid)))
}

// Stop an Ethernet device. The device can be restarted with a call to
// rte_eth_dev_start().
func (pid Port) Stop() {
	C.rte_eth_dev_stop(C.ushort(pid))
}

// Close a stopped Ethernet device. The device cannot be restarted!
// The function frees all port resources if the driver supports
// the flag RTE_ETH_DEV_CLOSE_REMOVE.
func (pid Port) Close() {
	C.rte_eth_dev_close(C.ushort(pid))
}

// PromiscEnable enables receipt in promiscuous mode for an Ethernet
// device.
func (pid Port) PromiscEnable() {
	C.rte_eth_promiscuous_enable(C.ushort(pid))
}

// PromiscDisable disables receipt in promiscuous mode for an Ethernet
// device.
func (pid Port) PromiscDisable() {
	C.rte_eth_promiscuous_disable(C.ushort(pid))
}

// SetLinkUp links up an Ethernet device.
//
// Set device link up will re-enable the device rx/tx
// functionality after it is previously set device linked down.
//
// Return codes:
//
//   - 0: Success, Ethernet device linked up.
//
//   - <0: Error code of the driver device link up function.
func (pid Port) SetLinkUp() error {
	return errno(C.rte_eth_dev_set_link_up(C.ushort(pid)))
}

// Link down an Ethernet device.
// The device rx/tx functionality will be disabled if success,
// and it can be re-enabled with a call to
// rte_eth_dev_set_link_up().
//
// Return codes:
//
//   - 0: Success, Ethernet device linked down.
//
//   - <0: Error code of the driver device link down function.
func (pid Port) SetLinkDown() error {
	return errno(C.rte_eth_dev_set_link_down(C.ushort(pid)))
}
