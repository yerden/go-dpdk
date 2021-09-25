package util

import (
	"bytes"
	"hash/crc32"
	"testing"
)

func makeGoHash32Func(complement bool) func([]byte, uint32) uint32 {
	tab := crc32.MakeTable(crc32.Castagnoli)
	if complement {
		return func(data []byte, acc uint32) uint32 {
			return ^crc32.Update(^acc, tab, data)
		}
	}

	{
		return func(data []byte, acc uint32) uint32 {
			return crc32.Update(acc, tab, data)
		}
	}
}

func onesComplement(fn func([]byte, uint32) uint32) func([]byte, uint32) uint32 {
	return func(data []byte, acc uint32) uint32 {
		return ^fn(data, ^acc)
	}
}

func BenchmarkComplementFalse(b *testing.B) {
	h := &Hash32{
		OnesComplement: false,
		Accum:          onesComplement(makeGoHash32Func(false)),
	}

	data := bytes.Repeat([]byte{0xab}, 31)

	for i := 0; i < b.N; i++ {
		h.Write(data)
	}
}

func BenchmarkComplementTrue(b *testing.B) {
	h := &Hash32{
		OnesComplement: true,
		Accum:          makeGoHash32Func(false),
	}

	data := bytes.Repeat([]byte{0xab}, 31)

	for i := 0; i < b.N; i++ {
		h.Write(data)
	}
}
