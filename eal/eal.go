package eal

/*
#include <stdlib.h>

#include <rte_config.h>
#include <rte_eal.h>
#include <rte_errno.h>
#include <rte_lcore.h>

static int eal_init(int argc, char **argv) {
	return rte_eal_init(argc, argv) < 0?rte_errno:0;
}

extern int lcoreFuncListener(void *arg);

static int go_eal_mp_remote_launch(
	lcore_function_t *f,
	uintptr_t arg,
	enum rte_rmt_call_master_t call_master) {
	return rte_eal_mp_remote_launch(f, (void*) arg, call_master);
}
*/
import "C"

import (
	"log"
	"os"
	"runtime"
	"strings"
	"unsafe"
)

// Maximum number of lcores configured during DPDK compile-time.
const (
	MaxLcore = C.RTE_MAX_LCORE
)

// The type of process in a linux, multi-process setup.
const (
	ProcAuto      = C.RTE_PROC_AUTO
	ProcPrimary   = C.RTE_PROC_PRIMARY
	ProcSecondary = C.RTE_PROC_SECONDARY
)

// LcoreFunc is the function prototype to be executed by EAL-owned
// threads.
type LcoreFunc func(*Lcore)

// Lcore is a per-lcore context and is supplied to LcoreFunc as an
// argument.
type Lcore struct {
	// Value is a user-specified context. You may change it as you
	// will and it will persist across LcoreFunc invocations.
	Value interface{}
	// Id is a this thread's CPU logical core id.
	ID uint
	// SocketID is a this thread's CPU socket id.
	SocketID uint

	// channel to receive LcoreFunc to execute.
	ch chan LcoreFunc
}

type ealConfig struct {
	lcores [MaxLcore]*Lcore
}

var (
	// goEAL is the storage for all EAL lcore threads configuration.
	goEAL = &ealConfig{}
)

// to run as lcore_function_t
//export lcoreFuncListener
func lcoreFuncListener(arg unsafe.Pointer) C.int {
	eal := (*ealConfig)(arg)
	lc := eal.lcores[LcoreID()]
	log.Println("started on lcore", LcoreID())

	for fn := range lc.ch {
		if fn == nil {
			break
		}
		lc.refresh()
		fn(lc)
	}
	return 0
}

func (lc *Lcore) refresh() {
	lc.ID = LcoreID()
	lc.SocketID = uint(C.rte_lcore_to_socket_id(C.uint(lc.ID)))
}

// ExecuteOnLcore sends fn to execute on CPU logical core lcoreID, i.e
// in EAL-owned thread on that lcore.
func ExecuteOnLcore(lcoreID uint, fn LcoreFunc) {
	if lc := goEAL.lcores[lcoreID]; lc != nil {
		lc.ch <- fn
	} else {
		panic("")
	}
}

// ExecuteOnMaster is a shortcut for ExecuteOnLcore with master lcore
// as a destination.
func ExecuteOnMaster(fn LcoreFunc) {
	ExecuteOnLcore(GetMasterLcore(), fn)
}

// ForeachLcore iterates through all CPU logical cores initialized by
// EAL. If skipMaster is true the iteration will skip master lcore.
func ForeachLcore(skipMaster bool, f func(lcoreID uint)) {
	i := ^C.uint(0)
	sm := C.int(0)

	if skipMaster {
		sm = 1
	}

	for {
		i = C.rte_get_next_lcore(i, sm, 0)
		if i >= C.RTE_MAX_LCORE {
			break
		}
		f(uint(i))
	}
}

// InitWithArgs initializes EAL as in rte_eal_init. Options are
// specified in a parsed command line string.
//
// This function initialized EAL and waits for executable functions on
// each of EAL-owned threads.
func InitWithArgs(argv []string) error {
	argv = append([]string{os.Args[0]}, argv...)
	argc := C.int(len(argv))
	cArgv := make([]*C.char, argc+1)
	for i, arg := range argv {
		cArgv[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(cArgv[i]))
	}

	ch := make(chan error, 1)
	go func() {
		// we should initialize EAL and run EAL threads in a separate
		// goroutine because its thread is going to be acquired by EAL
		runtime.LockOSThread()

		// initialize EAL
		err := errno(C.eal_init(argc, (**C.char)(&cArgv[0])))
		if err != nil {
			ch <- err
			return
		}

		// init per-lcore contexts
		ForeachLcore(false, func(lcoreID uint) {
			goEAL.lcores[lcoreID] = &Lcore{ch: make(chan LcoreFunc, 1)}
		})

		fn := (*C.lcore_function_t)(C.lcoreFuncListener)
		// nasty trick, but justified
		// since EAL struct is allocated globally, it won't be GC-ed,
		// so we may 'safely' cast the pointer to C.uintptr_t
		arg := C.uintptr_t(uintptr(unsafe.Pointer(goEAL)))

		// report that we're ok, otherwise we will panic
		ch <- nil
		defer log.Println("master lcore exited")
		// launch every EAL thread lcore function
		// it should be success since we've just called rte_eal_init()
		err = errno(C.go_eal_mp_remote_launch(fn, arg, C.SKIP_MASTER))
		if err != nil {
			ch <- err
			return
		}

		// run on master lcore
		lcoreFuncListener(unsafe.Pointer(goEAL))
	}()

	return <-ch
}

// Init initializes EAL as in rte_eal_init. Options are
// specified in a unparsed command line string. This string is parsed
// and InitWithArgs is then called upon.
func Init(argv string) error {
	return InitWithArgs(strings.Split(argv, " "))
}

// InitWithOpts initializes EAL as in rte_eal_init. Options are
// specified in array of Option-s. These options are then used to
// construct argv array and InitWithArgs is then called upon.
func InitWithOpts(opts ...Option) error {
	p := ealOptions{}
	for _, opt := range opts {
		opt.f(&p)
	}

	return InitWithArgs(p.argv())
}

// HasHugePages tells if huge pages are activated.
func HasHugePages() bool {
	return int(C.rte_eal_has_hugepages()) != 0
}

// ProcessType returns the current process type.
func ProcessType() int {
	return int(C.rte_eal_process_type())
}

// LcoreID returns CPU logical core id. This function must be called
// only in EAL thread.
func LcoreID() uint {
	return uint(C.rte_lcore_id())
}

// SocketID returns NUMA socket where the current thread resides. This
// function must be called only in EAL thread.
func SocketID() uint {
	return uint(C.rte_socket_id())
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
