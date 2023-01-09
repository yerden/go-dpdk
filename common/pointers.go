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

// CStruct is a GO structure representation of a C array.
// CStruct has a certain length, don't try to extend it.
type CStruct struct {
	// Ptr is a pointer to the beginning of the C array.
	Ptr unsafe.Pointer
	// Len is the size of memory area pointed by Ptr.
	Len int
}

// Init initializes CStruct instance with specified pointer ptr and
// length len.
func (cs *CStruct) Init(ptr unsafe.Pointer, len int) {
	cs.Ptr = ptr
	cs.Len = len
}

// Bytes converts C array into slice of bytes backed by this array.
// Returned slice cannot be extended or it will be reallocated into Go
// memory.
func (cs *CStruct) Bytes() (dst []byte) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&dst))
	sh.Data = uintptr(cs.Ptr)
	sh.Len = cs.Len
	sh.Cap = cs.Len
	return
}

// Memset initializes memory pointed by p and with length n.
func Memset(p unsafe.Pointer, init byte, n uintptr) {
	b := unsafe.Slice((*byte)(p), n)
	for i := range b {
		b[i] = init
	}
}
