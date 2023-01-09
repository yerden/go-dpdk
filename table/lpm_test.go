package table

import (
	"testing"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
)

func TestLPMParams(t *testing.T) {
	assert := common.Assert(t, true)

	eal.InitOnceSafe("test", 4)

	params := &LPMParams{
		Name:            "hello",
		NumTBL8:         16,
		Rules:           1024,
		Offset:          0,
		EntryUniqueSize: 2,
	}

	{
		params.TableOps = LPMOps
		ops := params.Ops()
		tbl := Create(-1, params, 2)
		assert(tbl != nil)
		err := tbl.Free(ops)
		assert(err == nil)
	}

	{
		params.TableOps = LPM6Ops
		ops := params.Ops()
		tbl := Create(-1, params, 2)
		assert(tbl != nil)
		err := tbl.Free(ops)
		assert(err == nil)
	}
}
