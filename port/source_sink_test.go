package port

import (
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
)

func TestPortSource(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		mp, err := mempool.CreateMbufPool(
			"hello",
			1024,
			2048,
		)
		assert(err == nil, err)
		defer mp.Free()

		params := &Source{
			Mempool: mp,
		}

		ops := params.InOps()
		assert(ops != nil)

		in := CreateIn(-1, params)
		assert(in != nil)

		assert(nil == in.Free(ops))
	})
	assert(err == nil, err)
}

func TestPortSink(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		params := &Sink{
			Filename: "/dev/null",
		}

		ops := params.OutOps()
		assert(ops != nil)

		out := CreateOut(-1, params)
		assert(out != nil)

		assert(nil == out.Free(ops))
	})
	assert(err == nil, err)
}
