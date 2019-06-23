package eal

import (
	"sync"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"golang.org/x/sys/unix"
)

func TestOptions(t *testing.T) {
	assert := common.Assert(t, true)
	var opts ealOptions

	// equivalent of "-c f -n 4 --socket-mem=1024,1024 --no-pci -d /path/to/so"
	options := []Option{
		OptLcores(MakeSet([]int{4, 5, 7, 9})),
		OptLcores(MakeSet([]int{0, 1, 2, 3})),
		OptMemoryChannels(4),
		OptSocketMemory(1024, 1024),
		OptNoPCI,
		OptLoadExternalPath("/path/to/so"),
	}

	for i := range options {
		options[i].f(&opts)
	}

	// first arg is the name of this executable
	argv := opts.argv()
	assert(len(argv) == 8)
	assert(argv[0] == "-c")
	assert(argv[1] == "f")
	assert(argv[2] == "-n")
	assert(argv[3] == "4")
	assert(argv[4] == "--socket-mem=1024,1024")
	assert(argv[5] == "--no-pci")
	assert(argv[6] == "-d")
	assert(argv[7] == "/path/to/so")
}

func TestEALInit(t *testing.T) {
	assert := common.Assert(t, true)
	var set unix.CPUSet
	assert(unix.SchedGetaffinity(0, &set) == nil)
	err := InitWithOpts(OptLcores(&set), OptNoHuge, OptNoPCI, OptMasterLcore(0))
	assert(err == nil)

	ch := make(chan uint, set.Count())
	assert(LcoreCount() == uint(set.Count()))
	var wg sync.WaitGroup
	ForeachLcore(false, func(lcoreId uint) {
		wg.Add(1)
		go ExecuteOnLcore(lcoreId, func(lc *Lcore) {
			defer wg.Done()
			assert(lc.ID == LcoreID())
			ch <- lc.ID
		})
	})
	wg.Wait()

	var myset unix.CPUSet
	for i := 0; i < set.Count(); i++ {
		myset.Set(int(<-ch))
	}

	select {
	case <-ch:
		assert(false)
	default:
	}

	assert(myset == set)

	// test panic
	ForeachLcore(false, func(lcoreID uint) {
		wg.Add(1)
		ExecuteOnLcore(lcoreID, func(lc *Lcore) {
			defer wg.Done()
			panic("emit panic")
		})
	})
	wg.Wait()
	ForeachLcore(false, func(lcoreID uint) {
		wg.Add(1)
		ExecuteOnLcore(lcoreID, func(lc *Lcore) {
			// lcore is fine
			defer wg.Done()
		})
	})
	wg.Wait()
}
