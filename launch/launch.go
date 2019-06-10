/*
Launch package allows to control execution of user-supplied functions
on specified logical CPU core.

This may have some advantages such as: reduce context switches, allows
to use non-preemptible algorithms etc.
*/
package launch

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
	// Lcore thread was called Stop() upon.
	ThreadStop
)

// ThreadFunc is a function prototype lcore thread can run.
type ThreadFunc func(*ThreadCtx) error

type job struct {
	fn ThreadFunc
}

// Thread controls execution of functions in a thread which
// belongs to specific logical CPU core.
type Thread struct {
	jobsCh   chan job
	state    int32
	ctrl     int32
	err      error
	wgJobs   sync.WaitGroup
	wgGlobal sync.WaitGroup
}

// ThreadCtx is supplied in ThreadFunc as a storage for lcore thread
// related data.
type ThreadCtx struct {
	// Value is a user-supplied data field which is persistent across
	// multiple lcore functions executed on the same lcore thread.
	Value    interface{}
	lcoreId  uint
	numaNode int
	ctrl     *int32
}

// LcoreID returns id of logical CPU core which lcore thread is tied
// to.
func (ctx *ThreadCtx) LcoreID() uint {
	return ctx.lcoreId
}

// SocketID returns id of CPU socket where lcore thread and logical
// core resides.
func (ctx *ThreadCtx) SocketID() int {
	return ctx.numaNode
}

// IsStop returns true if the lcore thread was requested to stop
// execution.
func (ctx *ThreadCtx) IsStop() bool {
	return ctx.ctrl != nil && atomic.LoadInt32(ctx.ctrl) == 1
}

// NewThread returns new Thread and possible error which may
// arise if the set-affinity system call is unsuccessful.
//
// id is the index of the desired logical CPU core.
func NewThread(id uint) (*Thread, error) {
	t := &Thread{
		jobsCh: make(chan job),
		state:  ThreadWait,
		ctrl:   0,
		err:    nil,
	}

	ch := make(chan error, 1)
	t.wgGlobal.Add(1)
	go func() {
		// we don't call matching UnlockOSThread() since we don't rely
		// on this particular thread state.
		runtime.LockOSThread()
		defer t.wgGlobal.Done()
		defer close(ch)

		wg := &t.wgJobs

		var s unix.CPUSet
		s.Set(int(id))
		err := unix.SchedSetaffinity(0, &s)
		if err != nil {
			atomic.StoreInt32(&t.state, ThreadStop)
			ch <- err
			return
		}

		defer atomic.StoreInt32(&t.state, ThreadStop)
		ch <- err

		ctx := &ThreadCtx{
			lcoreId:  id,
			numaNode: -1,
			ctrl:     &t.ctrl,
		}
		for j := range t.jobsCh {
			if j.fn == nil {
				return
			}
			atomic.StoreInt32(&t.state, ThreadExecute)
			if t.err = j.fn(ctx); ctx.IsStop() {
				continue
			}
			atomic.StoreInt32(&t.state, ThreadWait)
			wg.Done()
		}
	}()

	return t, <-ch
}

// Err returns error returned by lcore function after ending
// execution.
func (t *Thread) Err() error {
	return t.err
}

// Stop sends a signal to stop all activity on logical thread. After
// that, it waits for the current lcore function to finish execution.
func (t *Thread) Stop() {
	if atomic.LoadInt32(&t.state) != ThreadStop {
		atomic.StoreInt32(&t.ctrl, 1)
		t.jobsCh <- job{nil}
	}
	t.wgGlobal.Wait()
}

// State returns current state of logical thread, see constants
// Thread*.
func (t *Thread) State() int32 {
	return atomic.LoadInt32(&t.state)
}

// Launch sends new lcore function to execute. It returns true if the
// job was enqueued or false if logical thread is busy executing
// another lcore function.
func (t *Thread) Launch(fn ThreadFunc) bool {
	wg := &t.wgJobs
	wg.Add(1)
	select {
	case t.jobsCh <- job{fn}:
		// successfully enqueued
		// therefore we should decrement wg on remote launch
		return true
	default:
		wg.Done()
		return false
	}
}

// Wait blocks calling goroutine until current lcore function
// execution on logical thread is finished.
func (t *Thread) Wait() {
	wg := &t.wgJobs
	defer wg.Wait()
}

// type Grid struct {
// threads map[uint]*Thread
// }

// func NewGrid() *Grid {
// return &Grid{make(map[uint]*Thread)}
// }

// func (g *Grid) Set(lcore uint) {
// g.threads[lcore] = NewThread(lcore)
// }
