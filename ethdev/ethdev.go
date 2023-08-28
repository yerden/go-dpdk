/*
Package ethdev wraps RTE Ethernet Device API.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package ethdev

/*
#include <stdlib.h>
#include <net/if.h>

#include <rte_config.h>
#include <rte_errno.h>
#include <rte_memory.h>
#include <rte_ethdev.h>
#include <rte_version.h>
#include <rte_eth_ring.h>

// The max_rx_pkt_len changes occurred in commit: 1bb4a528c41f4af4847bd3d58cc2b2b9f1ec9a27.
#if RTE_VERSION < RTE_VERSION_NUM(21, 11, 0, 0)
enum {
	RX_MODE_LEN_OFF = offsetof(struct rte_eth_rxmode, max_rx_pkt_len),
};
#else
enum {
	RX_MODE_LEN_OFF = offsetof(struct rte_eth_rxmode, mtu),
};
#endif

static void set_tx_reject_tagged(struct rte_eth_txmode *txm) {
	txm->hw_vlan_reject_tagged = 1;
}

static void set_tx_reject_untagged(struct rte_eth_txmode *txm) {
	txm->hw_vlan_reject_untagged = 1;
}

static void set_tx_insert_pvid(struct rte_eth_txmode *txm) {
	txm->hw_vlan_insert_pvid = 1;
}

static int go_rte_get_ethdev_ports(uint16_t *ports, int n_ports) {
	uint16_t pid;
	int i = 0;

	RTE_ETH_FOREACH_DEV(pid) {
		if (i >= n_ports) {
			break;
		}
		ports[i++] = pid;
	}

	return i;
}

struct go_rte_eth_link {
	uint32_t link_speed;
	uint8_t link_duplex;
	uint8_t link_autoneg;
	uint8_t link_status;
};

static int go_rte_eth_link_get(uint16_t port_id, struct go_rte_eth_link *link)
{
	struct rte_eth_link data;
	int rc = rte_eth_link_get(port_id, &data);
	link->link_speed = data.link_speed;
	link->link_duplex= data.link_duplex;
	link->link_autoneg= data.link_autoneg;
	link->link_status= data.link_status;
	return rc;
}

static int go_rte_eth_link_get_nowait(uint16_t port_id, struct go_rte_eth_link *link)
{
	struct rte_eth_link data;
	int rc = rte_eth_link_get_nowait(port_id, &data);
	link->link_speed = data.link_speed;
	link->link_duplex= data.link_duplex;
	link->link_autoneg= data.link_autoneg;
	link->link_status= data.link_status;
	return rc;
}

*/
import "C"

import (
	"net"
	"reflect"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/mempool"
	"github.com/yerden/go-dpdk/ring"
)

// This enum indicates the flow control mode.
const (
	// Disable flow control.
	FcNone uint32 = C.RTE_FC_NONE
	// RX pause frame, enable flowctrl on TX side.
	FcRxPause uint32 = C.RTE_FC_RX_PAUSE
	// TX pause frame, enable flowctrl on RX side.
	FcTxPause uint32 = C.RTE_FC_TX_PAUSE
	// Enable flow control on both side.
	FcFull uint32 = C.RTE_FC_FULL
)

// Option represents device option which is then used by
// DevConfigure to setup Ethernet device.
type Option struct {
	f func(*ethConf)
}

// configuration options for RX/TX queue
type qConf struct {
	socket C.int
	rx     C.struct_rte_eth_rxconf
	tx     C.struct_rte_eth_txconf
}

// QueueOption represents an option which is used to setup RX/TX queue on
// Ethernet device.
type QueueOption struct {
	f func(*qConf)
}

// RxMode is used to configure Ethernet device through
// OptRxMode option.
type RxMode struct {
	// The multi-queue packet distribution mode to be used, e.g. RSS.
	// See MqRx* constants.
	MqMode uint
	// Requested MTU or MaxRxPktLen if JUMBO_FRAME enabled (for releases older than 21.11).
	MTU uint32
	// hdr buf size (header_split enabled).
	SplitHdrSize uint16
	// Per-port Rx offloads to be set using RxOffload* flags. Only
	// offloads set on rx_offload_capa field on rte_eth_dev_info
	// structure are allowed to be set.
	Offloads uint64
}

// TxMode is used to configure Ethernet device through
// OptTxMode option.
type TxMode struct {
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

// RssConf is a structure used to configure the Receive Side Scaling
// (RSS) feature of an Ethernet port.  If not nil, the Key points to
// an array holding the RSS key to use for hashing specific header
// fields of received packets.  Otherwise, a default random hash key
// is used by the device driver.
//
// To maintain compatibility the Key should be 40 bytes long.  To be
// compatible, this length will be checked in i40e only. Others assume
// 40 bytes to be used as before.
//
// The Hf field indicates the different types of IPv4/IPv6 packets to
// which the RSS hashing must be applied.  Supplying an *rss_hf* equal
// to zero disables the RSS feature.
type RssConf struct {
	/**< If not NULL, 40-byte hash key. */
	Key []byte
	/**< Hash functions to apply. */
	Hf uint64
}

// MACAddr is a universally administered address is uniquely assigned
// to a device by its manufacturer. The first three octets (in
// transmission order) contain the Organizationally Unique Identifier
// (OUI). The following three (MAC-48 and EUI-48) octets are assigned
// by that organization with the only constraint of uniqueness. A
// locally administered address is assigned to a device by a network
// administrator and does not contain OUIs.
type MACAddr C.struct_rte_ether_addr

// MACAddrGet retrieves the Ethernet address of an Ethernet device.
func (pid Port) MACAddrGet(addr *MACAddr) error {
	return errget(C.rte_eth_macaddr_get(C.ushort(pid), (*C.struct_rte_ether_addr)(addr)))
}

// HardwareAddr converts MACAddr into Go's native net.HardwareAddr.
func (addr *MACAddr) HardwareAddr() net.HardwareAddr {
	p := &addr.addr_bytes
	return (*[unsafe.Sizeof(*p)]byte)(unsafe.Pointer(&p[0]))[:]
}

func (addr *MACAddr) String() string {
	return addr.HardwareAddr().String()
}

// RssHashUpdate updates configuration of Receive Side Scaling hash
// computation of Ethernet device.
func (pid Port) RssHashUpdate(conf *RssConf) error {
	var rssConf C.struct_rte_eth_rss_conf

	rssConf.rss_key_len = C.uchar(len(conf.Key))
	rssConf.rss_hf = C.ulong(conf.Hf)

	if len(conf.Key) > 0 {
		p := C.CBytes(conf.Key)
		defer C.free(p)
		rssConf.rss_key = (*C.uchar)(p)
	}

	return errget(C.rte_eth_dev_rss_hash_update(C.ushort(pid), &rssConf))
}

// RssHashConfGet retrieves current configuration of Receive Side
// Scaling hash computation of Ethernet device.
func (pid Port) RssHashConfGet(conf *RssConf) error {
	var rssConf C.struct_rte_eth_rss_conf
	rc := C.rte_eth_dev_rss_hash_conf_get(C.ushort(pid), &rssConf)

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&conf.Key))
	sh.Data = uintptr(unsafe.Pointer(rssConf.rss_key))
	sh.Len = int(rssConf.rss_key_len)
	sh.Cap = sh.Len

	conf.Hf = uint64(rssConf.rss_hf)
	return errget(rc)
}

// RssRetaEntry64 is a structure used to configure 64 entries of
// Redirection Table of the Receive Side Scaling (RSS) feature of an
// Ethernet port. To configure more than 64 entries supported by
// hardware, an array of this structure is needed.
type RssRetaEntry64 C.struct_rte_eth_rss_reta_entry64

// Reta returns group of 64 redirection table entries. You may set
// elements of returned slice but no appends are allowed.
func (conf *RssRetaEntry64) Reta() (reta []uint16) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&reta))
	sh.Len = len(conf.reta)
	sh.Cap = len(conf.reta)
	sh.Data = uintptr(unsafe.Pointer(&conf.reta))
	return
}

// Mask returns bits indicating which entries need to be
// updated/queried.
func (conf *RssRetaEntry64) Mask() *uint64 {
	return (*uint64)(&conf.mask)
}

// RssRetaUpdate updates Redirection Table(RETA) of Receive Side
// Scaling of Ethernet device.
//
// conf is a RETA to update. retaSize is redirection table size. The
// table size can be queried by rte_eth_dev_info_get().
//
// Returns:
//
//	(0) if successful.
//	(-ENODEV) if port_id is invalid.
//	(-ENOTSUP) if hardware doesn't support.
//	(-EINVAL) if bad parameter.
//	(-EIO) if device is removed.
func (pid Port) RssRetaUpdate(conf []RssRetaEntry64, retaSize uint16) error {
	p := (*C.struct_rte_eth_rss_reta_entry64)(nil)
	if len(conf) > 0 {
		p = (*C.struct_rte_eth_rss_reta_entry64)(&conf[0])
	}
	return errget(C.rte_eth_dev_rss_reta_update(C.ushort(pid), p, C.ushort(retaSize)))
}

// RssRetaQuery queries Redirection Table(RETA) of Receive Side Scaling of Ethernet device.
//
// conf is a RETA to query. For each requested reta entry, corresponding bit in mask must be set.
// retaSize is a redirection table size. The table size can be queried by rte_eth_dev_info_get().
//
// Returns:
//
//	(0) if successful.
//	(-ENODEV) if port_id is invalid.
//	(-ENOTSUP) if hardware doesn't support.
//	(-EINVAL) if bad parameter.
//	(-EIO) if device is removed.
func (pid Port) RssRetaQuery(conf []RssRetaEntry64, retaSize uint16) error {
	p := (*C.struct_rte_eth_rss_reta_entry64)(nil)
	if len(conf) > 0 {
		p = (*C.struct_rte_eth_rss_reta_entry64)(&conf[0])
	}
	return errget(C.rte_eth_dev_rss_reta_query(C.ushort(pid), p, C.ushort(retaSize)))
}

// Thresh is a structure used to configure the ring threshold
// registers of an RX/TX queue for an Ethernet port.
type Thresh struct {
	// Ring prefetch threshold.
	PThresh uint8
	// Ring host threshold.
	HThresh uint8
	// Ring writeback threshold.
	WThresh uint8
}

// RxqConf is a s tructure used to configure an RX ring of an Ethernet
// port.
type RxqConf struct {
	Thresh
	// Drives the freeing of RX descriptors.
	FreeThresh uint16
	// Drop packets if no descriptors are available.
	DropEn uint8
	// Do not start queue with rte_eth_dev_start().
	DeferredStart uint8
	// Per-queue Rx offloads to be set using RxOffload* flags.
	// Only offloads set on rx_queue_offload_capa or rx_offload_capa
	// fields on rte_eth_dev_info structure are allowed to be set.
	Offloads uint64
}

// TxqConf is a structure used to configure a TX ring of an Ethernet
// port.
type TxqConf struct {
	Thresh
	// Drives the setting of RS bit on TXDs.
	RsThresh uint16
	// Start freeing TX buffers if there are less free descriptors
	// than this value.
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

func errget(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

func (t *Thresh) cThresh() (out C.struct_rte_eth_thresh) {
	out.pthresh = C.uchar(t.PThresh)
	out.hthresh = C.uchar(t.HThresh)
	out.wthresh = C.uchar(t.WThresh)
	return
}

// OptLinkSpeeds sets allowed speeds for the device.
// LinkSpeedFixed disables link autonegotiation, and a unique speed
// shall be set. Otherwise, the bitmap defines the set of speeds to be
// advertised. If the special value LinkSpeedAutoneg is used, all
// speeds supported are advertised.
func OptLinkSpeeds(speeds uint) Option {
	return Option{func(c *ethConf) {
		c.conf.link_speeds = C.uint(speeds)
	}}
}

// OptRxMode specifies port RX configuration.
func OptRxMode(conf RxMode) Option {
	return Option{func(c *ethConf) {
		c.conf.rxmode = C.struct_rte_eth_rxmode{
			mq_mode:        uint32(conf.MqMode),
			split_hdr_size: C.ushort(conf.SplitHdrSize),
			offloads:       C.ulong(conf.Offloads),
		}
		c.setRxPktLen(conf.MTU)
	}}
}

func (c *ethConf) setRxPktLen(n uint32) {
	rxptr := unsafe.Pointer((*C.struct_rte_eth_rxmode)(&c.conf.rxmode))
	p := unsafe.Pointer(uintptr(rxptr) + C.RX_MODE_LEN_OFF)
	*(*C.uint32_t)(p) = (C.uint32_t)(n)
}

type ethConf struct {
	conf  C.struct_rte_eth_conf
	cptrs []unsafe.Pointer // allocated C pointers
}

func (c *ethConf) free() {
	for _, p := range c.cptrs {
		C.free(p)
	}
}

// OptTxMode specifies port TX configuration.
func OptTxMode(conf TxMode) Option {
	return Option{func(c *ethConf) {
		c.conf.txmode = C.struct_rte_eth_txmode{
			mq_mode:  uint32(conf.MqMode),
			offloads: C.ulong(conf.Offloads),
			pvid:     C.ushort(conf.Pvid),
		}
		if conf.HwVlanRejectTagged {
			C.set_tx_reject_tagged(&c.conf.txmode)
		}
		if conf.HwVlanRejectUntagged {
			C.set_tx_reject_untagged(&c.conf.txmode)
		}
		if conf.HwVlanInsertPvid {
			C.set_tx_insert_pvid(&c.conf.txmode)
		}
	}}
}

// OptLoopbackMode specifies loopback operation mode. By default
// the value is 0, meaning the loopback mode is disabled.  Read the
// datasheet of given ethernet controller for details. The possible
// values of this field are defined in implementation of each driver.
func OptLoopbackMode(mode uint32) Option {
	return Option{func(c *ethConf) {
		c.conf.lpbk_mode = C.uint(mode)
	}}
}

// OptRss specifies RSS configuration.
func OptRss(conf RssConf) Option {
	return Option{func(c *ethConf) {
		rssConf := &c.conf.rx_adv_conf.rss_conf
		rssConf.rss_key_len = C.uchar(len(conf.Key))
		rssConf.rss_hf = C.ulong(conf.Hf)
		if conf.Key != nil && len(conf.Key) > 0 {
			p := C.CBytes(conf.Key)
			rssConf.rss_key = (*C.uchar)(p)
			c.cptrs = append(c.cptrs, p)
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
// Several Opt* options may be specified as well.
func (pid Port) DevConfigure(nrxq, ntxq uint16, opts ...Option) error {
	ec := &ethConf{}
	for i := range opts {
		opts[i].f(ec)
	}
	defer ec.free()

	return errget(C.rte_eth_dev_configure(C.ushort(pid), C.ushort(nrxq),
		C.ushort(ntxq), &ec.conf))
}

// OptRxqConf specifies the configuration an RX ring of an Ethernet
// port.
func OptRxqConf(conf RxqConf) QueueOption {
	return QueueOption{func(q *qConf) {
		q.rx = C.struct_rte_eth_rxconf{
			rx_thresh:         conf.Thresh.cThresh(),
			rx_free_thresh:    C.ushort(conf.FreeThresh),
			rx_drop_en:        C.uchar(conf.DropEn),
			rx_deferred_start: C.uchar(conf.DeferredStart),
			offloads:          C.ulong(conf.Offloads),
		}
	}}
}

// OptTxqConf specifies the configuration an TX ring of an Ethernet
// port.
func OptTxqConf(conf TxqConf) QueueOption {
	return QueueOption{func(q *qConf) {
		q.tx = C.struct_rte_eth_txconf{
			tx_thresh:         conf.Thresh.cThresh(),
			tx_rs_thresh:      C.ushort(conf.RsThresh),
			tx_free_thresh:    C.ushort(conf.FreeThresh),
			tx_deferred_start: C.uchar(conf.DeferredStart),
			offloads:          C.ulong(conf.Offloads),
		}
	}}
}

// OptSocket specifies the NUMA socket id for RX/TX queue.  The socket
// argument is the socket identifier in case of NUMA.  The value can
// be SOCKET_ID_ANY if there is no NUMA constraint for the DMA memory
// allocated for the receive/transmit descriptors of the ring.
func OptSocket(socket int) QueueOption {
	return QueueOption{func(q *qConf) {
		q.socket = C.int(socket)
	}}
}

// RxqSetup allocates and sets up a receive queue for an Ethernet
// device.
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
//
// opts specifies the configuration data to be used for the receive
// queue.  The *rx_conf* structure contains an *rx_thresh* structure
// with the values of the Prefetch, Host, and Write-Back threshold
// registers of the receive ring.  In addition it contains the
// hardware offloads features to activate using the DEV_RX_OFFLOAD_*
// flags.  If an offloading set in rx_conf->offloads hasn't been set
// in the input argument eth_conf->rxmode.offloads to
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
//   - -EINVAL: The memory pool pointer is null or the size of network
//     buffers which can be allocated from this memory pool does not
//     fit the various buffer sizes allowed by the device controller.
//
//   - -ENOMEM: Unable to allocate the receive ring descriptors or to
//     allocate network memory buffers from the memory pool when
//     initializing receive descriptors.
func (pid Port) RxqSetup(qid, nDesc uint16, mp *mempool.Mempool, opts ...QueueOption) error {
	conf := &qConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	return errget(C.rte_eth_rx_queue_setup(C.ushort(pid), C.ushort(qid),
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
// OptSocket specifies the socket identifier in case of NUMA.
// Its value can be *SOCKET_ID_ANY* if there is no NUMA constraint for
// the DMA memory allocated for the transmit descriptors of the ring.
//
// OptTxqConf specifies configuration data to be used for the transmit
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
func (pid Port) TxqSetup(qid, nDesc uint16, opts ...QueueOption) error {
	conf := &qConf{socket: C.SOCKET_ID_ANY}
	for i := range opts {
		opts[i].f(conf)
	}

	return errget(C.rte_eth_tx_queue_setup(C.ushort(pid), C.ushort(qid),
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
	return errget(C.rte_eth_dev_reset(C.ushort(pid)))
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
	return errget(C.rte_eth_dev_start(C.ushort(pid)))
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
func (pid Port) PromiscEnable() error {
	return errget(C.rte_eth_promiscuous_enable(C.ushort(pid)))
}

// PromiscDisable disables receipt in promiscuous mode for an Ethernet
// device.
func (pid Port) PromiscDisable() error {
	return errget(C.rte_eth_promiscuous_disable(C.ushort(pid)))
}

// SetLinkUp set link status to 'up' an Ethernet device.
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
	return errget(C.rte_eth_dev_set_link_up(C.ushort(pid)))
}

// SetLinkDown set link status to 'down' an Ethernet device.
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
	return errget(C.rte_eth_dev_set_link_down(C.ushort(pid)))
}

// SocketID returns the NUMA socket to which an Ethernet device is
// connected.  The function may return a default of zero if the socket
// could not be determined. -1 is returned if the port_id value is out
// of range.
func (pid Port) SocketID() int {
	return int(C.rte_eth_dev_socket_id(C.ushort(pid)))
}

// DevInfo is a structure used to retrieve the contextual information
// of an Ethernet device, such as the controlling driver of the
// device, etc...
type DevInfo C.struct_rte_eth_dev_info

// DriverName returns driver_name as a Go string.
func (info *DevInfo) DriverName() string {
	return C.GoString((*C.struct_rte_eth_dev_info)(info).driver_name)
}

// InterfaceName is the name of the interface in the system.
func (info *DevInfo) InterfaceName() string {
	var buf [C.IF_NAMESIZE]C.char
	return C.GoString(C.if_indextoname(info.if_index, &buf[0]))
}

// RetaSize returns Device redirection table size, the total number of
// entries.
func (info *DevInfo) RetaSize() uint16 {
	return uint16(info.reta_size)
}

// InfoGet retrieves the contextual information of an Ethernet device.
//
// As part of this function, a number of of fields in dev_info will be
// initialized as follows:
//
// rx_desc_lim = lim tx_desc_lim = lim
//
// Where lim is defined within the rte_eth_dev_info_get as
//
//	const struct rte_eth_desc_lim lim = { .nb_max = UINT16_MAX, .nb_min
//	= 0, .nb_align = 1, .nb_seg_max = UINT16_MAX, .nb_mtu_seg_max =
//	UINT16_MAX, };
//
//	device = dev->device min_mtu = RTE_ETHER_MIN_MTU max_mtu =
//	UINT16_MAX
//
// The following fields will be populated if support for
// dev_infos_get() exists for the device and the rte_eth_dev 'dev' has
// been populated successfully with a call to it:
//
//	driver_name = dev->device->driver->name nb_rx_queues =
//	dev->data->nb_rx_queues nb_tx_queues = dev->data->nb_tx_queues
//	dev_flags = &dev->data->dev_flags
func (pid Port) InfoGet(info *DevInfo) error {
	return errget(C.rte_eth_dev_info_get(C.ushort(pid), (*C.struct_rte_eth_dev_info)(info)))
}

// Name get the device name from port id. The device name is specified as below:
//
//   - PCIe address (Domain:Bus:Device.Function), for example- 0000:02:00.0
//   - SoC device name, for example- fsl-gmac0
//   - vdev dpdk name, for example- net_[pcap0|null0|tun0|tap0]
//
// (0) if successful.
// (-ENODEV) if port_id is invalid.
// (-EINVAL) on failure.
func (pid Port) Name() (string, error) {
	var buf [C.RTE_ETH_NAME_MAX_LEN]C.char
	err := errget(C.rte_eth_dev_get_name_by_port(C.ushort(pid), &buf[0]))
	return C.GoString(&buf[0]), err
}

// GetPortByName gets the port ID from device name. The device name
// should be specified as below:
//
//   - PCIe address (Domain:Bus:Device.Function), for example- 0000:02:00.0
//   - SoC device name, for example- fsl-gmac0
//   - vdev dpdk name, for example- net_[pcap0|null0|tun0|tap0]
//
// (0) if successful.
// (-ENODEV or -EINVAL) in case of failure.
func GetPortByName(name string) (Port, error) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	var p C.ushort
	return Port(p), errget(C.rte_eth_dev_get_port_by_name(cstr, &p))
}

// NbRxQueues returns number of configured RX queues for the device.
func (info *DevInfo) NbRxQueues() uint16 {
	return uint16(info.nb_rx_queues)
}

// NbTxQueues returns number of configured TX queues for the device.
func (info *DevInfo) NbTxQueues() uint16 {
	return uint16(info.nb_tx_queues)
}

// IsValid checks if port_id of device is attached.
func (pid Port) IsValid() bool {
	return C.rte_eth_dev_is_valid_port(C.ushort(pid)) != 0
}

// CountAvail gets the number of ports which are usable for the
// application.
//
// These devices must be iterated by using the macro
// RTE_ETH_FOREACH_DEV or RTE_ETH_FOREACH_DEV_OWNED_BY to deal with
// non-contiguous ranges of devices.
func CountAvail() int {
	return int(C.rte_eth_dev_count_avail())
}

// ValidPorts returns list of devices valid in EAL.
func ValidPorts() []Port {
	avail := CountAvail()
	if avail == 0 {
		return nil
	}

	ports := make([]Port, avail)
	n := C.go_rte_get_ethdev_ports((*C.uint16_t)(unsafe.Pointer(&ports[0])), C.int(avail))
	return ports[:n]
}

// CountTotal gets the total number of ports which are allocated.
//
// Some devices may not be available for the application.
func CountTotal() int {
	return int(C.rte_eth_dev_count_total())
}

// EthLink is a structure used to retrieve link-level information of
// an Ethernet port.
type EthLink C.struct_go_rte_eth_link

// Speed returns speed in Mbps.
func (link *EthLink) Speed() uint32 {
	return uint32(link.link_speed)
}

// Duplex returns true if the port is in full-duplex, otherwise it's
// in half-duplex.
func (link *EthLink) Duplex() bool {
	return link.link_duplex > 0
}

// AutoNeg returns true if auto-negotiation is true.
func (link *EthLink) AutoNeg() bool {
	return link.link_autoneg > 0
}

// Status returns true if link is on, false otherwise.
func (link *EthLink) Status() bool {
	return link.link_status > 0
}

// EthLinkGet retrieves the link status (up/down), the duplex mode
// (half/full), the negotiation (auto/fixed), and if available, the
// speed (Mbps).
//
// It might need to wait up to 9 seconds.
//
// Returns:
//
//	(0) if successful.
//	(-ENOTSUP) if the function is not supported in PMD driver.
//	(-ENODEV) if port_id invalid.
//	(-EINVAL) if bad parameter.
func (pid Port) EthLinkGet() (EthLink, error) {
	var d EthLink
	return d, errget(C.go_rte_eth_link_get(C.ushort(pid), (*C.struct_go_rte_eth_link)(&d)))
}

// EthLinkGetNowait retrieves the link status (up/down), the duplex
// mode (half/full), the negotiation (auto/fixed), and if available,
// the speed (Mbps).
//
// Returns:
//
//	(0) if successful.
//	(-ENOTSUP) if the function is not supported in PMD driver.
//	(-ENODEV) if port_id invalid.
//	(-EINVAL) if bad parameter.
func (pid Port) EthLinkGetNowait() (EthLink, error) {
	var d EthLink
	return d, errget(C.go_rte_eth_link_get_nowait(C.ushort(pid), (*C.struct_go_rte_eth_link)(&d)))
}

// FcConf is a structure used to configure Ethernet flow control
// parameter.  These parameters will be configured into the register
// of the NIC. Please refer to the corresponding data sheet for proper
// value.
type FcConf C.struct_rte_eth_fc_conf

// SetMode sets Flow Control mode.
func (conf *FcConf) SetMode(n uint32) {
	conf.mode = n
}

// FlowCtrlGet gets current status of the Ethernet link flow control for Ethernet device.
func (pid Port) FlowCtrlGet(conf *FcConf) error {
	return errget(C.rte_eth_dev_flow_ctrl_get(C.ushort(pid), (*C.struct_rte_eth_fc_conf)(conf)))
}

// FlowCtrlSet configures the Ethernet link flow control for Ethernet device.
func (pid Port) FlowCtrlSet(conf *FcConf) error {
	return errget(C.rte_eth_dev_flow_ctrl_set(C.ushort(pid), (*C.struct_rte_eth_fc_conf)(conf)))
}

func fromRings(rings []*ring.Ring) (p **C.struct_rte_ring, n C.uint) {
	if n = C.uint(len(rings)); n > 0 {
		p = (**C.struct_rte_ring)(unsafe.Pointer(&rings[0]))
	}
	return
}

// FromRings creates new ethdev port from a set of rings. rxqs and
// txqs represent RX and TX queues, socket is the NUMA node, name is
// to be given to the new port.
func FromRings(name string, rxqs, txqs []*ring.Ring, socket int) (Port, error) {
	rxq, nrx := fromRings(rxqs)
	txq, ntx := fromRings(txqs)
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	rc := C.rte_eth_from_rings(cname, rxq, nrx, txq, ntx, C.uint(socket))
	n, err := common.IntOrErr(rc)
	return Port(n), err
}

// FromRing creates new ethdev port from a ring.
func FromRing(r *ring.Ring) (Port, error) {
	rc := C.rte_eth_from_ring((*C.struct_rte_ring)(unsafe.Pointer(r)))
	n, err := common.IntOrErr(rc)
	return Port(n), err
}
