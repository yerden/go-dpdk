package mempool_test

import (
	"os"
	"sync"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
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
			err = eal.InitWithParams(os.Args[0],
				eal.NewParameter("-c", eal.NewMap(&set)),
				eal.NewParameter("-m", "128"),
				eal.NewParameter("--no-huge"),
				eal.NewParameter("--no-pci"),
			)
		}
		dpdkErr = err
	})
}

func TestMempoolCreateErr(t *testing.T) {
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
		n := uint32(10240)
		mp, err := mempool.CreateEmpty("test_mp",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(ctx.SocketID())),
			mempool.OptCacheSize(32000000), // too large
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(1024),
		)
		assert(mp == nil && err == syscall.EINVAL, err)
	})
	wg.Wait()
}

func TestMempoolPriv(t *testing.T) {
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
		n := uint32(10240)
		mp, err := mempool.CreateEmpty("test_mp",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(ctx.SocketID())),
			mempool.OptCacheSize(32), // too large
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(1024),
		)
		assert(mp != nil && err == nil, err)
		defer mp.Free()

		priv := mp.GetPrivBytes()
		assert(len(priv) == 1024, len(priv))
	})
	wg.Wait()
}

func TestMempoolCreate(t *testing.T) {
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
		n := uint32(10240)
		mp, err := mempool.CreateEmpty("test_mp",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(ctx.SocketID())),
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
		n = mp.ObjIter(func(obj []byte) {
			assert(obj != nil, "obj should be non-nil")
			assert(len(obj) == 2048, len(obj))
			k++
		})
		assert(m == int(n), m, n)
		assert(k == int(n), k, n)

		// lookup
		mymp, err := mempool.Lookup("test_mp")
		assert(err == nil, err)
		assert(mymp == mp, mymp)

		// lookup err
		_, err = mempool.Lookup("test_mp_nonexistent")
		assert(err != nil, err)
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
			mempool.OptSocket(int(ctx.SocketID())),
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
			mempool.OptSocket(int(ctx.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(63), // for each mbuf
		)
		assert(err != nil, err)
	})
	wg.Wait()
}
