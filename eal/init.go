package eal

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_eal.h>
#include <rte_errno.h>
#include <rte_lcore.h>

extern int lcoreFuncListener(void *arg);
*/
import "C"

import (
	"fmt"
	"io"
	"log"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

// Maximum number of lcores configured during DPDK compile-time.
const (
	MaxLcore = C.RTE_MAX_LCORE
)

const (
	lcoreJobsBuffer = 32
)

// The type of process in a linux, multi-process setup.
const (
	ProcAuto      int = C.RTE_PROC_AUTO
	ProcPrimary       = C.RTE_PROC_PRIMARY
	ProcSecondary     = C.RTE_PROC_SECONDARY
)

var (
	// goEAL is the storage for all EAL lcore threads configuration.
	goEAL = &ealConfig{make(map[uint]*LcoreCtx)}
)

type lcoreJob struct {
	fn  func(*LcoreCtx)
	ret chan<- error
}

// LcoreCtx is a per-lcore context and is supplied to function running to
// particular lcore.
type LcoreCtx struct {
	// Value is a user-specified context. You may change it as you
	// will and it will persist across function invocations on
	// particular lcore.
	Value interface{}

	// channel to receive functions to execute.
	ch chan *lcoreJob

	// signal to kill current thread
	done bool
}

type ealConfig struct {
	lcores map[uint]*LcoreCtx
}

func err(n ...interface{}) error {
	if len(n) == 0 {
		return common.RteErrno()
	}

	return common.IntToErr(n[0])
}

// LcoreID returns CPU logical core id. This function must be called
// only in EAL thread.
func (ctx *LcoreCtx) LcoreID() uint {
	return uint(C.rte_lcore_id())
}

// SocketID returns NUMA socket where the current thread resides. This
// function must be called only in EAL thread.
func (ctx *LcoreCtx) SocketID() uint {
	return uint(C.rte_socket_id())
}

// String implements fmt.Stringer.
func (ctx *LcoreCtx) String() string {
	return fmt.Sprintf("lcore=%d socket=%d", ctx.LcoreID(), ctx.SocketID())
}

// LcoreToSocket return socket id for given lcore LcoreID.
func LcoreToSocket(id uint) uint {
	return uint(C.rte_lcore_to_socket_id(C.uint(id)))
}

// ErrLcorePanic is an error returned by ExecOnLcore in case lcore
// function panics.
type ErrLcorePanic struct {
	Pc      []uintptr
	LcoreID uint
	errStr  string
}

// Error implements error interface.
func (e *ErrLcorePanic) Error() string {
	return e.errStr
}

// PrintStack prints calling stack of the error into specified writer.
func (e *ErrLcorePanic) PrintStack(w io.Writer) {
	frames := runtime.CallersFrames(e.Pc)
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		// skipping everything from runtime package.
		if strings.HasPrefix(frame.Function, "runtime.") {
			continue
		}
		fmt.Fprintf(w, "... at %s:%d, %s\n", frame.File, frame.Line,
			frame.Function)
	}
}

// ErrLcoreInvalid is returned by ExecOnLcore in case the desired
// lcore ID is invalid.
var ErrLcoreInvalid = fmt.Errorf("Invalid logical core")

// panicCatcher launches lcore function and returns possible panic as
// an error.
func panicCatcher(fn func(*LcoreCtx), ctx *LcoreCtx) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		errStr := fmt.Sprintf("panic on lcore %d: %v", ctx.LcoreID(), r)
		pc := make([]uintptr, 64)
		// this function is called from runtime package, so to
		// unwind the stack we may skip (1) runtime.Callers
		// function, (2) this caller function
		n := runtime.Callers(2, pc)
		err = &ErrLcorePanic{pc[:n], ctx.LcoreID(), errStr}
	}()
	fn(ctx)
	return err
}

// to run as lcore_function_t
//export lcoreFuncListener
func lcoreFuncListener(arg unsafe.Pointer) C.int {
	id := uint(C.rte_lcore_id())
	ctx := goEAL.lcores[id]
	log.Printf("lcore %d started", id)
	defer log.Printf("lcore %d exited", id)

	// wait group to signal successful launch
	// lcore is running
	wg := (*sync.WaitGroup)(arg)
	wg.Done()

	// run loop
	for job := range ctx.ch {
		err := panicCatcher(job.fn, ctx)
		if job.ret != nil {
			job.ret <- err
		}
		if ctx.done {
			break
		}
	}
	return 0
}

// call rte_eal_init and report its return value and rte_errno as an
// error. Should be run in master lcore thread only
func ealInit(args []string) (int, error) {
	mem := common.NewAllocatorSession(&common.StdAlloc{})
	defer mem.Flush()

	argc := C.int(len(args))
	argv := make([]*C.char, argc+1)
	for i := range args {
		argv[i] = (*C.char)(common.CString(mem, args[i]))
	}

	// initialize EAL
	n := int(C.rte_eal_init(argc, &argv[0]))
	if n < 0 {
		return n, err()
	}
	return n, nil
}

// launch lcoreFuncListener on all slave lcores
// should be run in master lcore thread only
//
// each lcore must do wg.Done() upon successful launch.
func ealLaunch(wg *sync.WaitGroup) {
	// init per-lcore contexts
	for _, id := range Lcores() {
		goEAL.lcores[id] = &LcoreCtx{ch: make(chan *lcoreJob, lcoreJobsBuffer)}
	}

	// lcore function
	fn := (*C.lcore_function_t)(C.lcoreFuncListener)

	// launch every EAL thread lcore function
	// it should be success since we've just called rte_eal_init()
	C.rte_eal_mp_remote_launch(fn, unsafe.Pointer(wg), C.CALL_MASTER)
}

// Init initializes EAL as in rte_eal_init. Options are specified in a
// parsed command line string.
//
// This function initialized EAL and waits for executable functions on
// each of EAL-owned threads.
//
// Returns number of parsed args and an error.
func Init(args []string) (n int, err error) {
	log.Println("EAL parameters:", args)

	// This WaitGroup is used to notify caller that lcoreFuncListener
	// is successfully executed on every EAL lcore and thus must be
	// released (i.e. unlock Wait()) upon finishing master lcore setup
	// or returning EAL init error.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// we should initialize EAL and run EAL threads in a separate
		// goroutine because its thread is going to be acquired by EAL
		// and become master lcore thread
		runtime.LockOSThread()

		// initialize EAL
		if n, err = ealInit(args); err != nil {
			wg.Done()
			return
		}

		// we're about to launch lcore functions so we add slave
		// lcores to WaitGroup
		wg.Add(int(LcoreCount() - 1))

		// ealLaunch runs lcoreFuncListener on all slave lcores and
		// master lcore. It will block until lcoreFuncListener stops
		// on master lcore, see StopLcores.
		ealLaunch(&wg)
	}()
	wg.Wait()
	return
}

// Cleanup releases EAL-allocated resources, ensuring that no hugepage
// memory is leaked. It is expected that all DPDK applications call
// rte_eal_cleanup() before exiting. Not calling this function could
// result in leaking hugepages, leading to failure during
// initialization of secondary processes.
func Cleanup() error {
	return err(C.rte_eal_cleanup())
}
