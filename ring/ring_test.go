package ring_test

import (
	"sync"
	"syscall"
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
		err = eal.InitWithOpts(eal.OptLcores(&set), eal.OptMemory(128),
			eal.OptNoHuge, eal.OptNoPCI)
		assert(err == nil, err)
	})
}

func TestRingCreate(t *testing.T) {
	assert := common.Assert(t, true)
	initEAL(t)
	var wg sync.WaitGroup

	wg.Add(1)
	eal.ExecuteOnMaster(func(lc *eal.Lcore) {
		defer wg.Done()
		r, err := ring.Create("test_ring", 1024, ring.OptSC,
			ring.OptSP, ring.OptSocket(lc.SocketID()))
		assert(r != nil && err == nil, err)
		defer r.Free()
	})
	wg.Wait()
}

func TestRingInit(t *testing.T) {
	assert := common.Assert(t, true)
	initEAL(t)
	var wg sync.WaitGroup

	wg.Add(1)
	eal.ExecuteOnMaster(func(lc *eal.Lcore) {
		defer wg.Done()
		_, err := ring.GetMemSize(1023)
		assert(err == syscall.EINVAL) // invalid count
		n, err := ring.GetMemSize(1024)
		assert(n > 0 && err == nil, err)

		ringData := make([]byte, n)
		r := (*ring.Ring)(unsafe.Pointer(&ringData[0]))
		err = r.Init("test_ring", 1024, ring.OptSC,
			ring.OptSP, ring.OptSocket(lc.SocketID()))
		assert(err == nil, err)
	})
	wg.Wait()
}

func TestRingNew(t *testing.T) {
	assert := common.Assert(t, true)

	n := 1024
	r, err := ring.New("test_ring", uint(n))
	assert(r != nil && err == nil, err)
	defer r.Free() // should have no effect

	var objIn uintptr

	assert(r.IsEmpty())
	assert(!r.IsFull())
	assert(r.Cap() == uint(n)-1, r.Cap())

	for i := 0; i < int(r.Cap()); i++ {
		objIn = uintptr(i)
		assert(r.SpEnqueue(objIn), i)
	}

	assert(!r.IsEmpty())
	assert(r.IsFull())
	for i := 0; i < int(r.Cap()); i++ {
		objOut, ok := r.ScDequeue()
		assert(uintptr(i) == objOut && ok)
	}

	assert(r.IsEmpty())
	assert(!r.IsFull())
	_, ok := r.ScDequeue()
	assert(!ok)
}

func TestRingNewErr(t *testing.T) {
	assert := common.Assert(t, true)

	r, err := ring.New("test_ring", 1023)
	assert(r == nil && err == syscall.EINVAL)
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
