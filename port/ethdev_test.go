package port

import (
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
)

func TestPortEthdevRx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		params := &EthdevRx{
			PortID:  0,
			QueueID: 0,
		}

		ops := params.InOps()
		assert(ops != nil)

		in := CreateIn(-1, params)
		assert(in != nil)

		assert(nil == in.Free(ops))
	})
	assert(err == nil, err)
}

func TestPortEthdevTx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		params := &EthdevTx{
			PortID:      0,
			QueueID:     0,
			TxBurstSize: 32,
			NoDrop:      true,
			Retries:     4,
		}

		ops := params.OutOps()
		assert(ops != nil)

		out := CreateOut(-1, params)
		assert(out != nil)

		assert(nil == out.Free(ops))
	})
	assert(err == nil, err)
}
