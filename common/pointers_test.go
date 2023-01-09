package common

import (
	"testing"
	"unsafe"
)

const someByte = 0xDC

func initBytes(b []byte) {
	for i := range b {
		b[i] = someByte
	}
}

func checkBytesNil(t testing.TB, b []byte, n int) {
	t.Helper()

	for i := 0; i < n; i++ {
		if b[i] != 0 {
			t.Fatalf("byte array of size %d not zeroed\n", n)
		}
	}

	if b[n] != someByte {
		t.Fatalf("invalid initialization of array of size %d\n", n)
	}
}

func TestMemset(t *testing.T) {
	b := make([]byte, 1024)

	for i := 1; i < len(b)-1; i++ {
		initBytes(b)
		Memset(unsafe.Pointer(&b[0]), 0, uintptr(i))
		checkBytesNil(t, b, i)
	}
}

var globalSlice []byte

func benchmarkMemsetN(b *testing.B, n int) {
	b.Helper()
	data := make([]byte, n)
	initBytes(data)

	for i := 0; i < b.N; i++ {
		Memset(unsafe.Pointer(&data[0]), 0, uintptr(len(data)))
	}

	globalSlice = data
}

func benchmarkPlainInitN(b *testing.B, n int) {
	b.Helper()
	data := make([]byte, n)
	initBytes(data)

	for i := 0; i < b.N; i++ {
		for j := range data {
			data[j] = 0
		}
	}

	globalSlice = data
}

func BenchmarkMemset4(b *testing.B) {
	benchmarkMemsetN(b, 4)
}

func BenchmarkPlainInit4(b *testing.B) {
	benchmarkPlainInitN(b, 4)
}

func BenchmarkMemset6(b *testing.B) {
	benchmarkMemsetN(b, 6)
}

func BenchmarkPlainInit6(b *testing.B) {
	benchmarkPlainInitN(b, 6)
}

func BenchmarkMemset8(b *testing.B) {
	benchmarkMemsetN(b, 8)
}

func BenchmarkPlainInit8(b *testing.B) {
	benchmarkPlainInitN(b, 8)
}

func BenchmarkMemset12(b *testing.B) {
	benchmarkMemsetN(b, 12)
}

func BenchmarkPlainInit12(b *testing.B) {
	benchmarkPlainInitN(b, 12)
}

func BenchmarkMemset31(b *testing.B) {
	benchmarkMemsetN(b, 31)
}

func BenchmarkPlainInit31(b *testing.B) {
	benchmarkPlainInitN(b, 31)
}

func BenchmarkMemset511(b *testing.B) {
	benchmarkMemsetN(b, 511)
}

func BenchmarkPlainInit511(b *testing.B) {
	benchmarkPlainInitN(b, 511)
}
