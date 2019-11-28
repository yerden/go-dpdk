package common

import (
	"encoding/binary"
	"unsafe"
)

// CopyFromBytes copies no more than max bytes from src to an area
// pointed to by dst.
func CopyFromBytes(dst unsafe.Pointer, src []byte, max int) int {
	return copy(MakeSlice(dst, max), src)
}

// CopyFromBytes copies no more than max bytes from an area pointed to
// by src to dst.
func CopyToBytes(dst []byte, src unsafe.Pointer, max int) int {
	return copy(dst, MakeSlice(src, max))
}

// PutUint16 stores uint16 value into an area pointed to dst.
func PutUint16(b binary.ByteOrder, dst unsafe.Pointer, d uint16) {
	buf := MakeSlice(dst, 2)
	b.PutUint16(buf, d)
}

// PutUint32 stores uint32 value into an area pointed to dst.
func PutUint32(b binary.ByteOrder, dst unsafe.Pointer, d uint32) {
	buf := MakeSlice(dst, 4)
	b.PutUint32(buf, d)
}

// PutUint64 stores uint64 value into an area pointed to dst.
func PutUint64(b binary.ByteOrder, dst unsafe.Pointer, d uint64) {
	buf := MakeSlice(dst, 8)
	b.PutUint64(buf, d)
}
