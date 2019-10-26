package eal

import (
	_ "log"
	"sync"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"golang.org/x/sys/unix"
)

func TestEALInit(t *testing.T) {
	assert := common.Assert(t, true)
	var set unix.CPUSet
	assert(unix.SchedGetaffinity(0, &set) == nil)

	err := InitWithParams(
		NewParameter("-c", NewMap(&set)),
		NewParameter("--no-huge"),
		NewParameter("--no-pci"),
		NewParameter("--master-lcore", "0"),
	)
	assert(err == nil)

	ch := make(chan uint, set.Count())
	assert(LcoreCount() == uint(set.Count()))
	var wg sync.WaitGroup
	for _, id := range Lcores(false) {
		wg.Add(1)
		ExecuteOnLcore(id, func(id uint) func(*Lcore) {
			return func(lc *Lcore) {
				defer wg.Done()
				assert(id == lc.ID())
				ch <- lc.ID()
			}
		}(id))
	}
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
	for _, id := range Lcores(false) {
		wg.Add(1)
		ExecuteOnLcore(id, func(lc *Lcore) {
			defer wg.Done()
			panic("emit panic")
		})
	}
	wg.Wait()
	for _, id := range Lcores(false) {
		wg.Add(1)
		ExecuteOnLcore(id, func(lc *Lcore) {
			// lcore is fine
			defer wg.Done()
		})
	}
	wg.Wait()
}

func TestParseCmd(t *testing.T) {
	assert := common.Assert(t, true)

	res, err := parseCmd("hello bitter world")
	assert(err == nil, err)
	assert(res[0] == "hello")
	assert(res[1] == "bitter")
	assert(res[2] == "world")

	res, err = parseCmd("hello --bitter world")
	assert(err == nil, err)
	assert(res[0] == "hello")
	assert(res[1] == "--bitter", res[1])
	assert(res[2] == "world")
}
