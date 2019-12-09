package memzone_test

import (
	"sync"
	"syscall"
	"testing"
	// "unsafe"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/memzone"
)

var (
	dpdk    sync.Once
	dpdkErr error
)

func initEAL() {
	dpdk.Do(func() {
		var set unix.CPUSet
		err := unix.SchedGetaffinity(0, &set)
		if err == nil {
			err = eal.InitWithParams(
				eal.NewParameter("-c", eal.NewMap(&set)),
				eal.NewParameter("-m", "128"),
				eal.NewParameter("--no-huge"),
				eal.NewParameter("--no-pci"),
			)
		}
		dpdkErr = err
	})
}

func TestMemzoneCreateErr(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	initEAL()
	assert(dpdkErr == nil, dpdkErr)

	var wg sync.WaitGroup
	wg.Add(1)
	// create and test mempool on master lcore
	eal.ExecuteOnMaster(func(ctx *eal.Lcore) {
		defer wg.Done()
		// create empty mempool
		n := 100000000
		mz, err := memzone.Reserve("test_mz",
			uintptr(n),            // size of zone
			memzone.OptSocket(10), // incorrect, ctx.SocketID()),
			memzone.OptFlag(memzone.PageSizeHintOnly))
		assert(mz == nil && err == syscall.ENOMEM, mz, err)
	})
	wg.Wait()
}

func TestMemzoneCreate(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	initEAL()
	assert(dpdkErr == nil, dpdkErr)

	var wg sync.WaitGroup
	wg.Add(1)
	// create and test mempool on master lcore
	eal.ExecuteOnMaster(func(ctx *eal.Lcore) {
		defer wg.Done()
		// create empty mempool
		n := 100000000
		mz, err := memzone.Reserve("test_mz",
			uintptr(n), // size of zone
			memzone.OptSocket(ctx.SocketID()),
			memzone.OptFlag(memzone.PageSizeHintOnly))
		assert(mz != nil && err == nil, err)

		mz1, err := memzone.Lookup("test_mz")
		assert(mz == mz1 && err == nil)

		err = mz.Free()
		assert(err == nil, err)

		_, err = memzone.Lookup("test_mz")
		assert(err != nil)

	})
	wg.Wait()
}

func TestMemzoneWriteTo(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	initEAL()
	assert(dpdkErr == nil, dpdkErr)

	var wg sync.WaitGroup
	wg.Add(1)
	// create and test mempool on master lcore
	eal.ExecuteOnMaster(func(ctx *eal.Lcore) {
		defer wg.Done()
		// create empty mempool
		n := 100000000
		mz, err := memzone.Reserve("test_mz",
			uintptr(n), // size of zone
			memzone.OptSocket(ctx.SocketID()),
			memzone.OptFlag(memzone.PageSizeHintOnly))
		assert(mz != nil && err == nil, err)
		defer mz.Free()

		b := mz.Bytes()
		assert(len(b) == n)
		assert(n == copy(b, make([]byte, n)))
		assert("test_mz" == mz.Name())
	})
	wg.Wait()
}

func TestMemzoneAligned(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	initEAL()
	assert(dpdkErr == nil, dpdkErr)

	var wg sync.WaitGroup
	wg.Add(1)
	// create and test mempool on master lcore
	eal.ExecuteOnMaster(func(ctx *eal.Lcore) {
		defer wg.Done()
		// create empty mempool
		n := 100000000
		mz, err := memzone.Reserve("test_mz",
			uintptr(n), // size of zone
			memzone.OptSocket(ctx.SocketID()),
			memzone.OptFlag(memzone.PageSizeHintOnly),
			memzone.OptAligned(1024))
		assert(mz != nil && err == nil, err)
		defer mz.Free()

		b := mz.Bytes()
		assert(len(b) == n)
		assert(n == copy(b, make([]byte, n)))
		assert("test_mz" == mz.Name())

		var mz1 *memzone.Memzone
		memzone.Walk(func(mz *memzone.Memzone) {
			if mz.Name() == "test_mz" {
				mz1 = mz
			}
		})
		assert(mz == mz1, mz1)
	})
	wg.Wait()
}
