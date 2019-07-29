/*
Package lcore allows to control execution of user-supplied functions
on specified logical CPU core.

This may have some advantages such as: reduce context switches, allows
to use non-preemptible algorithms etc.

Please note that thread function is entirely specified by user. It is
up to user to define how this function would exit.
*/
package lcore

import (
	// "fmt"
	"runtime"
	"sync"

	"golang.org/x/sys/unix"
)

type Thread chan<- func()

func NewLockedThread(ch chan func()) Thread {
	go func() {
		runtime.LockOSThread()
		for f := range ch {
			f()
		}
		runtime.UnlockOSThread()
	}()

	return ch
}

func (t Thread) Exec(wait bool, fn func()) {
	if !wait {
		t <- fn
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	t <- func() {
		fn()
		wg.Done()
	}
	wg.Wait()
}

func (t Thread) Gettid() (tid int) {
	t.Exec(true, func() {
		tid = unix.Gettid()
	})
	return
}

func (t Thread) SetAffinity(id uint) error {
	var s unix.CPUSet
	s.Set(int(id))
	return unix.SchedSetaffinity(t.Gettid(), &s)
}

func (t Thread) GetAffinity() (s unix.CPUSet, err error) {
	return s, unix.SchedGetaffinity(t.Gettid(), &s)
}

func (t Thread) Close() {
	close(t)
}
