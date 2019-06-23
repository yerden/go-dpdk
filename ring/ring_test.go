package ring_test

import (
	"sync"
	"testing"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ring"
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

func TestRingCreate(t *testing.T) {
	assert := common.Assert(t, true)
	initEAL(t)

	eal.ExecuteOnMaster(func(lc *eal.Lcore) {
		r, err := ring.Create("test_ring", 1024, ring.OptSC,
			ring.OptSP, ring.OptSocket(lc.SocketID))
		assert(r != nil && err == nil, err)
		defer r.Free()
	})
}

func TestRingInit(t *testing.T) {
	assert := common.Assert(t, true)
	initEAL(t)

	eal.ExecuteOnMaster(func(lc *eal.Lcore) {
		_, err := ring.GetMemSize(1023)
		assert(err != nil) // invalid count
		n, err := ring.GetMemSize(1024)
		assert(n > 0 && err == nil, err)

		ringData := make([]byte, n)
		r := (*ring.Ring)(unsafe.Pointer(&ringData[0]))
		err = r.Init("test_ring", 1024, ring.OptSC,
			ring.OptSP, ring.OptSocket(lc.SocketID))
		assert(err == nil, err)
	})
}

func TestRingNew(t *testing.T) {
	assert := common.Assert(t, true)

	r, err := ring.New("test_ring", 1024)
	assert(r != nil && err == nil, err)

	objIn := uintptr(0xAABBCCDD)
	assert(r.SpEnqueue(objIn))
	objOut, ok := r.ScDequeue()
	assert(objIn == objOut && ok)
	_, ok = r.ScDequeue()
	assert(!ok)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func benchmarkRingUintptr(b *testing.B, burst int) {
	var wg sync.WaitGroup
	initEAL(b)
	assert := common.Assert(b, true)

	r, err := ring.New("hello", 1024)
	assert(r != nil && err == nil, err)

	sender := func(n int) {
		defer wg.Done()
		i := 0
		buf := make([]uintptr, burst)
		for i < n {
			k, _ := r.SpEnqueueBulk(buf[:burst])
			i += int(k)
		}
	}

	receiver := func(n int) {
		defer wg.Done()
		buf := make([]uintptr, burst)
		i := 0
		for i < n {
			k, _ := r.ScDequeueBulk(buf)
			i += int(k)
		}
	}

	wg.Add(2)
	go sender(b.N)
	go receiver(b.N)
	wg.Wait()
}

func BenchmarkRingUintptr1(b *testing.B) {
	benchmarkRingUintptr(b, 1)
}

func BenchmarkRingUintptr20(b *testing.B) {
	benchmarkRingUintptr(b, 20)
}

func BenchmarkRingUintptr128(b *testing.B) {
	benchmarkRingUintptr(b, 128)
}

func BenchmarkRingUintptr512(b *testing.B) {
	benchmarkRingUintptr(b, 512)
}

func BenchmarkChanNonblockUintptr(b *testing.B) {
	var wg sync.WaitGroup

	ch := make(chan uintptr, 1024)

	sender := func(n int) {
		defer wg.Done()
		i := 0
		for i < n {
			select {
			case ch <- uintptr(0xAABBCCDD):
				i++
			default:
			}
		}
	}

	receiver := func(n int) {
		defer wg.Done()
		i := 0
		for i < n {
			select {
			case <-ch:
				i++
			default:
			}
		}
	}

	wg.Add(2)
	go sender(b.N)
	go receiver(b.N)
	wg.Wait()
}

func BenchmarkChanBlockUintptr(b *testing.B) {
	var wg sync.WaitGroup

	ch := make(chan uintptr, 1024)

	sender := func(n int) {
		defer wg.Done()
		i := 0
		for i < n {
			ch <- uintptr(0xAABBCCDD)
			i++
		}
	}

	receiver := func(n int) {
		defer wg.Done()
		i := 0
		for i < n {
			<-ch
			i++
		}
	}

	wg.Add(2)
	go sender(b.N)
	go receiver(b.N)
	wg.Wait()
}
