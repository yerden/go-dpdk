package port

import (
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ring"
)

func TestPortRingRx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		confRx := &RingRx{Ring: r, Multi: true}

		ops := confRx.InOps()
		assert(ops != nil)

		arg, dtor := confRx.Transform(alloc)
		defer dtor(arg)

		in := CreateIn(-1, confRx)
		assert(in != nil)

		err = in.Free(ops)
		assert(err == nil, err)
	})
	assert(err == nil, err)
}

func TestPortRingTx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		confTx := &RingTx{Ring: r, Multi: true, NoDrop: false, TxBurstSize: 64}

		ops := confTx.OutOps()
		tx := CreateOut(-1, confTx)
		assert(tx != nil)

		err = tx.Free(ops)
		assert(err == nil, err)
	})
	assert(err == nil, err)
}

func TestPortRingCreateRx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		confRx := &RingRx{Ring: r, Multi: true}
		ops := confRx.InOps()
		rx := CreateIn(-1, confRx)
		assert(rx != nil)
		assert(rx.Free(ops) == nil)
	})
	assert(err == nil, err)
}

func TestPortRingCreateTx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		r, err := ring.Create("test_ring", 1024)
		assert(err == nil, err)
		defer r.Free()

		confTx := &RingTx{Ring: r, Multi: true, NoDrop: false, TxBurstSize: 64}
		ops := confTx.OutOps()
		tx := CreateOut(-1, confTx)
		assert(tx != nil, tx)
		assert(tx.Free(ops) == nil)
	})
	assert(err == nil, err)
}
