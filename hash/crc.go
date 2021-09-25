/*
Package hash wraps RTE hash libraries.

Please refer to DPDK Programmer's Guide for reference and caveats.
*/
package hash

/*
#include <rte_config.h>
#include <rte_hash_crc.h>
*/
import "C"
import (
	"hash"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/yerden/go-dpdk/util"
)

const (
	// Crc32SW flag means 'Don't use SSE4.2 intrinsics'.
	Crc32SW uint8 = C.CRC32_SW

	// Crc32Sse42 means 'Use SSE4.2 intrinsics if available'.
	Crc32Sse42 uint8 = C.CRC32_SSE42

	// Crc32Sse42X64 means 'Use 64-bit SSE4.2 intrinsic if available
	// (default)'.
	Crc32Sse42X64 uint8 = C.CRC32_SSE42_x64
)

// CrcSetAlg allow or disallow use of SSE4.2 instrinsics for CRC32
// hash calculation. Specify OR of declared Crc32* flags.
func CrcSetAlg(alg uint8) {
	C.rte_hash_crc_set_alg(C.uint8_t(alg))
}

// CrcUpdate calculate CRC32 hash on user-supplied byte array.
func CrcUpdate(data []byte, acc uint32) uint32 {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	p, length := unsafe.Pointer(sh.Data), C.uint(sh.Len)
	//p, length := unsafe.Pointer(&data[0]), C.uint(len(data))
	crc := uint32(C.rte_hash_crc(p, length, C.uint(acc)))
	runtime.KeepAlive(data)
	return crc
}

func complementFunc(f func([]byte, uint32) uint32) func([]byte, uint32) uint32 {
	return func(data []byte, acc uint32) uint32 {
		return ^f(data, ^acc)
	}
}

// NewCrcHash creates new hash.Hash32 implemented on top of CrcUpdate.
// If complement is true then accumulator argument is complemented to
// 1 along with the output. This gives identical output as native Go
// crc32's implementation.
func NewCrcHash(seed uint32, complement bool) hash.Hash32 {
	h := &util.Hash32{
		Seed:  seed,
		Value: seed,
		Accum: CrcUpdate,
		Block: 8,
	}

	if complement {
		h.Accum = complementFunc(CrcUpdate)
	}

	return h
}
