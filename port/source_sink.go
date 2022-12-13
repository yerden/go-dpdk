package port

/*
#include <rte_config.h>
#include <rte_port_source_sink.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/mempool"
)

// compile time checks
var _ = []InParams{
	&Source{},
}

var _ = []OutParams{
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

// InOps implements InParams interface.
func (rd *Source) InOps() *InOps {
	return (*InOps)(&C.rte_port_source_ops)
}

// Transform implements common.Transformer interface.
func (rd *Source) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	// port
	var params *C.struct_rte_port_source_params
	params = (*C.struct_rte_port_source_params)(alloc.Malloc(unsafe.Sizeof(*params)))
	params.mempool = (*C.struct_rte_mempool)(unsafe.Pointer(rd.Mempool))
	params.n_bytes_per_pkt = C.uint32_t(rd.BytesPerPacket)

	if rd.Filename != "" {
		params.file_name = (*C.char)(common.CString(alloc, rd.Filename))
	}

	return unsafe.Pointer(params), func(arg unsafe.Pointer) {
		params := (*C.struct_rte_port_source_params)(arg)
		alloc.Free(unsafe.Pointer(params.file_name))
		alloc.Free(arg)
	}
}

// Sink is an output port that drops all packets written to it.
type Sink struct {
	// The full path of the pcap file to write the packets to.
	Filename string

	// The maximum number of packets write to the pcap file. If this value is
	// 0, the "infinite" write will be carried out.
	MaxPackets uint32
}

// OutOps implements OutParams interface.
func (wr *Sink) OutOps() *OutOps {
	return (*OutOps)(&C.rte_port_sink_ops)
}

// Transform implements common.Transformer interface.
func (wr *Sink) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	var params *C.struct_rte_port_sink_params
	params = (*C.struct_rte_port_sink_params)(alloc.Malloc(unsafe.Sizeof(*params)))
	params.max_n_pkts = C.uint32_t(wr.MaxPackets)

	if wr.Filename != "" {
		params.file_name = (*C.char)(common.CString(alloc, wr.Filename))
	}

	return unsafe.Pointer(params), func(arg unsafe.Pointer) {
		params := (*C.struct_rte_port_sink_params)(arg)
		alloc.Free(unsafe.Pointer(params.file_name))
		alloc.Free(unsafe.Pointer(arg))
	}
}
