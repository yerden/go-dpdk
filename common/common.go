/*
Package common contains basic routines needed by other modules in
go-dpdk package.
*/
package common

/*
#include <rte_memory.h>
#include <rte_errno.h>
static int rteErrno() {
	return rte_errno;
}
*/
import "C"

import (
	"errors"
	"reflect"
	"sync"
	"syscall"
)

// SocketIDAny represents selection for any NUMA socket.
const (
	SocketIDAny = int(C.SOCKET_ID_ANY)
)

// Custom RTE induced errors.
var (
	ErrNoConfig  = errors.New("Missing rte_config")
	ErrSecondary = errors.New("Operation not allowed in secondary processes")
)

func errno(n int64) error {
	if n == 0 {
		return nil
	} else if n < 0 {
		n = -n
	}

	if n == int64(C.E_RTE_NO_CONFIG) {
		return ErrNoConfig
	}

	if n == int64(C.E_RTE_SECONDARY) {
		return ErrSecondary
	}

	return syscall.Errno(int(n))
}

// RteErrno returns rte_errno variable.
func RteErrno() error {
	return errno(int64(C.rteErrno()))
}

// IntOrErr returns error as in Errno in case n is negative.
// Otherwise, the value itself with nil error will be returned.
//
// If n is nil, then n = RteErrno()
// if n is not nil and not a signed integer, function panics.
func IntOrErr(n interface{}) (int, error) {
	x := reflect.ValueOf(n).Int()
	if x >= 0 {
		return int(x), nil
	}
	return 0, errno(x)
}

// IntToErr converts n into an 'errno' error. If n is not a signed
// integer it will panic.
func IntToErr(n interface{}) error {
	x := reflect.ValueOf(n).Int()
	return errno(x)
}

// DoOnce decorates fn in a way that it will effectively run only once
// returning the resulting error value in this and all subsequent
// calls. Useful in unit testing when initialization must be performed
// only once.
func DoOnce(fn func() error) func() error {
	var once sync.Once
	var err error
	return func() error {
		once.Do(func() { err = fn() })
		return err
	}
}

func Err(n ...interface{}) error {
	if len(n) == 0 {
		return RteErrno()
	}

	return IntToErr(n[0])
}