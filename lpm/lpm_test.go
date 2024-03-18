package lpm_test

import (
	"net/netip"
	"syscall"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/lpm"
)

var cidrs = []string{
	"192.168.0.1/24",
	"10.0.1.2/8",
	"172.16.0.30/27",
	"2001:db8:a0b:12f0::1/32",
}

func TestLpmCreate(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	var myaddr *lpm.LPM
	var e error
	cfg := &lpm.Config{
		MaxRules:    128,
		NumberTbl8s: 1 << 8,
	}
	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		myaddr, e = lpm.Create("test_lpm", -1, cfg)
	})
	assert(err == nil, err)
	assert(e == nil, e)
	defer myaddr.Free()

	// find existing
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		e = lpm.FindExisting("test_lpm_lalal", &myaddr)
		if e != syscall.ENOENT {
			return
		}
		e = lpm.FindExisting("test_lpm", &myaddr)
	})
	assert(err == nil, err)
	assert(e == nil, e)

	// populate LPM object with IPv4
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		for _, s := range cidrs {
			if ipnet := netip.MustParsePrefix(s); ipnet.Addr().Is4() {
				if e = myaddr.Add(ipnet, 1); e != nil {
					return
				}
			}
		}
	})
	assert(err == nil, err)
	assert(e == nil, err)

	// test lookup
	var hop uint32
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		if hop, e = myaddr.Lookup(netip.AddrFrom4([4]byte{10, 0, 0, 0})); e != nil {
			return
		}
		if hop, e = myaddr.Lookup(netip.AddrFrom4([4]byte{10, 1, 0, 0})); e != nil {
			return
		}
		if hop, e = myaddr.Lookup(netip.AddrFrom4([4]byte{10, 2, 0, 0})); e != nil {
			return
		}
	})
	assert(err == nil && e == nil, err, e)
	assert(hop == 1)

	// test lookup
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		hop, e = myaddr.Lookup(netip.AddrFrom4([4]byte{11, 0, 0, 0}))
	})
	assert(err == nil && e == syscall.ENOENT, err, e)
}

func TestLpm6Create(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	var myaddr *lpm.LPM6
	var e error
	cfg := &lpm.Config6{
		MaxRules:    128,
		NumberTbl8s: 1 << 8,
	}
	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		myaddr, e = lpm.Create6("test_lpm6", -1, cfg)
	})
	assert(err == nil, err)
	assert(e == nil, e)
	defer myaddr.Free()

	// find existing
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		e = lpm.FindExisting("test_lpm6_lalal", &myaddr)
		if e != syscall.ENOENT {
			return
		}
		e = lpm.FindExisting("test_lpm6", &myaddr)
	})
	assert(err == nil, err)
	assert(e == nil, e)

	// populate LPM object with IPv6
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		for _, s := range cidrs {
			if ipnet := netip.MustParsePrefix(s); ipnet.Addr().Is6() {
				if e = myaddr.Add(ipnet, 1); e != nil {
					return
				}
			}
		}
	})
	assert(err == nil, err)
	assert(e == nil, err)

	// test lookup
	var hop uint32
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		if hop, e = myaddr.Lookup(netip.MustParseAddr("2001:db8:a0b:12f0::1")); e != nil {
			return
		}
		if hop, e = myaddr.Lookup(netip.MustParseAddr("2001:db8:a0b:12f0::2")); e != nil {
			return
		}
		if hop, e = myaddr.Lookup(netip.MustParseAddr("2001:db8:a0b:12f0:0a0a::2")); e != nil {
			return
		}
	})
	assert(err == nil && e == nil, err, e)
	assert(hop == 1)

	// test lookup
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		hop, e = myaddr.Lookup(netip.MustParseAddr("2001:db9:a0b:12f0:0a0a::2"))
	})
	assert(err == nil && e == syscall.ENOENT, err, e)
}

func ExampleFindExisting() {
	// find LPM for IPv4 lookups.
	var myLPM *lpm.LPM
	if err := lpm.FindExisting("my_lpm_object", &myLPM); err != nil {
		panic(err)
	}

	// find LPM6 for IPv6 lookups.
	var myLPM6 *lpm.LPM6
	if err := lpm.FindExisting("my_lpm6_object", &myLPM6); err != nil {
		panic(err)
	}
}
