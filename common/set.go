package common

import (
	"reflect"
	"sort"
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

func hex(a []int) string {
	var out []byte
	for _, n := range a {
		i, r := n/4, uint(n&3)
		if i >= len(out) {
			add := make([]byte, i-len(out)+1)
			for k := range add {
				add[k] = '0'
			}
			out = append(add, out...)
		}
		i = len(out) - 1 - i
		out[i] = getHexByte(getHexIndex(out[i]) | (1 << r))
	}

	return string(out)
}

// Set represents a set of integer numbers.
type Set interface {
	// IsSet tests if given int is inside the Set.
	IsSet(int) bool
	// Count returns the number of integers in the Set.
	Count() int
	// Set stores integer in the Set
	Set(int)
}

// Map is an []int array-based implementation of a Set.
type Map struct {
	array []int
}

func (m *Map) find(n int) (int, bool) {
	k := sort.SearchInts(m.array, n)
	return k, k < len(m.array) && m.array[k] == n
}

// Set implements Set interface.
func (m *Map) Set(n int) {
	if k, ok := m.find(n); !ok {
		m.array = append(m.array, n)
		copy(m.array[k+1:], m.array[k:])
		m.array[k] = n
	}
}

// IsSet implements Set interface.
func (m *Map) IsSet(n int) bool {
	_, ok := m.find(n)
	return ok
}

// Zero zeroes out Map.
func (m *Map) Zero() {
	m.array = m.array[:0]
}

// Count implements Set interface.
func (m *Map) Count() int {
	return len(m.array)
}

// String implements fmt.Stringer interface.
func (m *Map) String() string {
	return hex(m.array)
}

// copySet copies non-negative members of src to dst.
func copySet(dst, src Set) int {
	var n int
	for i := 0; n < src.Count(); i++ {
		if src.IsSet(i) {
			dst.Set(i)
			n++
		}
	}
	return n
}

// NewMap creates instance of a Map.
//
// i may represent a Set, an array or a slice of integers, a map with
// integer keys. Otherwise, the function would panic.
func NewMap(i interface{}) *Map {
	m := &Map{}

	if i == nil {
		return m
	}

	if s, ok := i.(Set); ok {
		copySet(m, s)
		return m
	}

	intType := reflect.ValueOf(int(0)).Type()

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for _, k := range keys {
			m.Set(int(k.Convert(intType).Int()))
		}
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		for n := 0; n < v.Len(); n++ {
			elem := v.Index(n).Convert(intType).Int()
			m.Set(int(elem))
		}
	default:
		elem := v.Convert(intType).Int()
		m.Set(int(elem))
	}
	return m
}
