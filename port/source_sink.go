package port

/*
#include <rte_config.h>
#include <rte_port_source_sink.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/mempool"
)

// compile time checks
var _ = []RxFactory{
	&Source{},
}

var _ = []TxFactory{
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

// CreateRx implements RxFactory interface.
func (rd *Source) CreateRx(socket int) (*Rx, error) {
	rx := &Rx{ops: &C.rte_port_source_ops}

	// port
	params := &C.struct_rte_port_source_params{
		mempool:         (*C.struct_rte_mempool)(unsafe.Pointer(rd.Mempool)),
		n_bytes_per_pkt: C.uint32_t(rd.BytesPerPacket),
	}

	if rd.Filename != "" {
		params.file_name = (*C.char)(C.CString(rd.Filename))
		defer C.free(unsafe.Pointer(params.file_name))
	}

	return rx, rx.doCreate(socket, unsafe.Pointer(params))
}

// Sink is an output port that drops all packets written to it.
type Sink struct {
	// The full path of the pcap file to write the packets to.
	Filename string

	// The maximum number of packets write to the pcap file. If this value is
	// 0, the "infinite" write will be carried out.
	MaxPackets uint32
}

// CreateTx implements TxFactory interface.
func (wr *Sink) CreateTx(socket int) (*Tx, error) {
	tx := &Tx{ops: &C.rte_port_sink_ops}

	// port
	params := &C.struct_rte_port_sink_params{
		max_n_pkts: C.uint32_t(wr.MaxPackets),
	}

	if wr.Filename != "" {
		params.file_name = (*C.char)(C.CString(wr.Filename))
		defer C.free(unsafe.Pointer(params.file_name))
	}

	return tx, tx.doCreate(socket, unsafe.Pointer(params))
}
