package lcore_test

import (
	// "errors"
	"fmt"
	"sync"
	"testing"
	// "time"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/lcore"
)

func TestNewThread(t *testing.T) {
	thd := lcore.NewLockedThread(make(chan func()))
	defer thd.Close()
	err := thd.SetAffinity(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var s unix.CPUSet
	var wg sync.WaitGroup
	wg.Add(1)
	thd.Exec(false, func() {
		defer wg.Done()
		err = unix.SchedGetaffinity(0, &s)
	})
	wg.Wait()
	if err != nil || !s.IsSet(0) || s.Count() != 1 {
		t.Error("core did not launch")
		t.FailNow()
	}

	// execute and wait
	thd.Exec(true, func() {
		s = unix.CPUSet{}
	})

	if s.Count() != 0 {
		t.Error("core did not launch")
		t.FailNow()
	}
}

func TestNewThreadFail(t *testing.T) {
	thd := lcore.NewLockedThread(make(chan func()))
	defer thd.Close()
	err := thd.SetAffinity(64)
	if err == nil {
		t.Error(err)
		t.FailNow()
	}
}

func ExampleThread_Exec() {
	// create a thread with new channel
	thd := lcore.NewLockedThread(make(chan func()))
	defer thd.Close()

	var a int
	thd.Exec(true, func() { a = 1 })
	fmt.Println(a)
	// Output: 1
}

func ExampleNewLockedThread() {
	// create a thread with new channel
	thd := lcore.NewLockedThread(make(chan func()))
	defer thd.Close()

	// Set affinity to lcore 0
	err := thd.SetAffinity(0)
	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var a int
	thd.Exec(false, func() {
		defer wg.Done()
		a = 1
	})
	wg.Wait()

	fmt.Println(a)
	// Output: 1
}
