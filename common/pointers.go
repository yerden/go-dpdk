package common

import (
	"reflect"
	"unsafe"
)

// MakeSlice returns byte slice specified by pointer and of len max.
func MakeSlice(buf unsafe.Pointer, max int) []byte {
	var dst []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&dst))
	sh.Data = uintptr(buf)
	sh.Len = max
	sh.Cap = max
	return dst
}

// CArray is a GO structure representation of a C array.
// Ptr is a pointer to the beginning of the C array and Len is its maximum length.
// CArray has a certain length, don't try to extend it.
type CArray struct {
	Ptr unsafe.Pointer
	Len int
}
