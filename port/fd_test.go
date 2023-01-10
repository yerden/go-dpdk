package port

import (
	"os"
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
)

func TestPortFdRx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	f, err := os.CreateTemp("", "testfile")
	assert(err == nil)
	defer os.Remove(f.Name())
	defer f.Close()

	err = eal.ExecOnMain(func(*eal.LcoreCtx) {
		mp, err := mempool.CreateMbufPool(
			"hello",
			1024,
			2048,
		)
		assert(err == nil, err)
		defer mp.Free()

		params := &FdRx{
			Mempool: mp,
			Fd:      f.Fd(),
			MTU:     1 << 12,
		}

		ops := params.InOps()
		assert(ops != nil)

		in := CreateIn(-1, params)
		assert(in != nil)

		assert(nil == in.Free(ops))
	})
	assert(err == nil, err)
}

func TestPortFdOut(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	f, err := os.Open("/dev/null")
	assert(err == nil)
	defer f.Close()

	err = eal.ExecOnMain(func(*eal.LcoreCtx) {
		params := &FdTx{
			Fd:      f.Fd(),
			BurstSize: 32,
			NoDrop:  true,
			Retries: 32,
		}

		ops := params.OutOps()
		assert(ops != nil)

		out := CreateOut(-1, params)
		assert(out != nil)

		assert(nil == out.Free(ops))
	})
	assert(err == nil, err)
}
