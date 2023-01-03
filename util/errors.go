package util

import (
	"bytes"
	"fmt"
)

// ErrWithMessage wraps Err and buffer to store annotation.
type ErrWithMessage struct {
	Err error
	B   bytes.Buffer
}

func (e *ErrWithMessage) Error() string {
	return fmt.Sprintf("%v: %v", e.B.String(), e.Err)
}

func (e *ErrWithMessage) Unwrap() error {
	return e.Err
}

// ErrWrapf annotates specified err with formatted string.
func ErrWrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	e := &ErrWithMessage{Err: err}
	fmt.Fprintf(&e.B, format, args...)
	return e
}

// ErrWrap annotates specified err with message string.
func ErrWrap(err error, msg string) error {
	return ErrWrapf(err, msg)
}
