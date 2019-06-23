package mempool_test

import (
	"sync"
	"testing"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
)

var (
	dpdk sync.Once
)

func initEAL(t testing.TB) {
	assert := common.Assert(t, true)
	var set unix.CPUSet
	err := unix.SchedGetaffinity(0, &set)
	assert(err == nil, err)
	dpdk.Do(func() {
		err = eal.InitWithOpts(eal.OptLcores(&set), eal.OptMemory(1024),
			eal.OptNoHuge, eal.OptNoPCI)
		assert(err == nil, err)
	})
}

func TestCreateMempool(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	initEAL(t)

	var wg sync.WaitGroup
	wg.Add(1)
	// create and test mempool on master lcore
	eal.ExecuteOnMaster(func(ctx *eal.Lcore) {
		defer wg.Done()
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateEmpty("test_mp",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(ctx.SocketID)),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(1024),
		)
		assert(err == nil, err)
		assert(mp != nil)
		defer mp.Free()

		// change ops to ring
		err = mp.SetOpsByName("ring_mp_mc", nil)
		assert(err == nil, err)

		// populate by default
		m, err := mp.PopulateDefault()
		assert(err == nil, err)
		assert(m == int(n), m, n)

		// iterate all objects
		k := 0
		n = mp.ObjIter(func(obj unsafe.Pointer) {
			assert(obj != nil, "obj should be non-nil")
			k++
		})
		assert(m == int(n), m, n)
		assert(k == int(n), k, n)
	})
	wg.Wait()

	wg.Add(1)
	// create and test mempool on master lcore
	eal.ExecuteOnMaster(func(ctx *eal.Lcore) {
		defer wg.Done()
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateMbufPool("test_mbuf_pool",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(ctx.SocketID)),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(64), // for each mbuf
		)
		assert(err == nil, err)
		assert(mp != nil)
		defer mp.Free()

		mp, err = mempool.CreateMbufPool("test_mbuf_pool_err",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(ctx.SocketID)),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(63), // for each mbuf
		)
		assert(err != nil, err)
	})
	wg.Wait()
}
