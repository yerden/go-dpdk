package port

import (
	"testing"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ring"
)

var initEAL = common.DoOnce(func() error {
	var set unix.CPUSet
	err := unix.SchedGetaffinity(0, &set)
	if err == nil {
		_, err = eal.Init([]string{"test",
			"-c", common.NewMap(&set).String(),
			"-m", "128",
			"--no-huge",
			"--no-pci",
			"--master-lcore", "0"})
	}
	return err
})

func TestPortRingIn(t *testing.T) {
	assert := common.Assert(t, true)
	m := common.NewAllocatorSession(&common.StdAlloc{})
	defer m.Flush()

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	err := eal.ExecOnMaster(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		var ops *InOps
		var port *In
		var arg *InArg
		confIn := &RingIn{Ring: r, Multi: true}

		arg = confIn.Arg(m)
		assert(arg != nil)

		ops = confIn.Ops()
		assert(ops != nil, ops)

		port = ops.Create(-1, arg)

		assert(port != nil, port)
		err = ops.Free(port)
		assert(err == nil, err)
	})
	assert(err == nil, err)
}

func TestPortRingOut(t *testing.T) {
	assert := common.Assert(t, true)
	m := common.NewAllocatorSession(&common.StdAlloc{})
	defer m.Flush()

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	err := eal.ExecOnMaster(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		var ops *OutOps
		var port *Out
		var arg *OutArg
		confOut := &RingOut{Ring: r, Multi: true, NoDrop: false, TxBurstSize: 64}

		arg = confOut.Arg(m)
		assert(arg != nil)

		ops = confOut.Ops()
		assert(ops != nil, ops)

		port = ops.Create(-1, arg)

		assert(port != nil, port)
		err = ops.Free(port)
		assert(err == nil, err)
	})
	assert(err == nil, err)
}

func TestPortRingCreateIn(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	err := eal.ExecOnMaster(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		var ops *InOps
		var port *In
		confIn := &RingIn{Ring: r, Multi: true}
		ops, port = CreateIn(confIn, -1)
		assert(ops != nil, ops)
		assert(port != nil, port)
		assert(ops.Free(port) == nil)
	})
	assert(err == nil, err)
}

func TestPortRingCreateOut(t *testing.T) {
	assert := common.Assert(t, true)
	m := common.NewAllocatorSession(&common.StdAlloc{})
	defer m.Flush()

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	err := eal.ExecOnMaster(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		var ops *OutOps
		var port *Out
		confOut := &RingOut{Ring: r, Multi: true, NoDrop: false, TxBurstSize: 64}
		ops, port = CreateOut(confOut, -1)
		assert(ops != nil, ops)
		assert(port != nil, port)
		assert(ops.Free(port) == nil)
	})
	assert(err == nil, err)
}
