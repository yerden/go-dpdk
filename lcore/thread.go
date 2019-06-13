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
	"sync/atomic"

	"golang.org/x/sys/unix"
)

// Lcore thread state.
const (
	// Lcore thread is waiting for function to execute.
	ThreadWait int32 = iota
	// Lcore thread is executing a function.
	ThreadExecute
	// Lcore thread was called Exit() upon.
	ThreadExit
)

// ThreadFunc is a function prototype lcore thread can run.
type ThreadFunc func(*ThreadCtx) error

type job struct {
	fn ThreadFunc
}

// Thread controls execution of functions in a thread which
// belongs to specific logical CPU core.
type Thread struct {
	jobsCh chan job
	state  int32
	err    error
	wg     sync.WaitGroup
}

// ThreadCtx is supplied in ThreadFunc as a storage for lcore thread
// related data.
type ThreadCtx struct {
	// Value is a user-supplied data field which is persistent across
	// multiple lcore functions executed on the same lcore thread.
	Value interface{}

	lcoreID  uint
	numaNode int
}

// LcoreID returns id of logical CPU core which lcore thread is tied
// to.
func (ctx *ThreadCtx) LcoreID() uint {
	return ctx.lcoreID
}

// SocketID returns id of CPU socket where lcore thread and logical
// core resides.
func (ctx *ThreadCtx) SocketID() int {
	return ctx.numaNode
}

// NewThread returns new Thread and possible error which may
// arise if the set-affinity system call is unsuccessful.
//
// id is the index of the desired logical CPU core.
func NewThread(id uint) (*Thread, error) {
	t := &Thread{
		jobsCh: make(chan job),
		state:  ThreadWait,
		err:    nil,
	}

	ch := make(chan error, 1)
	t.wg.Add(1)
	go func() {
		// we don't want this thread to become usable by Go and we'd
		// like it to be destroyed. So, we don't call matching
		// UnlockOSThread().
		runtime.LockOSThread()
		defer t.wg.Done()
		defer close(ch)

		// Set affinity to specified CPU for current thread.
		var s unix.CPUSet
		s.Set(int(id))
		err := unix.SchedSetaffinity(0, &s)
		if err != nil {
			// set state and return error
			atomic.StoreInt32(&t.state, ThreadExit)
			ch <- err
			return
		}

		defer atomic.StoreInt32(&t.state, ThreadExit)
		ch <- err

		ctx := &ThreadCtx{
			lcoreID:  id,
			numaNode: NumaNode(id),
		}

		for j := range t.jobsCh {
			if j.fn == nil {
				return
			}
			atomic.StoreInt32(&t.state, ThreadExecute)
			t.err = j.fn(ctx)
			atomic.StoreInt32(&t.state, ThreadWait)
		}
	}()

	return t, <-ch
}

// Err returns error returned by lcore function after ending
// execution. It is safe to call only when thread has finished
// executing a function.
func (t *Thread) Err() error {
	return t.err
}

// Exit sends a signal to stop all activity on the thread. After
// that, it waits for the current lcore function to finish execution.
// After that, only Err() and State() may be called to check upon the
// error and the state of the thread. Other calls are prohibited.
func (t *Thread) Exit() {
	t.jobsCh <- job{nil}
	t.wg.Wait()
}

// State returns current state of the thread, see Thread*
// constants.
func (t *Thread) State() int32 {
	return atomic.LoadInt32(&t.state)
}

// Execute sends new lcore function to execute. This function will
// block until previous lcore function finishes.
func (t *Thread) Execute(fn ThreadFunc) {
	if fn != nil {
		t.jobsCh <- job{fn}
	}
}
