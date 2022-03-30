package ethdev

import (
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
)

var initEAL = common.DoOnce(func() error {
	_, err := eal.Init([]string{"test",
		"--lcores", "(0-3)@0",
		"--vdev", "net_null0",
		"-d", eal.PmdPath,
		"-m", "128", "--no-huge", "--no-pci",
		"--main-lcore", "0"})
	return err
})

func assert(t testing.TB, expected bool, args ...interface{}) {
	if !expected {
		t.Helper()
		t.Fatal(args...)
	}
}

func TestEthXstats(t *testing.T) {
	assert(t, initEAL() == nil)

	var names []XstatName
	var ealErr error
	namesByID := map[uint64]string{}

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		pid := Port(0)
		names, ealErr = pid.XstatNames()
		assert(t, ealErr == nil)

		namesByID, ealErr = pid.XstatNameIDs()
		assert(t, ealErr == nil)
	})

	assert(t, err == nil, err)

	assert(t, len(names) > 0)
	assert(t, len(namesByID) == len(names), namesByID)
}

// TestEthStats tests that EthStats can be casted to Stats.
func TestEthStats(t *testing.T) {
	var srcA [10]Stats

	b := (*[unsafe.Sizeof(srcA)]byte)(unsafe.Pointer(&srcA))[:]
	rand.Seed(time.Now().UnixNano())
	rand.Read(b)

	for i := range srcA {
		src := &srcA[i]
		dst := src.Cast()
		if uint64(src.ipackets) != dst.Ipackets {
			t.Fatal("uint64(src.ipackets) != dst.Ipackets", i, uint64(src.ipackets), dst.Ipackets)
		}

		if uint64(src.opackets) != dst.Opackets {
			t.Fatal("uint64(src.opackets) != dst.Opackets", i, uint64(src.opackets), dst.Opackets)
		}

		if uint64(src.ibytes) != dst.Ibytes {
			t.Fatal("uint64(src.ibytes) != dst.Ibytes", i, uint64(src.ibytes), dst.Ibytes)
		}

		if uint64(src.obytes) != dst.Obytes {
			t.Fatal("uint64(src.obytes) != dst.Obytes", i, uint64(src.obytes), dst.Obytes)
		}

		if uint64(src.imissed) != dst.Imissed {
			t.Fatal("uint64(src.imissed) != dst.Imissed", i, uint64(src.imissed), dst.Imissed)
		}

		if uint64(src.ierrors) != dst.Ierrors {
			t.Fatal("uint64(src.ierrors) != dst.Ierrors", i, uint64(src.ierrors), dst.Ierrors)
		}

		if uint64(src.oerrors) != dst.Oerrors {
			t.Fatal("uint64(src.oerrors) != dst.Oerrors", i, uint64(src.oerrors), dst.Oerrors)
		}

		if uint64(src.rx_nombuf) != dst.RxNoMbuf {
			t.Fatal("uint64(src.rx_nombuf) != dst.RxNoMbuf", i, uint64(src.rx_nombuf), dst.RxNoMbuf)
		}
	}
}

// TestEthStats tests that EthStats can be casted to Stats.
func TestEthXstat(t *testing.T) {
	var srcA [10]cXstat

	b := (*[unsafe.Sizeof(srcA)]byte)(unsafe.Pointer(&srcA))[:]
	rand.Seed(time.Now().UnixNano())
	rand.Read(b)

	dstA := (*[10]Xstat)(unsafe.Pointer(&srcA))

	for i := range srcA {
		src := &srcA[i]
		dst := dstA[i]
		if uint64(src.id) != dst.Index {
			t.Fatal("uint64(src.id) != dst.Index", i, uint64(src.id), dst.Index)
		}
		if uint64(src.value) != dst.Value {
			t.Fatal("uint64(src.value) != dst.Value", i, uint64(src.value), dst.Value)
		}
	}
}
