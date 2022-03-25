package ethdev

/*
#include <rte_config.h>
#include <rte_errno.h>
#include <rte_memory.h>
#include <rte_ethdev.h>
*/
import "C"
import (
	"os"
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

func (pid Port) xstatIDByName(x *XstatName) (id uint64, err error) {
	return id, errget(C.rte_eth_xstats_get_id_by_name(C.ushort(pid), &x.name[0], (*C.uint64_t)(&id)))
}

// XstatNameIDs returns names of extended statistics mapped by their
// ID and an error. This is supposed to be called to cache counter
// ids and not used in hot path.
func (pid Port) XstatNameIDs() (map[uint64]string, error) {
	names, err := pid.XstatNames()
	if err != nil {
		return nil, err
	}

	var id uint64
	namesID := map[uint64]string{}
	for i := range names {
		x := &names[i]

		if id, err = pid.xstatIDByName(x); err != nil {
			return nil, err
		}

		if _, ok := namesID[id]; ok {
			return nil, os.ErrExist
		}

		namesID[id] = x.String()
	}

	return namesID, nil
}

// XstatGetByID retrieves extended statistics of an Ethernet device.
//
// ids is the IDs array given by app to retrieve specific statistics.
// May be nil to retrieve all available statistics or, if values is
// nil as well, just the number of available statistics.
//
// values is the array to be filled in with requested device
// statistics. Must not be nil if ids are specified (not nil).
//
// Returns:
//
// A positive value lower or equal to len(values): success. The return
// value is the number of entries filled in the stats table.
//
// A positive value higher than len(values): success: The given
// statistics table is too small. The return value corresponds to the
// size that should be given to succeed. The entries in the table are
// not valid and shall not be used by the caller.
//
// Otherwise, error is returned.
func (pid Port) XstatGetByID(ids, values []uint64) (int, error) {
	cids := (*C.uint64_t)(nil)
	cvalues := cids

	if len(ids) != 0 {
		cids = (*C.uint64_t)(&ids[0])
	}

	if len(values) != 0 {
		cvalues = (*C.uint64_t)(&values[0])
	}

	rc := C.rte_eth_xstats_get_by_id(C.ushort(pid), cids, cvalues, C.uint(len(values)))
	if rc < 0 {
		return 0, errget(rc)
	}

	return int(rc), nil
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
