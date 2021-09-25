package util

import "hash"

var _ hash.Hash32 = (*Hash32)(nil)

// Hash32 wraps a function Accum to implement hash.Hash32 interface.
type Hash32 struct {
	// Seed is the initial hash32 value.
	Seed uint32

	// Value is the calculated hash32 value. To reset the Hash32
	// state, simply reset this to Seed.
	Value uint32

	// Optimal block size. See BlockSize method of hash.Hash32
	// interface.
	Block int

	// If true then previous hash value is complemented to one before
	// specifying to Accum and the output is then complemented to one.
	OnesComplement bool

	// Accum is the backend hash implementation. It takes input byte
	// array and calculates new hash value based on previous one.
	Accum func([]byte, uint32) uint32
}

// Write implements io.Writer interface needed for hash.Hash32
// implementation.
func (h *Hash32) Write(p []byte) (int, error) {
	if h.OnesComplement {
		h.Value = ^h.Accum(p, ^h.Value)
	} else {
		h.Value = h.Accum(p, h.Value)
	}
	return len(p), nil
}

// Sum implements hash.Hash32 interface.
func (h *Hash32) Sum(b []byte) []byte {
	s := h.Value
	return append(b, byte(s>>24), byte(s>>16), byte(s>>8), byte(s))
}

// Sum32 implements hash.Hash32 interface.
func (h *Hash32) Sum32() uint32 {
	return h.Value
}

// Reset implements hash.Hash32 interface.
func (h *Hash32) Reset() {
	h.Value = h.Seed
}

// Size implements hash.Hash32 interface. Returns 4.
func (h *Hash32) Size() int {
	return 4
}

// BlockSize implements hash.Hash32 interface.
func (h *Hash32) BlockSize() int {
	return h.Block
}
