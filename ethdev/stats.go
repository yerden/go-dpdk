package ethdev

/*
#include <rte_config.h>
#include <rte_errno.h>
#include <rte_memory.h>
#include <rte_ethdev.h>
*/
import "C"
import (
	"unsafe"
)

// for testing
type cXstat C.struct_rte_eth_xstat

// Stats represents basic statistics counters C struct.
type Stats C.struct_rte_eth_stats

// Xstat reflect C rte_eth_xstats struct.
type Xstat struct {
	// Index of an element in XstatsNames array.
	Index uint64

	// Value of stats.
	Value uint64
}

// XstatName reflect C rte_eth_xstats_name struct.
type XstatName C.struct_rte_eth_xstat_name

// String returns name of the stat.
func (x *XstatName) String() string {
	return C.GoString(&x.name[0])
}

// GoStats contains counters from rte_eth_stats.
type GoStats struct {
	Ipackets uint64 `metric:"ipackets" type:"counter"`
	Opackets uint64 `metric:"opackets" type:"counter"`
	Ibytes   uint64 `metric:"ibytes" type:"counter"`
	Obytes   uint64 `metric:"obytes" type:"counter"`
	Imissed  uint64 `metric:"imissed" type:"counter"`
	Ierrors  uint64 `metric:"ierrors" type:"counter"`
	Oerrors  uint64 `metric:"oerrors" type:"counter"`
	RxNoMbuf uint64 `metric:"rxnombuf" type:"counter"`
}

// Cast transforms Stats pointer to GoStats.
func (s *Stats) Cast() *GoStats {
	return (*GoStats)(unsafe.Pointer(s))
}

// StatsGet retrieves statistics from ethernet device.
func (pid Port) StatsGet(stats *Stats) error {
	return err(C.rte_eth_stats_get(C.ushort(pid), (*C.struct_rte_eth_stats)(stats)))
}

// XstatNames returns list of names for custom eth dev stats.
func (pid Port) XstatNames() ([]XstatName, error) {
	n := C.rte_eth_xstats_get_names(C.ushort(pid), nil, 0)
	if n <= 0 {
		return nil, err(n)
	}

	names := make([]XstatName, n)
	C.rte_eth_xstats_get_names(C.ushort(pid), (*C.struct_rte_eth_xstat_name)(&names[0]), C.uint(n))
	return names, nil
}

// XstatsGet retrieves xstat from eth dev. Returns number of retrieved
// statistics and possible error. If len(*pout) is insufficient to
// store all stats it gets extended.
func (pid Port) XstatsGet(pout *[]Xstat) (int, error) {
	for {
		out := *pout

		var p *C.struct_rte_eth_xstat
		if len(out) != 0 {
			p = (*C.struct_rte_eth_xstat)(unsafe.Pointer(&out[0]))
		}

		n := int(C.rte_eth_xstats_get(C.ushort(pid), p, C.uint(len(out))))
		if n < 0 {
			return 0, err(n)
		}

		if n >= len(out) {
			return n, nil
		}

		out = append(out[:0], make([]Xstat, n)...)
		*pout = out
	}
}

// XstatsReset resets xstats counters.
func (pid Port) XstatsReset() error {
	return err(C.rte_eth_xstats_reset(C.ushort(pid)))
}

// StatsReset resets stats counters.
func (pid Port) StatsReset() error {
	return err(C.rte_eth_stats_reset(C.ushort(pid)))
}
