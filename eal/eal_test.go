package eal

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"golang.org/x/sys/unix"
)

func TestEALInit(t *testing.T) {
	assert := common.Assert(t, true)
	var set unix.CPUSet
	assert(unix.SchedGetaffinity(0, &set) == nil)

	_, err := Init([]string{"test", "--some-invalid-option"})
	assert(err != nil)

	n, err := Init([]string{"test",
		"-c", common.NewMap(&set).String(),
		"-m", "128",
		"--no-huge",
		"--no-pci",
		"--main-lcore", "0"})
	assert(n == 8, n)
	assert(err == nil)

	ch := make(chan uint, set.Count())
	assert(LcoreCount() == uint(set.Count()))
	for _, id := range Lcores() {
		ExecOnLcore(id, func(id uint) func(*LcoreCtx) {
			return func(ctx *LcoreCtx) {
				assert(id == LcoreID())
				ch <- LcoreID()
			}
		}(id))
	}

	ExecOnMain(func(*LcoreCtx) {
		assert(HasPCI() == false)
		assert(HasHugePages() == false)
	})

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
	PanicAsErr = true

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

	// test panic returning arbitrary error
	err = ExecOnMain(func(*LcoreCtx) {
		panic(flag.ErrHelp)
	})
	assert(err != nil)
	e, ok := err.(*ErrLcorePanic)
	assert(ok)
	assert(e.Unwrap() == flag.ErrHelp)

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

func ExampleInit_flags() {
	// If your executable is to be launched like a DPDK example:
	//   /path/to/exec -c 3f -m 1024 <other-eal-options> -- <go-flags>
	// then you may do the following:
	n, err := Init(os.Args)
	if err != nil {
		panic(err)
	}

	// to be able to do further flag processing, i.e. pretend the cmd was:
	//   /path/to/exec <go-flags>
	os.Args[n], os.Args = os.Args[0], os.Args[n:]
	flag.Parse()
}

func ExampleExecOnLcore() {
	// Lcore ID 1
	lid := uint(1)

	err := ExecOnLcore(lid, func(ctx *LcoreCtx) {
		log.Printf("this is lcore #%d\n", LcoreID())
	})

	if err == ErrLcoreInvalid {
		// lid doesn't exist
		log.Fatalf("invalid lcore %d\n", lid)
	}

	if e, ok := err.(*ErrLcorePanic); ok {
		// lcore panicked
		log.Fatalln(e)
	}
}

func ExampleExecOnLcore_error() {
	// Lcore ID 1
	lid := uint(1)

	someErr := fmt.Errorf("lcore error")
	err := ExecOnLcore(lid, func(ctx *LcoreCtx) {
		if 2+2 != 4 {
			panic(someErr)
		}
	})

	if e, ok := err.(*ErrLcorePanic); ok {
		// lcore panicked
		if err := e.Unwrap(); err == someErr {
			log.Fatalln("check the math")
		}
	}

	// or, as of Go 1.13
	//   if errors.Is(err, someErr) { ...
}
