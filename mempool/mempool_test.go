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

func assert(t testing.TB, expected bool, args ...interface{}) {
	if !expected {
		t.Helper()
		t.Fatal(args...)
	}
}

func doOnMain(t *testing.T, fn func(p *mempool.Mempool, data []byte)) {
	// Initialize EAL on all cores
	assert(t, initEAL() == nil)

	data := []byte("hello from Mbuf")
	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		// create empty mempool
		n := uint32(10240)
		mp, err := mempool.CreateMbufPool("test_mbuf_pool",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(64), // for each Mbuf
		)
		assert(t, err == nil, err)
		assert(t, mp != nil)
		defer mp.Free()
		fn(mp, data)
	})
	assert(t, err == nil)
}

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

	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		p, err = mempool.CreateMbufPool("test_mbuf_pool_err",
			10240, // elements count
			2048,  // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(63), // for each mbuf
		)
		assert(err != nil, err)
	})
}

func TestMempoolPriv(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		priv := p.GetPrivBytes()
		assert(t, len(priv) == 64, len(priv))
	})
}

// TODO: add docs
func TestMbufMethods(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		myMbuf := mbuf.PktMbufAlloc(p)
		mbuf.PktMbufReset(myMbuf)
		mbuf.PktMbufAppend(myMbuf, data)
		assert(t, bytes.Equal(myMbuf.Data(), data))
		memp := myMbuf.GetPool()
		mymp, err := mempool.Lookup("test_mbuf_pool")
		assert(t, mymp == memp)
		assert(t, err == nil)
		defer mbuf.PktMbufFree(myMbuf)

		var mbufArr []*mbuf.Mbuf
		mbufArr = make([]*mbuf.Mbuf, 4)
		err = mbuf.PktMbufAllocBulk(p, mbufArr)
		assert(t, err == nil)
		for _, m := range mbufArr {
			mbuf.PktMbufAppend(m, data)
			assert(t, bytes.Equal(m.Data(), data))
		}

		for _, m := range mbufArr {
			mbuf.PktMbufReset(m)
			assert(t, bytes.Equal(m.Data(), []byte{}))
		}

		var mbufArrEmpty []*mbuf.Mbuf
		err = mbuf.PktMbufAllocBulk(p, mbufArrEmpty)
		assert(t, err == nil)
		assert(t, len(mbufArrEmpty) == 0)
	})
}

func TestMbufpoolPriv(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		myMbuf := mbuf.PktMbufAlloc(p)
		mData := myMbuf.GetPrivData()
		assert(t, mData.Len == int(myMbuf.GetPrivSize()))
		ptrSlice := common.MakeSlice(mData.Ptr, mData.Len)

		copy(ptrSlice, data)
		newData := myMbuf.GetPrivData()
		newSlice := common.MakeSlice(newData.Ptr, newData.Len)
		assert(t, bytes.Equal(data, newSlice[:len(data)]))
		assert(t, newData.Len == int(myMbuf.GetPrivSize()))
		mbuf.PktMbufFree(myMbuf)
	})
}

func TestAllocResetAppend(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		myMbuf := mbuf.PktMbufAlloc(p)
		myMbuf.ResetAndAppend(data)
		assert(t, bytes.Equal(myMbuf.Data(), data))

		newMbuf := mbuf.AllocResetAndAppend(p, data)
		assert(t, newMbuf != nil)
		assert(t, bytes.Equal(newMbuf.Data(), data))
	})

}

func BenchmarkAllocFromChannel(b *testing.B) {
	initEAL()
	n := uint32(10240)
	ch := make(chan *mbuf.Mbuf, n)
	data := []byte("Some data for test")
	eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		mp, _ := mempool.CreateMbufPool("test_mbuf_pool",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(64), // for each mbuf
		)
		defer mp.Free()
		mbufArr := make([]*mbuf.Mbuf, n)
		mbuf.PktMbufAllocBulk(mp, mbufArr)
		for _, m := range mbufArr {
			ch <- m
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			myMbuf := <-ch
			myMbuf.ResetAndAppend(data)
			ch <- myMbuf
		}
	})
	close(ch)
}

func BenchmarkAllocFromMempool(b *testing.B) {
	initEAL()
	n := uint32(10240)
	data := []byte("Some data for test")
	eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		mp, _ := mempool.CreateMbufPool("test_mbuf_pool",
			n,    // elements count
			2048, // size of element
			mempool.OptSocket(int(eal.SocketID())),
			mempool.OptCacheSize(32),
			mempool.OptOpsName("stack"),
			mempool.OptPrivateDataSize(64), // for each mbuf
		)
		defer mp.Free()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			newMbuf := mbuf.AllocResetAndAppend(mp, data)
			mbuf.PktMbufFree(newMbuf)
		}
	})
}
