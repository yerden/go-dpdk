package eal_test

import (
	"golang.org/x/sys/unix"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
)

func TestCommonSets(t *testing.T) {
	assert := common.Assert(t, true)

	var s eal.Set
	s = eal.MakeSet(1)
	assert(s.IsSet(1))
	assert(!s.IsSet(0))
	assert(!s.IsSet(2))

	s = eal.MakeSet([]int{1, 2, 3})
	assert(s.IsSet(1))
	assert(s.IsSet(2))
	assert(s.IsSet(3))
	assert(!s.IsSet(4))
	assert(!s.IsSet(0))

	s = eal.MakeSet([]uint16{4, 5, 6})
	assert(s.IsSet(4))
	assert(s.IsSet(5))
	assert(s.IsSet(6))
	assert(!s.IsSet(1))
	assert(!s.IsSet(10))

	s = eal.MakeSet(map[uint16]bool{
		11: true,
		22: true,
		32: true,
	})
	assert(s.IsSet(11))
	assert(s.IsSet(22))
	assert(s.IsSet(32))
	assert(!s.IsSet(124))
	assert(!s.IsSet(3))
}

func TestCommonSetToHex(t *testing.T) {
	assert := common.Assert(t, true)
	var set unix.CPUSet

	set.Zero()
	set.Set(0)
	assert(eal.SetToHex(&set, 128) == "1")

	set.Set(1)
	assert(eal.SetToHex(&set, 128) == "3")

	set.Set(2)
	assert(eal.SetToHex(&set, 128) == "7")
}
