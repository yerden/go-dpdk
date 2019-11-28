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
