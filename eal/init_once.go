package eal

import (
	"fmt"
	"sync"
)

// SafeEALArgs returns safe parameters to be used for testing
// purposes. Specify cmdname as the binary name in argv[0] and number
// of lcores. All lcores will be assigned to core 0.
func SafeEALArgs(cmdname string, lcores int) []string {
	return []string{
		cmdname,
		"--lcores", fmt.Sprintf("(0-%d)@0", lcores-1),
		"--vdev", "net_null0",
		"-m", "128",
		"--no-huge",
		"--no-pci",
		"--main-lcore", "0",
	}
}

// multiple calls guard
var ealOnce struct {
	sync.Mutex
	already bool
}

// InitOnce calls Init guarded with global lock.
//
// If Init returns error it panics. If Init was already called it
// simply returns.
//
// It's mostly intended to use in tests.
func InitOnce(args []string) {
	x := &ealOnce

	x.Lock()
	defer x.Unlock()

	if !x.already {
		if _, err := Init(args); err != nil {
			panic(err)
		}
	}

	x.already = true
}

// InitOnceSafe calls Init guarded with global lock on arguments
// returned by SafeEALArgs.
//
// If Init returns error it panics. If Init was already called it
// simply returns.
//
// It's mostly intended to use in tests.
func InitOnceSafe(cmdname string, lcores int) {
	InitOnce(SafeEALArgs(cmdname, lcores))
}
