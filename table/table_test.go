package table

import (
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
)

func TestTableHash(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		hash := &HashParams{
			TableOps:   HashCuckooOps,
			Name:       "hello",
			KeySize:    32,
			KeysNum:    1024,
			BucketsNum: 1024,
		}

		hash.CuckooHash.Func = Crc32Hash

		table1 := Create(0, hash, 20)
		assert(table1 != nil)

		err := table1.Free(hash.TableOps)
		assert(err == nil)
	})
	assert(err == nil, err)
}
