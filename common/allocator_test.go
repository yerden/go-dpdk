package common

import (
	// "fmt"
	// "reflect"
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
