package common

/*
struct void_ptr {
	void *p;
};
*/
import "C"

import (
	"reflect"
	"unsafe"
)

const (
	// Size of the void pointer.
	SizeOfCPtr = C.sizeof_struct_void_ptr
)

// Uintptr returns pointer as an integer number for arbitraty input
// value. The attempt to interpret the value as convertible is made.
func Uintptr(i interface{}) uintptr {
	if p, ok := i.(unsafe.Pointer); ok {
		return uintptr(p)
	}

	switch v := reflect.ValueOf(i); v.Kind() {
	case reflect.Ptr:
		fallthrough
	case reflect.Slice:
		return v.Pointer()
	default:
		panic("hey you")
	}
}
