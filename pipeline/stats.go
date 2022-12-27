package pipeline

/*
#include <rte_pipeline.h>
*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/port"
	"github.com/yerden/go-dpdk/table"
)

// PortInStats is the pipeline input port stats.
type PortInStats struct {
	port.InStats

	// Number of packets dropped by action handler.
	PacketsDroppedByAH uint64
}

// PortOutStats is the pipeline output port stats.
type PortOutStats struct {
	port.OutStats

	// Number of packets dropped by action handler.
	PacketsDroppedByAH uint64
}

// TableStats is the pipeline table stats.
type TableStats struct {
	table.Stats

	PacketsDroppedByHitAH  uint64
	PacketsDroppedByMissAH uint64
	PacketsDroppedByHit    uint64
	PacketsDroppedByMiss   uint64
}

var _ = []uintptr{
	unsafe.Sizeof(PortInStats{}) - unsafe.Sizeof(C.struct_rte_pipeline_port_in_stats{}),
	unsafe.Sizeof(C.struct_rte_pipeline_port_in_stats{}) - unsafe.Sizeof(PortInStats{}),
	unsafe.Sizeof(PortOutStats{}) - unsafe.Sizeof(C.struct_rte_pipeline_port_out_stats{}),
	unsafe.Sizeof(C.struct_rte_pipeline_port_out_stats{}) - unsafe.Sizeof(PortOutStats{}),
	unsafe.Sizeof(TableStats{}) - unsafe.Sizeof(C.struct_rte_pipeline_table_stats{}),
	unsafe.Sizeof(C.struct_rte_pipeline_table_stats{}) - unsafe.Sizeof(TableStats{}),
}

// PortOutStatsRead reads stats on output port registered in the
// pipeline.
//
// If clear s true then clear stats after reading.
func (pl *Pipeline) PortOutStatsRead(port PortOut, s *PortOutStats, clear bool) error {
	var c C.int
	if clear {
		c = 1
	}

	return common.IntErr(int64(C.rte_pipeline_port_out_stats_read(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(port),
		(*C.struct_rte_pipeline_port_out_stats)(unsafe.Pointer(s)),
		c)))
}

// PortInStatsRead reads stats on input port registered in the
// pipeline.
//
// If clear s true then clear stats after reading.
func (pl *Pipeline) PortInStatsRead(port PortIn, s *PortInStats, clear bool) error {
	var c C.int
	if clear {
		c = 1
	}

	return common.IntErr(int64(C.rte_pipeline_port_in_stats_read(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(port),
		(*C.struct_rte_pipeline_port_in_stats)(unsafe.Pointer(s)),
		c)))
}

// TableStatsRead reads stats on table registered in the pipeline.
//
// If clear s true then clear stats after reading.
func (pl *Pipeline) TableStatsRead(table Table, s *TableStats, clear bool) error {
	var c C.int
	if clear {
		c = 1
	}

	return common.IntErr(int64(C.rte_pipeline_table_stats_read(
		(*C.struct_rte_pipeline)(pl),
		C.uint32_t(table),
		(*C.struct_rte_pipeline_table_stats)(unsafe.Pointer(s)),
		c)))
}
