package flow

/*
#include <stdint.h>
#include <rte_config.h>
#include <rte_flow.h>
*/
import "C"
import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

func off(p unsafe.Pointer, d uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + d)
}

func beU16(n uint16, p unsafe.Pointer) {
	var d []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	sh.Data = uintptr(p)
	sh.Len = 2
	sh.Cap = sh.Len
	binary.BigEndian.PutUint16(d, n)
}

func beU32(n uint32, p unsafe.Pointer) {
	var d []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&d))
	sh.Data = uintptr(p)
	sh.Len = 4
	sh.Cap = sh.Len
	binary.BigEndian.PutUint32(d, n)
}
