package hash

import "C"
import (
	"hash"
	"hash/crc32"
	"math/rand"
	"testing"
	"time"

	"github.com/yerden/go-dpdk/util"
)

func makeGoHash32Func() func([]byte, uint32) uint32 {
	tab := crc32.MakeTable(crc32.Castagnoli)
	return func(data []byte, acc uint32) uint32 {
		return ^crc32.Update(^acc, tab, data)
	}
}

func testHash32(t *testing.T, src *rand.Rand, h1, h2 hash.Hash32, n int) {
	data := make([]byte, n)
	src.Read(data)

	h1.Reset()
	h1.Write(data)

	h2.Reset()
	h2.Write(data)

	if v2, v1 := h2.Sum32(), h1.Sum32(); v2 != v1 {
		t.Helper()
		t.Error("hash values not equal for n=", n, v2, "!=", v1)
	}
}

func testHash32Upd(t *testing.T, src *rand.Rand, h1, h2 func([]byte, uint32) uint32, n int) {
	data := make([]byte, n)
	src.Read(data)
	seed := src.Uint32()

	if v1, v2 := h1(data, seed), h2(data, seed); v2 != v1 {
		t.Helper()
		t.Error("hash values not equal for n=", n, v2, "!=", v1)
	}
}

func testAdditivity(t *testing.T, src *rand.Rand, hf func([]byte, uint32) uint32, n int) {
	data := make([]byte, n)
	src.Read(data)
	seed := src.Uint32()

	for i := 1; i < n-1; i++ {
		part1 := data[:i]
		part2 := data[i:]

		v1 := seed
		v1 = hf(part1, v1)
		v1 = hf(part2, v1)

		v2 := hf(data, seed)
		if v1 != v2 {
			t.Error("not additive")
		}
	}
}

func TestCrcAdditivity15(t *testing.T) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 2; i < 100; i++ {
		testAdditivity(t, src, CrcUpdate, i)
	}
}

func TestCrcGoAdditivity15(t *testing.T) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 2; i < 100; i++ {
		testAdditivity(t, src, makeGoHash32Func(), i)
	}
}

func TestCrc(t *testing.T) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))

	for seed := uint32(0); seed < 1000; seed++ {
		rteCrc32 := &util.Hash32{
			Seed:  seed,
			Block: 8,
			Accum: CrcUpdate,
		}

		goHashCrc32 := &util.Hash32{
			Seed:  seed,
			Block: 8,
			Accum: makeGoHash32Func(),
		}

		for n := 0; n < 100; n++ {
			testHash32(t, src, rteCrc32, goHashCrc32, n)
		}
	}
}

func TestCrcComplement(t *testing.T) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))

	rteCrc32 := &util.Hash32{
		Seed:  0,
		Block: 8,
		Accum: complementFunc(CrcUpdate),
	}

	goHashCrc32 := crc32.New(crc32.MakeTable(crc32.Castagnoli))

	for n := 0; n < 100; n++ {
		testHash32(t, src, rteCrc32, goHashCrc32, n)
	}
}

func TestCrcUpdate(t *testing.T) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))

	goCrc32 := makeGoHash32Func()

	for n := 1; n < 100; n++ {
		testHash32Upd(t, src, CrcUpdate, goCrc32, n)
	}
}

func benchmarkHash32(b *testing.B, src *rand.Rand, h hash.Hash32, block int) {
	data := make([]byte, block)
	src.Read(data)

	h.Reset()

	for i := 0; i < b.N; i++ {
		h.Write(data)
	}
}

func benchmarkHash32Func(b *testing.B, f func([]byte, uint32) uint32, block int) {
	src := rand.New(rand.NewSource(time.Now().UnixNano()))

	testHash := &util.Hash32{
		Seed:  0,
		Block: 8,
		Accum: f,
	}

	benchmarkHash32(b, src, testHash, block)
}

func BenchmarkCrcUpdate8(b *testing.B) {
	benchmarkHash32Func(b, CrcUpdate, 8)
}

func BenchmarkCrcUpdate31(b *testing.B) {
	benchmarkHash32Func(b, CrcUpdate, 31)
}

func BenchmarkGoCrc32Update8(b *testing.B) {
	benchmarkHash32Func(b, makeGoHash32Func(), 8)
}

func BenchmarkGoCrc32Update31(b *testing.B) {
	benchmarkHash32Func(b, makeGoHash32Func(), 31)
}

func CrcUint32(ip, seed uint32) uint32 {
	return uint32(C.rte_hash_crc_4byte(C.uint32_t(ip), C.uint32_t(seed)))
}
