package eal

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"unsafe"
)

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
