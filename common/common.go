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
	"syscall"
)

// SocketIDAny represents selection for any NUMA socket.
const (
	SocketIDAny = int(C.SOCKET_ID_ANY)
)

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
func RteErrno() int {
	return int(C.rteErrno())
}

func tryErrno(n interface{}) (x int64) {
	if n == nil {
		x = int64(C.rteErrno())
	} else {
		x = reflect.ValueOf(n).Int()
	}
	return x
}

// IntOrErr returns error as in Errno in case n is negative.
// Otherwise, the value itself with nil error will be returned.
func IntOrErr(n interface{}) (int, error) {
	x := tryErrno(n)
	if x >= 0 {
		return int(x), nil
	}
	return 0, errno(x)
}

// Errno converts return value of C function into meaningful error.
func Errno(n interface{}) error {
	return errno(tryErrno(n))
}
