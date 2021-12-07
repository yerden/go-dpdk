package mempool_test

import (
	"bytes"
	"syscall"
	"testing"

	"github.com/yerden/go-dpdk/mbuf"
	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
)

var initEAL = common.DoOnce(func() error {
	var set unix.CPUSet
	err := unix.SchedGetaffinity(0, &set)
	if err == nil {
		_, err = eal.Init([]string{"test",
			"-c", common.NewMap(&set).String(),
			"-d", eal.PmdPath,
			"-m", "128",
			"--no-huge",
			"--no-pci",
			"--main-lcore", "0"})
	}
	return err
})

func TestMempoolCreateErr(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateEmpty("test_mp",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32000000), // too large
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(1024),
		)
		assert(mp == nil && err == syscall.EINVAL, err)
	})
	assert(err == nil, err)
}

func TestMempoolPriv(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateEmpty("test_mp",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32), // too large
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(1024),
		)
		assert(mp != nil && err == nil, err)
		defer mp.Free()

		priv := mp.GetPrivBytes()
		assert(len(priv) == 1024, len(priv))
	})
	assert(err == nil, err)
}

func TestMempoolCreate(t *testing.T) {
	assert := common.Assert(t, false)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateEmpty("test_mp",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			//mempool.OptOpsName("stack"),
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
	assert(err == nil, err)

	// create and test mempool on main lcore
	err = eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateMbufPool("test_mbuf_pool",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(64), // for each mbuf
		)
		assert(err == nil, err)
		assert(mp != nil)
		defer mp.Free()

		data := []byte("hello from mbuf")
		myMbuf := mbuf.PktMbufAlloc(mp)
		mbuf.PktMbufAppend(myMbuf, data)
		assert(bytes.Equal(myMbuf.Data(), data))
		memp := myMbuf.GetPool()
		mymp, err := mempool.Lookup("test_mbuf_pool")
		assert(mymp == memp)
		assert(err == nil)
		defer mbuf.PktMbufFree(myMbuf)

		var mbufArr []*mbuf.Mbuf
		mbufArr = make([]*mbuf.Mbuf, 4)
		err = mbuf.PktMbufAllocBulk(mp, mbufArr)
		assert(err == nil)
		for _, m := range mbufArr {
			mbuf.PktMbufAppend(m, data)
			assert(bytes.Equal(m.Data(), data))
		}

		for _, m := range mbufArr {
			mbuf.PktMbufReset(m)
			assert(bytes.Equal(m.Data(), []byte{}))
		}

		var mbufArrEmpty []*mbuf.Mbuf
		err = mbuf.PktMbufAllocBulk(mp, mbufArrEmpty)
		assert(err == nil)
		assert(len(mbufArrEmpty) == 0)

		mp, err = mempool.CreateMbufPool("test_mbuf_pool_err",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(63), // for each mbuf
		)
		assert(err != nil, err)
	})
	assert(err == nil, err)
}

func TestMbufpoolPriv(t *testing.T) {
	assert := common.Assert(t, false)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateMbufPool("test_mbuf_pool_priv",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(64), // for each mbuf
		)
		assert(err == nil, err)
		assert(mp != nil)
		defer mp.Free()
		var mpPrivData mbuf.MpPrivateData
		mpPrivData.SetFrom(mp)

		data := []byte("hello from private area")
		myMbuf := mbuf.PktMbufAlloc(mp)
		mData := mpPrivData.PrivData(myMbuf)
		assert(len(mData) == int(mpPrivData.MbufPrivSize))

		copy(mData, data)
		newData := mpPrivData.PrivData(myMbuf)
		assert(bytes.Equal(data, newData[:len(data)]))
		mbuf.PktMbufFree(myMbuf)
	})
	assert(err == nil, err)
}
