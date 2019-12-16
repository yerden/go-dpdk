package common_test

import (
	"fmt"
	"testing"

	"github.com/yerden/go-dpdk/common"

	"golang.org/x/sys/unix"
)

func TestCommonMapCreate(t *testing.T) {
	assert := common.Assert(t, true)

	a := common.NewMap([]int{0, 1, 2, 3})
	assert(a.IsSet(0))
	assert(a.IsSet(1))
	assert(a.IsSet(2))
	assert(a.IsSet(3))
	assert(a.Count() == 4)
	assert(a.String() == "f")

	a = common.NewMap([]int{4, 5, 7, 9})
	assert(a.IsSet(4))
	assert(a.IsSet(5))
	assert(!a.IsSet(6))
	assert(a.IsSet(7))
	assert(!a.IsSet(8))
	assert(a.IsSet(9))
	assert(a.Count() == 4)
	assert(a.String() == "2b0")

	a = common.NewMap([]int{6, 8})
	assert(a.String() == "140")
}

func TestCommonMapCreate2(t *testing.T) {
	assert := common.Assert(t, true)

	s := common.NewMap(1)
	assert(s.String() == "2", s.String())

	s = common.NewMap([]int{1, 2, 3})
	assert(s.String() == "e", s.String())

	s = common.NewMap([]int{4, 5, 6})
	assert(s.String() == "70", s.String())

	s = common.NewMap(map[uint16]bool{
		11: true,
		22: true,
		32: true,
	})
	assert(s.Count() == 3)
	assert(s.IsSet(11))
	assert(s.IsSet(22))
	assert(s.IsSet(32))
}

func TestCommonMapSet(t *testing.T) {
	assert := common.Assert(t, true)

	a := common.NewMap([]int{0, 1, 2, 3})
	assert(a.IsSet(0))
	assert(a.IsSet(1))
	assert(a.IsSet(2))
	assert(a.IsSet(3))
	assert(a.Count() == 4)
	assert(a.String() == "f")

	a.Set(1)
	assert(a.Count() == 4)
	a.Set(2)
	assert(a.Count() == 4)
	a.Set(4) // new
	assert(a.Count() == 5)
	assert(a.IsSet(4))

	a.Zero()
	assert(a.Count() == 0)
}

func TestCommonMapFromSet(t *testing.T) {
	assert := common.Assert(t, true)

	a := common.NewMap(nil)
	assert(a.Count() == 0)

	var set unix.CPUSet
	err := unix.SchedGetaffinity(0, &set)
	assert(err == nil, err)

	a = common.NewMap(&set)
	assert(a.Count() == set.Count())
}

func ExampleNewMap() {
	x := common.NewMap([]int{0, 1, 2, 3})
	fmt.Println(x)
	// Output: f
}
