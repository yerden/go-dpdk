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

// Thread is a channel which jobs can be sent to.
type Thread chan<- func()

// NewLockedThread creates new Thread with user specified channel.
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

// Exec sends new job to the Thread. If wait is true this function
// blocks until job finishes execution.
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

// Gettid returns Thread id.
func (t Thread) Gettid() (tid int) {
	t.Exec(true, func() {
		tid = unix.Gettid()
	})
	return
}

// SetAffinity pins the thread to specified CPU core id.
func (t Thread) SetAffinity(id uint) error {
	var s unix.CPUSet
	s.Set(int(id))
	return unix.SchedSetaffinity(t.Gettid(), &s)
}

// GetAffinity retrieves CPU affinity of the Thread.
func (t Thread) GetAffinity() (s unix.CPUSet, err error) {
	return s, unix.SchedGetaffinity(t.Gettid(), &s)
}

// Close sends a signal to Thread to finish.
func (t Thread) Close() {
	close(t)
}
