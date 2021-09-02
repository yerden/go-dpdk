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

	"github.com/yerden/go-dpdk/common"
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

type statsAsU64 [unsafe.Sizeof(Stats{}) / 8]uint64

// Diff computes delta between s and old.
func (s *Stats) Diff(old, delta *Stats) {
	cur := (*statsAsU64)(unsafe.Pointer(s))
	src := (*statsAsU64)(unsafe.Pointer(old))
	d := (*statsAsU64)(unsafe.Pointer(delta))

	for i := range cur {
		if cur[i] > src[i] && src[i] != 0 {
			d[i] = cur[i] - src[i]
		}
	}
}

// StatsGet retrieves statistics from ethernet device.
func (pid Port) StatsGet(stats *Stats) error {
	return errget(C.rte_eth_stats_get(C.ushort(pid), (*C.struct_rte_eth_stats)(stats)))
}

// XstatNames returns list of names for custom eth dev stats.
func (pid Port) XstatNames() ([]XstatName, error) {
	n := C.rte_eth_xstats_get_names(C.ushort(pid), nil, 0)
	if n <= 0 {
		return nil, errget(n)
	}

	names := make([]XstatName, n)
	C.rte_eth_xstats_get_names(C.ushort(pid), (*C.struct_rte_eth_xstat_name)(&names[0]), C.uint(n))
	return names, nil
}

// XstatsGet retrieves xstat from eth dev. Returns number of retrieved
// statistics and possible error. If returned number is greater than
// len(out) extend the slice and try again.
func (pid Port) XstatsGet(out []Xstat) (int, error) {
	var p *C.struct_rte_eth_xstat
	if len(out) != 0 {
		p = (*C.struct_rte_eth_xstat)(unsafe.Pointer(&out[0]))
	}

	n := C.rte_eth_xstats_get(C.ushort(pid), p, C.uint(len(out)))
	return common.IntOrErr(n)
}

// XstatsReset resets xstats counters.
func (pid Port) XstatsReset() error {
	return errget(C.rte_eth_xstats_reset(C.ushort(pid)))
}

// StatsReset resets stats counters.
func (pid Port) StatsReset() error {
	return errget(C.rte_eth_stats_reset(C.ushort(pid)))
}

// return position of an Xstat with Index equal to idx, or -1 if not
// found.
func searchXstat(in []Xstat, idx uint64) int {
	for n := range in {
		if in[n].Index == idx {
			return n
		}
	}
	return -1
}

// XstatDiff computes delta from incoming stats array (incoming) and
// old stats array (old) with the result placed in delta.
//
// If old contains zero-value Xstat the resulting delta is zero to
// filter out outliers.
//
// If old doesn't contain Xstat contained in incoming new Xstat is
// appended to old.
//
// If delta doesn't contain new Xstat it is appended.
//
// The resulting "old" (now updated) and delta are returned.
func XstatDiff(incoming, old, delta []Xstat) (newOld, newDelta []Xstat) {
	for _, d := range incoming {
		if n := searchXstat(old, d.Index); n < 0 {
			old = append(old, d)
			d.Value = 0
		} else if f := old[n].Value; f > 0 && f < d.Value {
			old[n].Value = d.Value
			d.Value -= f
		} else {
			old[n].Value = d.Value
			d.Value = 0
		}

		if k := searchXstat(delta, d.Index); k < 0 {
			delta = append(delta, d)
		} else {
			delta[k].Value = d.Value
		}
	}

	return old, delta
}
