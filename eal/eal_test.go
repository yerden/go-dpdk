package eal

import (
	_ "log"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"golang.org/x/sys/unix"
)

func TestEALInit(t *testing.T) {
	assert := common.Assert(t, true)
	var set unix.CPUSet
	assert(unix.SchedGetaffinity(0, &set) == nil)

	n, err := Init([]string{"test",
		"-c", common.NewMap(&set).String(),
		"-m", "128",
		"--no-huge",
		"--no-pci",
		"--master-lcore", "0"})
	assert(n == 8, n)
	assert(err == nil)

	ch := make(chan uint, set.Count())
	assert(LcoreCount() == uint(set.Count()))
	for _, id := range Lcores() {
		ExecOnLcore(id, func(id uint) func(*LcoreCtx) {
			return func(ctx *LcoreCtx) {
				assert(id == ctx.LcoreID())
				ch <- ctx.LcoreID()
			}
		}(id))
	}

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
	for _, id := range Lcores() {
		err := ExecOnLcore(id, func(ctx *LcoreCtx) {
			panic("emit panic")
		})
		e, ok := err.(*ErrLcorePanic)
		assert(ok)
		assert(e.LcoreID == id)
		assert(len(e.Pc) > 0)
	}
	for _, id := range Lcores() {
		ok := false
		err := ExecOnLcore(id, func(ctx *LcoreCtx) {
			// lcore is fine
			ok = true
		})
		assert(ok && err == nil, err)
	}

	// invalid lcore
	assert(ExecOnLcore(uint(1024), func(ctx *LcoreCtx) {}) == ErrLcoreInvalid)

	// stop all lcores
	StopLcores()

	err = Cleanup()
	assert(err == nil, err)
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
