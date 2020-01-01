package common

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"
)

// Assert allows to perform one-lined tests and, optionally, print
// some diagnostic info if the test failed.
//
// If fail is true, test failure will cause panic and cease test
// execution.
func Assert(t testing.TB, fail bool) func(bool, ...interface{}) {
	return func(expected bool, v ...interface{}) {
		if !expected {
			t.Helper()
			if t.Error(v...); fail {
				t.FailNow()
			}
		}
	}
}

// FprintStackFrames prints calling stack of the error into specified
// writer. Program counters are specified in pc.
func FprintStackFrames(w io.Writer, pc []uintptr) {
	frames := runtime.CallersFrames(pc)
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
