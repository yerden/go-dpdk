package common

/*
#include <rte_errno.h>
*/
import "C"

import (
	"errors"
	"syscall"
)

var (
	ErrNoConfig  = errors.New("Missing rte_config")
	ErrSecondary = errors.New("Operation not allowed in secondary processes")
)

// IntOrErr returns error as in Errno in case n is negative.
// Otherwise, the value itself with nil error will be returned.
func IntOrErr(n int) (int, error) {
	if n >= 0 {
		return n, nil
	}
	return 0, Errno(n)
}

// Errno converts return value of C function into meaningful error.
func Errno(n int) error {
	if n == 0 {
		return nil
	} else if n < 0 {
		n = -n
	}

	if n == int(C.E_RTE_NO_CONFIG) {
		return ErrNoConfig
	}

	if n == int(C.E_RTE_SECONDARY) {
		return ErrSecondary
	}

	return syscall.Errno(n)
}
