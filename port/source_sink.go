package port

/*
#include <rte_config.h>
#include <rte_port_source_sink.h>
*/
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/yerden/go-dpdk/mempool"
)

// compile time checks
var _ = []ReaderParams{
	&Source{},
}

var _ = []WriterParams{
	&Sink{},
}

// Source is an input port that can be used to generate packets.
type Source struct {
	// Pre-initialized buffer pool.
	*mempool.Mempool

	// The full path of the pcap file to read packets from.
	Filename string

	// The number of bytes to be read from each packet in the pcap file. If
	// this value is 0, the whole packet is read; if it is bigger than packet
	// size, the generated packets will contain the whole packet.
	BytesPerPacket uint32
}

// ReaderOps implements ReaderParams interface.
func (rd *Source) ReaderOps() (*ReaderOps, unsafe.Pointer) {
	ops := (*ReaderOps)(&C.rte_port_source_ops)
	rc := &C.struct_rte_port_source_params{}
	rc.mempool = (*C.struct_rte_mempool)(unsafe.Pointer(rd.Mempool))
	rc.n_bytes_per_pkt = C.uint32_t(rd.BytesPerPacket)
	if rd.Filename != "" {
		rc.file_name = C.CString(rd.Filename)
		// file_name no longer needed once rc is out of use, it
		// doesn't persist in port itself, so we may free it as long
		// as it's out of reach.
		runtime.SetFinalizer(rc, func(rc *C.struct_rte_port_source_params) {
			C.free(unsafe.Pointer(rc.file_name))
		})
	}
	return ops, unsafe.Pointer(rc)
}

// Sink is an output port that drops all packets written to it.
type Sink struct {
	// The full path of the pcap file to write the packets to.
	Filename string

	// The maximum number of packets write to the pcap file. If this value is
	// 0, the "infinite" write will be carried out.
	MaxPackets uint32
}

// WriterOps implements WriterParams interface.
func (wr *Sink) WriterOps() (*WriterOps, unsafe.Pointer) {
	ops := (*WriterOps)(&C.rte_port_sink_ops)
	rc := &C.struct_rte_port_sink_params{}
	rc.max_n_pkts = C.uint32_t(wr.MaxPackets)
	if wr.Filename != "" {
		rc.file_name = C.CString(wr.Filename)
		// file_name no longer needed once rc is out of use, it
		// doesn't persist in port itself, so we may free it as long
		// as it's out of reach.
		runtime.SetFinalizer(rc, func(rc *C.struct_rte_port_sink_params) {
			C.free(unsafe.Pointer(rc.file_name))
		})
	}
	return ops, unsafe.Pointer(rc)
}
