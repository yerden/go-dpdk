package eal

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"reflect"
	"unsafe"
)

var (
	hexMap = []byte{
		'0', '1', '2', '3', '4', '5', '6', '7',
		'8', '9', 'a', 'b', 'c', 'd', 'e', 'f',
	}
)

func getHexByte(n int) byte { return hexMap[n] }
func getHexIndex(c byte) int {
	for i, x := range hexMap {
		if x == c {
			return i
		}
	}
	return 0
}

// Set tells if given number (for example, CPU logical core id) is
// present.
type Set interface {
	IsSet(int) bool
}

// SetToHex converts Set into hex string representation. For example,
// [0 1 2 3] converts to "f".
func SetToHex(mask Set, max int) string {
	var out []byte
	for n := 0; n < max; n++ {
		if !mask.IsSet(n) {
			continue
		}

		i, r := n/4, uint(n&3)
		if i >= len(out) {
			out = append(make([]byte, i-len(out)+1), out...)
		}
		i = len(out) - 1 - i
		out[i] = getHexByte(getHexIndex(out[i]) | (1 << r))
	}

	return string(out)
}

// FuncSet is a function which mimics Set interface.
type FuncSet func(int) bool

// IsSet implements Set interface.
func (f FuncSet) IsSet(x int) bool { return f(x) }

// MakeSet attempts to convert various kinds of data into Set
// interface. Currently supported: slices, arrays of integer or
// unsigned integer values, maps with integer/unsigned integer keys.
// If nothing suits the function panics.
func MakeSet(i interface{}) Set {
	if a, ok := i.(Set); ok {
		return a
	}

	intType := reflect.ValueOf(int(0)).Type()

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Map:
		t := v.Type().Key()
		return FuncSet(func(x int) bool {
			key := reflect.ValueOf(x).Convert(t)
			return v.MapIndex(key).IsValid()
		})
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		return FuncSet(func(x int) bool {
			for n := 0; n < v.Len(); n++ {
				elem := v.Index(n).Convert(intType).Int()
				if x == int(elem) {
					return true
				}
			}
			return false
		})
	default:
		return FuncSet(func(x int) bool {
			elem := v.Convert(intType).Int()
			return x == int(elem)
		})
	}
}

func nextArg(argv **C.char) **C.char {
	return (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(argv)) + unsafe.Sizeof(*argv)))
}

// free strings in NULL-terminated array
func freeArgv(argv **C.char) {
	for x := argv; *x != nil; x = nextArg(x) {
		C.free(unsafe.Pointer(*x))
		*x = nil
	}
}

func copyArgv(argv **C.char) **C.char {
	var elem []*C.char
	for x := argv; *x != nil; x = nextArg(x) {
		elem = append(elem, *x)
	}
	elem = append(elem, nil)
	return &elem[0]
}

func makeArgcArgv(argv []string) (C.int, **C.char) {
	argc := C.int(len(argv))
	elem := make([]*C.char, argc, argc+1) // last elem is NULL
	firstElem := (**C.char)(&elem[0])
	for i, arg := range argv {
		elem[i] = C.CString(arg)
	}
	return argc, firstElem
}
