package launch

import (
	"errors"
	"testing"
	"time"
)

func TestNewThread(t *testing.T) {
	lt, err := NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer lt.Stop()

	a := 1
	ok := lt.Launch(func(ctx *ThreadCtx) error {
		a = 2
		return nil
	})

	if !ok {
		t.Error("unable to launch")
		t.FailNow()
	}

	lt.Wait()
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
	lt, err := NewThread(64)
	if err == nil {
		t.FailNow()
	}
	defer lt.Stop() // should be no panic
	lt.Stop()
	if lt.State() != ThreadStop {
		t.FailNow()
	}
}

func TestCtxValue(t *testing.T) {
	lt, err := NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer lt.Stop()

	ok := lt.Launch(func(ctx *ThreadCtx) error {
		data := []int{1, 2}
		ctx.Value = data
		return nil
	})
	lt.Wait()

	var data []int
	ok = lt.Launch(func(ctx *ThreadCtx) error {
		data = ctx.Value.([]int)
		return nil
	})

	lt.Wait()
	ok = len(data) == 2 && data[0] == 1 && data[1] == 2
	if !ok {
		t.FailNow()
	}
}

func TestError(t *testing.T) {
	lt, err := NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer lt.Stop()

	someErr := errors.New("some error")
	ok := lt.Launch(func(ctx *ThreadCtx) error {
		return someErr
	})
	lt.Wait()
	if !ok || lt.Err() != someErr {
		t.FailNow()
	}
}

func TestState(t *testing.T) {
	lt, err := NewThread(0)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if lt.State() != ThreadWait {
		lt.Stop()
		t.FailNow()
	}

	ok := lt.Launch(func(ctx *ThreadCtx) error {
		time.Sleep(time.Second)
		return nil
	})

	time.Sleep(100 * time.Millisecond)
	if !ok || lt.State() != ThreadExecute {
		lt.Stop()
		t.FailNow()
	}

	lt.Wait()
	if lt.State() != ThreadWait {
		lt.Stop()
		t.FailNow()
	}

	lt.Stop()
	if lt.State() != ThreadStop {
		t.FailNow()
	}
}
