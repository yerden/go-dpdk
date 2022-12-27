package table

/*
#include <rte_hash.h>
#include <rte_hash_crc.h>
#include <rte_table_hash.h>
#include <rte_table_hash_cuckoo.h>

rte_hash_function go_crc32f = rte_hash_crc;

*/
import "C"

import (
	"unsafe"

	"github.com/yerden/go-dpdk/common"
)

var (
	// Crc32Hash is the DPDK CRC32 hash function.
	Crc32Hash HashFunc = (HashFunc)(C.go_crc32f)
)

// HashOpHash is the signature of hash function used for Hash table
// implementation.
type HashOpHash C.rte_table_hash_op_hash

// HashFunc is the signature of hash function used for Cuckoo-Hash
// table implementation.
type HashFunc C.rte_hash_function

// Table ops implementations.
var (
	HashExtOps    = (*Ops)(&C.rte_table_hash_ext_ops)
	HashLruOps    = (*Ops)(&C.rte_table_hash_lru_ops)
	HashCuckooOps = (*Ops)(&C.rte_table_hash_cuckoo_ops)
)

// HashParams is the Hash table parameters.
type HashParams struct {
	// Set this to desired Hash*Ops variable.
	TableOps *Ops

	// Name.
	Name string

	// Key size (number of bytes).
	KeySize uint32

	// Byte offset within packet meta-data where the key is located
	KeyOffset uint32

	// Key mask.
	KeyMask []byte

	// Number of keys.
	KeysNum uint32

	// Number of buckets.
	BucketsNum uint32

	// Hash function.
	Hash struct {
		Func HashOpHash
		Seed uint64
	}

	// Hash function for HashCuckooOps.
	CuckooHash struct {
		Func HashFunc
		Seed uint32
	}
}

// Transform implements common.Transformer interface.
func (p *HashParams) Transform(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	switch p.TableOps {
	case HashExtOps:
		fallthrough
	case HashLruOps:
		return p.tformToOrdinary(alloc)
	case HashCuckooOps:
		return p.tformToCuckoo(alloc)
	}

	panic("unsupported Hash Table ops")
}

func (p *HashParams) tformToOrdinary(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	var params *C.struct_rte_table_hash_params
	params = (*C.struct_rte_table_hash_params)(alloc.Malloc(unsafe.Sizeof(*params)))
	params.name = (*C.char)(common.CString(alloc, p.Name))
	params.key_size = C.uint32_t(p.KeySize)
	params.key_offset = C.uint32_t(p.KeyOffset)
	params.key_mask = (*C.uint8_t)(common.CBytes(alloc, p.KeyMask))
	params.n_keys = C.uint32_t(p.KeysNum)
	params.n_buckets = C.uint32_t(p.BucketsNum)
	params.f_hash = C.rte_table_hash_op_hash(p.Hash.Func)
	params.seed = C.uint64_t(p.Hash.Seed)
	return unsafe.Pointer(params), func(p unsafe.Pointer) {
		params := (*C.struct_rte_table_hash_params)(p)
		alloc.Free(unsafe.Pointer(params.name))
		alloc.Free(unsafe.Pointer(params.key_mask))
		alloc.Free(p)
	}
}

func (p *HashParams) tformToCuckoo(alloc common.Allocator) (unsafe.Pointer, func(unsafe.Pointer)) {
	var params *C.struct_rte_table_hash_cuckoo_params
	params = (*C.struct_rte_table_hash_cuckoo_params)(alloc.Malloc(unsafe.Sizeof(*params)))
	params.name = (*C.char)(common.CString(alloc, p.Name))
	params.key_size = C.uint32_t(p.KeySize)
	params.key_offset = C.uint32_t(p.KeyOffset)
	params.key_mask = (*C.uint8_t)(common.CBytes(alloc, p.KeyMask))
	params.n_keys = C.uint32_t(p.KeysNum)
	params.n_buckets = C.uint32_t(p.BucketsNum)
	params.f_hash = C.rte_hash_function(p.CuckooHash.Func)
	params.seed = C.uint32_t(p.CuckooHash.Seed)
	return unsafe.Pointer(params), func(p unsafe.Pointer) {
		params := (*C.struct_rte_table_hash_cuckoo_params)(p)
		alloc.Free(unsafe.Pointer(params.name))
		alloc.Free(unsafe.Pointer(params.key_mask))
		alloc.Free(p)
	}
}

// Ops implements Params interface.
func (p *HashParams) Ops() *Ops {
	return p.TableOps
}
