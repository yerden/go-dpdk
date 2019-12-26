package common

import (
	// "fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestAllocatorMallocT(t *testing.T) {
	assert := Assert(t, true)
	var p *int
	m := &StdAlloc{}

	obj := MallocT(m, &p)
	*p = 123
	assert(*p == 123)
	assert(unsafe.Pointer(p) == obj)

	m.Free(obj)
}

func TestAllocatorSession(t *testing.T) {
	assert := Assert(t, true)
	var p *int
	m := &StdAlloc{}
	s := NewAllocatorSession(m)
	defer s.Flush()

	obj := MallocT(s, &p)
	*p = 123
	assert(*p == 123)
	assert(unsafe.Pointer(p) == obj)
}

func TestAllocatorSessionCalloc(t *testing.T) {
	assert := Assert(t, true)
	var p *int
	m := &StdAlloc{}
	s := NewAllocatorSession(m)
	defer s.Flush()

	obj := CallocT(s, &p, 2)
	*p = 123
	assert(*p == 123)
	assert(unsafe.Pointer(p) == obj)

	p = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(*p)))
	*p = 456
	assert(*p == 456)
}

func TestAllocatorSessionCalloc2(t *testing.T) {
	assert := Assert(t, true)

	// allocator
	m := &StdAlloc{}
	s := NewAllocatorSession(m)
	defer s.Flush()

	// declare slice
	var p []int
	var ptr *int

	// allocate it
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	sh.Data = uintptr(CallocT(s, &ptr, 2))
	sh.Len = 2
	sh.Cap = 2

	p[0], p[1] = 123, 456
	assert(p[0] == 123)
	assert(p[1] == 456)
	assert(sh.Data == uintptr(unsafe.Pointer(ptr)))
}
