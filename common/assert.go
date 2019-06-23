package common

import (
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
