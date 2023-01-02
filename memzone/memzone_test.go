package memzone_test

import (
	"syscall"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/memzone"
)

func TestMemzoneCreateErr(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	var mz *memzone.Memzone
	var err error
	// create and test mempool on main lcore
	execErr := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := 100000000
		mz, err = memzone.Reserve("test_mz",
			uintptr(n),            // size of zone
			memzone.OptSocket(10), // incorrect, ctx.SocketID()),
			memzone.OptFlag(memzone.PageSizeHintOnly))
	})
	if eal.HasHugePages() {
		assert(mz == nil && err == syscall.ENOMEM, mz, err)
	} else {
		assert(mz != nil && err == nil)
		mz.Free()
	}
	assert(execErr == nil, execErr)
}

func TestMemzoneCreate(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := 100000000
		mz, err := memzone.Reserve("test_mz",
			uintptr(n), // size of zone
			memzone.OptSocket(eal.SocketID()),
			memzone.OptFlag(memzone.PageSizeHintOnly))
		assert(mz != nil && err == nil, err)

		mz1, err := memzone.Lookup("test_mz")
		assert(mz == mz1 && err == nil)

		err = mz.Free()
		assert(err == nil, err)

		_, err = memzone.Lookup("test_mz")
		assert(err != nil)

	})
	assert(err == nil, err)
}

func TestMemzoneWriteTo(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := 100000000
		mz, err := memzone.Reserve("test_mz",
			uintptr(n), // size of zone
			memzone.OptSocket(eal.SocketID()),
			memzone.OptFlag(memzone.PageSizeHintOnly))
		assert(mz != nil && err == nil, err)
		defer mz.Free()

		b := mz.Bytes()
		assert(len(b) == n)
		assert(n == copy(b, make([]byte, n)))
		assert("test_mz" == mz.Name())
	})
	assert(err == nil, err)
}

func TestMemzoneAligned(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := 100000000
		mz, err := memzone.Reserve("test_mz",
			uintptr(n), // size of zone
			memzone.OptSocket(eal.SocketID()),
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
	assert(err == nil, err)
}
