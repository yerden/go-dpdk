package mempool_test

import (
	"bytes"
	"log"

	"reflect"
	"syscall"
	"testing"
	"unsafe"

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
	assert(t, err == nil, err)
}

type TestStruct struct {
	a int
	b string
}

func doOnMainStruct(t *testing.T, fn func(p *mempool.Mempool, str []TestStruct)) {
	// Initialize EAL on all cores
	assert(t, initEAL() == nil)

	data := []TestStruct{}
	for i := 0; i < 10; i++ {
		data = append(data, TestStruct{
			a: i,
			b: "asd",
		})
	}

	// create and test mempool on main lcore
	err := eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
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
	assert(t, err == nil, err)
}

func TestGetGoStruct(t *testing.T) {
	strSize := int(unsafe.Sizeof(TestStruct{}))

	const expectedTotalLengthOfSlice = 240

	doOnMainStruct(t, func(p *mempool.Mempool, data []TestStruct) {
		cstr := &common.CStruct{}
		cstr.Ptr = unsafe.Pointer(&data[0])
		cstr.Len = len(data) * strSize
		myMbuf := mbuf.AllocResetAndAppend(p, cstr)

		mbufData := myMbuf.Data()

		if len(myMbuf.Data()) == expectedTotalLengthOfSlice {
			// cast slice to byte array
			//byteArray = *(*[10]byte)(myMbuf.Data())
			// cast array to array of structs
			expectedSlice := *(*[10]TestStruct)(unsafe.Pointer(&mbufData[0]))
			log.Println(expectedSlice)
		} else {
			// other length
		}
	})
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
		priv := p.PrivBytes()
		assert(t, len(priv) == 64, len(priv))
	})
}

func TestResetAndAppendErr(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// test with slice of byte
		cArr := &common.CStruct{}
		cArr.Init(unsafe.Pointer(&data[0]), 2049)

		myMbuf := mbuf.PktMbufAlloc(p)
		defer myMbuf.PktMbufFree()
		assert(t, myMbuf != nil)
		err := myMbuf.ResetAndAppend(cArr)
		assert(t, err != nil, err)
	})
}

func TestAllocResetAndAppendErr(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// test with slice of byte
		cArr := &common.CStruct{}
		cArr.Init(unsafe.Pointer(&data[0]), 2049)

		myMbuf := mbuf.AllocResetAndAppend(p, cArr)
		defer myMbuf.PktMbufFree()
		assert(t, myMbuf == nil)
	})
}

func TestHeadRoomSize(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// allocate mbuf from mempool
		myMbuf := mbuf.PktMbufAlloc(p)

		// reset the fields
		myMbuf.PktMbufReset()
		assert(t, myMbuf.BufLen() == 2048)
		assert(t, myMbuf.HeadRoomSize() == 128)
		assert(t, myMbuf.HeadRoomSize() == myMbuf.PktMbufHeadRoomSize())
	})
}

func TestTailRoomSize(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// allocate mbuf from mempool
		myMbuf := mbuf.PktMbufAlloc(p)

		// reset the fields
		myMbuf.PktMbufReset()
		assert(t, myMbuf.BufLen() == 2048)
		assert(t, myMbuf.TailRoomSize() == 1920)
		assert(t, myMbuf.TailRoomSize() == myMbuf.PktMbufTailRoomSize())
		assert(t, myMbuf.TailRoomSize() == myMbuf.BufLen()-myMbuf.HeadRoomSize()-uint16(len(myMbuf.Data())))
	})
}

func TestBufLen(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// allocate mbuf from mempool
		myMbuf := mbuf.PktMbufAlloc(p)

		// reset the fields
		myMbuf.PktMbufReset()
		assert(t, myMbuf.BufLen() == 2048)
	})
}

func TestMbufMethods(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// allocate mbuf from mempool
		myMbuf := mbuf.PktMbufAlloc(p)

		// reset the fields
		myMbuf.PktMbufReset()
		assert(t, myMbuf.BufLen() == 2048)
		assert(t, myMbuf.HeadRoomSize() == 128)
		assert(t, myMbuf.TailRoomSize() == myMbuf.BufLen()-myMbuf.HeadRoomSize()-uint16(len(myMbuf.Data())))

		// append data to the mbuf
		myMbuf.PktMbufAppend(data)
		assert(t, myMbuf.BufLen() == 2048)
		assert(t, myMbuf.HeadRoomSize() == 128)
		assert(t, myMbuf.TailRoomSize() == myMbuf.BufLen()-myMbuf.HeadRoomSize()-uint16(len(myMbuf.Data())))
		assert(t, bytes.Equal(myMbuf.Data(), data))

		memp := myMbuf.Mempool()
		mymp, err := mempool.Lookup("test_mbuf_pool")
		assert(t, mymp == memp)
		assert(t, err == nil)
		defer myMbuf.PktMbufFree()

		// allocate a bulk of mbufs and append the data to them
		var mbufArr []*mbuf.Mbuf
		mbufArr = make([]*mbuf.Mbuf, 4)
		err = mbuf.PktMbufAllocBulk(p, mbufArr)
		assert(t, err == syscall.Errno(0))
		for _, m := range mbufArr {
			m.PktMbufAppend(data)
			assert(t, bytes.Equal(m.Data(), data))
		}

		// reset the fields
		for _, m := range mbufArr {
			m.PktMbufReset()
			assert(t, bytes.Equal(m.Data(), []byte{}))
		}

		// allocation to the empty array
		var mbufArrEmpty []*mbuf.Mbuf
		err = mbuf.PktMbufAllocBulk(p, mbufArrEmpty)
		assert(t, err == syscall.Errno(0))
		assert(t, len(mbufArrEmpty) == 0)
	})
}

func TestMbufpoolPriv(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// alloc mbuf from mempool
		myMbuf := mbuf.PktMbufAlloc(p)

		mData := &common.CStruct{}
		// get the content of mbuf private area specified by pointer and len
		myMbuf.PrivData(mData)
		assert(t, mData.Len == int(myMbuf.PrivSize()))

		// create byte slice from pointer to the private area of mbuf
		ptrSlice := mData.Bytes()
		// copy the data to mbufs private area
		copy(ptrSlice, data)

		newData := &common.CStruct{}
		myMbuf.PrivData(newData)
		newSlice := newData.Bytes()
		assert(t, bytes.Equal(data, newSlice[:len(data)]))
		assert(t, newData.Len == int(myMbuf.PrivSize()))
		myMbuf.PktMbufFree()
	})
}

type someStruct struct {
	intField    int
	stringField string
	bytes4Field [4]byte
	uint8Fields uint8
}

func TestAllocResetAppend(t *testing.T) {
	doOnMain(t, func(p *mempool.Mempool, data []byte) {
		// test with slice of byte
		cArr := &common.CStruct{}
		cArr.Init(unsafe.Pointer(&data[0]), len(data))

		myMbuf := mbuf.PktMbufAlloc(p)
		assert(t, myMbuf != nil)
		err := myMbuf.ResetAndAppend(cArr)
		assert(t, err == nil)
		assert(t, len(myMbuf.Data()) == len(data))

		d := myMbuf.Data()
		buf := *(*[]byte)(unsafe.Pointer(&d))
		assert(t, bytes.Equal(buf, data))
		myMbuf.PktMbufFree()

		newMbuf := mbuf.AllocResetAndAppend(p, cArr)
		assert(t, newMbuf != nil)
		assert(t, len(newMbuf.Data()) == len(data))
		newMbuf.PktMbufFree()

		// test with struct
		testdata := someStruct{
			intField:    250,
			stringField: "hello from mbuf",
			bytes4Field: [4]byte{10, 15, 20, 30},
			uint8Fields: 128,
		}
		var arr []someStruct
		for i := 0; i < 2; i++ {
			arr = append(arr, testdata)
		}
		cArr.Ptr = unsafe.Pointer(&arr[0])
		cArr.Len = int(uintptr(len(arr)) * reflect.TypeOf(arr).Elem().Size())

		m := mbuf.AllocResetAndAppend(p, cArr)
		assert(t, m != nil)
		assert(t, len(m.Data()) == cArr.Len)

		d = m.Data()
		str := *(*[]someStruct)(unsafe.Pointer(&d))
		assert(t, str[0] == testdata)
		assert(t, str[1] == testdata)
	})

}

func BenchmarkAllocFromChannel(b *testing.B) {
	initEAL()
	n := uint32(10240)
	ch := make(chan *mbuf.Mbuf, n)
	data := []byte("Some data for test")
	cArr := &common.CStruct{}
	cArr.Init(unsafe.Pointer(&data[0]), len(data))
	defer close(ch)
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
			myMbuf.ResetAndAppend(cArr)
			ch <- myMbuf
		}
	})
}

func BenchmarkAllocFromMempool(b *testing.B) {
	initEAL()
	n := uint32(10240)
	data := []byte("Some data for test")
	cArr := &common.CStruct{}
	cArr.Init(unsafe.Pointer(&data[0]), len(data))
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
			newMbuf := mbuf.AllocResetAndAppend(mp, cArr)
			newMbuf.PktMbufFree()
		}
	})
}
