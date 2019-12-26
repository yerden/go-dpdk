/*
Package eal wraps EAL initialization and provides some additional functionality
on top of that. Every CPU's logical core which is setup by EAL runs its own
function which essentially receives functions to execute via Go channel. So you
may run arbitrary Go code in the context of EAL thread.

EAL may be initialized via command line string, parsed command line string or a
set of Options.

Please note that some functions may be called only in EAL thread because of TLS
(Thread Local Storage) dependency.

API is a subject to change. Be aware.
*/
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
	"bufio"
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

type ealConfig struct {
	lcores map[uint]*LcoreCtx
}

var (
	// goEAL is the storage for all EAL lcore threads configuration.
	goEAL = &ealConfig{make(map[uint]*LcoreCtx)}
)

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

// StopLcores sends signal to EAL threads to finish execution of
// go-dpdk lcore function executor.
//
// Warning: it will block until all lcore threads finish execution.
func StopLcores() {
	lcores := Lcores()
	ch := make(chan error, len(lcores))

	for _, id := range lcores {
		ExecOnLcoreAsync(id, ch, func(ctx *LcoreCtx) {
			ctx.done = true
		})
	}

	for range lcores {
		<-ch
	}
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

type lcoreJob struct {
	fn  func(*LcoreCtx)
	ret chan<- error
}

// ExecOnLcoreAsync sends fn to execute on CPU logical core lcoreID,
// i.e. in EAL-owned thread on that lcore.
//
// Possible panic in lcore function will be intercepted and returned
// as an error of type ErrLcorePanic through ret channel specified by
// caller. If lcoreID is invalid, ErrLcoreInvalid error will be
// returned the same way.
//
// The function returns ret. You may specify ret to be nil, in which
// case no error will be reported.
func ExecOnLcoreAsync(lcoreID uint, ret chan error, fn func(*LcoreCtx)) <-chan error {
	if ctx, ok := goEAL.lcores[lcoreID]; ok {
		ctx.ch <- &lcoreJob{fn, ret}
	} else if ret != nil {
		ret <- ErrLcoreInvalid
	}
	return ret
}

// ExecOnLcore sends fn to execute on CPU logical core lcoreID, i.e.
// in EAL-owned thread on that lcore. Then it waits for the execution
// to finish and returns the execution result.
//
// Possible panic in lcore function will be intercepted and returned
// as an error of type ErrLcorePanic. If lcoreID is invalid,
// ErrLcoreInvalid error will be returned.
func ExecOnLcore(lcoreID uint, fn func(*LcoreCtx)) error {
	return <-ExecOnLcoreAsync(lcoreID, make(chan error, 1), fn)
}

// ExecOnMasterAsync is a shortcut for ExecOnLcoreAsync with master
// lcore as a destination.
func ExecOnMasterAsync(ret chan error, fn func(*LcoreCtx)) <-chan error {
	return ExecOnLcoreAsync(GetMasterLcore(), ret, fn)
}

// ExecOnMaster is a shortcut for ExecOnLcore with master lcore as a
// destination.
func ExecOnMaster(fn func(*LcoreCtx)) error {
	return ExecOnLcore(GetMasterLcore(), fn)
}

type lcoresIter struct {
	i  C.uint
	sm C.int
}

func (iter *lcoresIter) next() bool {
	iter.i = C.rte_get_next_lcore(iter.i, iter.sm, 0)
	return iter.i < C.RTE_MAX_LCORE
}

// If skipMaster is 0, master lcore will be included in the result.
// Otherwise, it will miss the output.
func getLcores(skipMaster int) (out []uint) {
	c := &lcoresIter{i: ^C.uint(0), sm: C.int(skipMaster)}
	for c.next() {
		out = append(out, uint(c.i))
	}
	return out
}

// Lcores returns all lcores registered in EAL.
func Lcores() []uint {
	return getLcores(0)
}

// LcoresSlave returns all slave lcores registered in EAL.
// Lcore is slave if it is not master.
func LcoresSlave() []uint {
	return getLcores(1)
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
func ealLaunch(wg *sync.WaitGroup) {
	// init per-lcore contexts
	for _, id := range Lcores() {
		goEAL.lcores[id] = &LcoreCtx{ch: make(chan *lcoreJob, lcoreJobsBuffer)}
	}

	// lcore function
	fn := (*C.lcore_function_t)(C.lcoreFuncListener)

	wg.Add(len(LcoresSlave()))
	// launch every EAL thread lcore function
	// it should be success since we've just called rte_eal_init()
	C.rte_eal_mp_remote_launch(fn, unsafe.Pointer(wg), C.CALL_MASTER)
}

// Cleanup releases EAL-allocated resources, ensuring that no hugepage
// memory is leaked. It is expected that all DPDK applications call
// rte_eal_cleanup() before exiting. Not calling this function could
// result in leaking hugepages, leading to failure during
// initialization of secondary processes.
func Cleanup() error {
	return err(C.rte_eal_cleanup())
}

func parseCmd(input string) ([]string, error) {
	s := bufio.NewScanner(strings.NewReader(input))
	s.Split(common.SplitFunc(common.DefaultSplitter))

	var argv []string
	for s.Scan() {
		argv = append(argv, s.Text())
	}
	return argv, s.Err()
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
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// we should initialize EAL and run EAL threads in a separate
		// goroutine because its thread is going to be acquired by EAL
		// and become master lcore thread
		runtime.LockOSThread()

		// initialize EAL and launch lcoreFuncListener on all slave
		// lcores, then report
		if n, err = ealInit(args); err != nil {
			wg.Done()
		} else {
			ealLaunch(&wg)
		}
	}()
	wg.Wait()
	return
}

// InitCmd initializes EAL as in rte_eal_init. Options are specified
// in a unparsed command line string. This string is parsed and
// Init is then called upon.
func InitCmd(input string) (int, error) {
	argv, err := parseCmd(input)
	if err != nil {
		return 0, err
	}
	return Init(argv)
}

// HasHugePages tells if huge pages are activated.
func HasHugePages() bool {
	return int(C.rte_eal_has_hugepages()) != 0
}

// HasPCI tells whether EAL is using PCI bus. Disabled by –no-pci
// option.
func HasPCI() bool {
	return int(C.rte_eal_has_pci()) != 0
}

// ProcessType returns the current process type.
func ProcessType() int {
	return int(C.rte_eal_process_type())
}

// LcoreCount returns number of CPU logical cores configured by EAL.
func LcoreCount() uint {
	return uint(C.rte_lcore_count())
}

// GetMasterLcore returns CPU logical core id where the master thread
// is executed.
func GetMasterLcore() uint {
	return uint(C.rte_get_master_lcore())
}
