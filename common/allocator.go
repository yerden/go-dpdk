package common

import "C"

import (
	"reflect"
	"unsafe"
)

// Allocator provides allocating and freeing of objects. It should be
// used with Cgo to withstand the rule of not allowing Go pointers
// inside a Go pointer. The allocator allows to defer freeing of
// objects so instead of freeing objects individually you may delete
// them by Flush at once and abandon allocator instance.
type Allocator interface {
	// Malloc allocates memory of length size.
	Malloc(size uintptr) unsafe.Pointer
	// Free releases previously allocated memory pointed to by p.
	Free(p unsafe.Pointer)
	// Realloc allocates memory of length size.
	Realloc(p unsafe.Pointer, size uintptr) unsafe.Pointer
}

// AllocatorSession wraps Allocator and storage for allocated
// pointers. Useful to perform allocations and free them with one
// call.
type AllocatorSession struct {
	mem  Allocator
	ptrs map[unsafe.Pointer]struct{}
}

// NewAllocatorSession creates new AllocatorSession.
func NewAllocatorSession(mem Allocator) *AllocatorSession {
	return &AllocatorSession{mem, make(map[unsafe.Pointer]struct{})}
}

// Malloc implements Allocator.
func (s *AllocatorSession) Malloc(size uintptr) unsafe.Pointer {
	p := s.mem.Malloc(size)
	s.ptrs[p] = struct{}{}
	return p
}

// Free implements Allocator.
func (s *AllocatorSession) Free(p unsafe.Pointer) {
	s.mem.Free(p)
	delete(s.ptrs, p)
}

// Realloc implements Allocator.
func (s *AllocatorSession) Realloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	p1 := s.mem.Realloc(p, size)
	if p != p1 {
		s.ptrs[p1] = struct{}{}
		delete(s.ptrs, p)
	}
	return p1
}

// Flush releases all previously allocated memory in this session.
func (s *AllocatorSession) Flush() {
	for p := range s.ptrs {
		s.mem.Free(p)
	}
	s.ptrs = make(map[unsafe.Pointer]struct{})
}

// compile-time assertion
var _ Allocator = (*AllocatorSession)(nil)

// MallocT allocates an object by its type. The type and its size is
// derived from ptr which is a pointer to pointer of required type
// where new object will be stored. For example:
//   var x *int
//   a := NewAlloc()
//   defer a.Flush()
//   MallocT(a, &x)
//   /* x is now an allocated pointer */
func MallocT(a Allocator, ptr interface{}) unsafe.Pointer {
	// get the type of value to allocate;
	// v should be the pointer to pointer,
	// hence twice Elem
	v := reflect.ValueOf(ptr)
	t := v.Type().Elem().Elem()
	p := a.Malloc(t.Size())
	reflect.Indirect(v).Set(reflect.NewAt(t, p))
	return p
}

// CBytes creates a copy of byte slice with given Allocator. It's
// analogous to C.CBytes.
func CBytes(a Allocator, b []byte) unsafe.Pointer {
	p := a.Malloc(uintptr(len(b)))
	copy(MakeSlice(p, len(b)), b)
	return p
}

// CString a copy of a string with given Allocator. It's analogous to
// C.CString.
func CString(a Allocator, s string) *C.char {
	p := a.Malloc(uintptr(len(s) + 1))
	dst := MakeSlice(p, len(s)+1)
	copy(dst, s)
	dst[len(s)] = 0
	return (*C.char)(p)
}
