package lcore_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yerden/go-dpdk/lcore"
)

func TestNewThread(t *testing.T) {
	lt, err := lcore.NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer lt.Exit()

	var wg sync.WaitGroup
	a := 1

	wg.Add(1)
	lt.Execute(func(ctx *lcore.ThreadCtx) error {
		defer wg.Done()
		a = 2
		return nil
	})

	wg.Wait()
	if lt.Err() != nil {
		t.Error("error is not nil")
		t.FailNow()
	}

	if a != 2 {
		t.Error("core did not launch: a=", a)
		t.FailNow()
	}
}

func TestNewThreadFail(t *testing.T) {
	lt, err := lcore.NewThread(64)
	if err == nil {
		t.FailNow()
	}
	if lt.State() != lcore.ThreadExit {
		t.FailNow()
	}
}

func TestCtxValue(t *testing.T) {
	lt, err := lcore.NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer lt.Exit()
	var wg sync.WaitGroup

	wg.Add(1)
	lt.Execute(func(ctx *lcore.ThreadCtx) error {
		defer wg.Done()
		data := []int{1, 2}
		ctx.Value = data
		return nil
	})
	wg.Wait()

	var data []int
	wg.Add(1)
	lt.Execute(func(ctx *lcore.ThreadCtx) error {
		defer wg.Done()
		data = ctx.Value.([]int)
		return nil
	})

	wg.Wait()
	ok := len(data) == 2 && data[0] == 1 && data[1] == 2
	if !ok {
		t.FailNow()
	}
}

func TestError(t *testing.T) {
	lt, err := lcore.NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer lt.Exit()
	var wg sync.WaitGroup

	someErr := errors.New("some error")
	wg.Add(1)
	lt.Execute(func(ctx *lcore.ThreadCtx) error {
		defer wg.Done()
		return someErr
	})
	wg.Wait()
	if lt.Err() != someErr {
		t.FailNow()
	}
}

func TestState(t *testing.T) {
	lt, err := lcore.NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if lt.State() != lcore.ThreadWait {
		lt.Exit()
		t.FailNow()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	lt.Execute(func(ctx *lcore.ThreadCtx) error {
		defer wg.Done()
		time.Sleep(time.Second)
		return nil
	})

	time.Sleep(100 * time.Millisecond)
	if lt.State() != lcore.ThreadExecute {
		lt.Exit()
		t.FailNow()
	}

	wg.Wait()
	if lt.State() != lcore.ThreadWait {
		lt.Exit()
		t.FailNow()
	}

	lt.Exit()
	if lt.State() != lcore.ThreadExit {
		t.FailNow()
	}
}

func ExampleThread_Execute() {
	lt, err := lcore.NewThread(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lt.Exit()

	var wg sync.WaitGroup
	wg.Add(1)
	errSomeBad := errors.New("some error")
	lt.Execute(func(ctx *lcore.ThreadCtx) error {
		defer wg.Done()
		fmt.Printf("core %d on socket %d\n", ctx.LcoreID(), ctx.SocketID())
		return errSomeBad
	})
	wg.Wait()

	if lt.Err() != errSomeBad {
		fmt.Println("well, it should be")
		return
	}

	fmt.Println("success")
	// Output: success
}
