package table

import (
	"testing"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
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
			"--main-lcore", "0"})
	}
	return err
})

func TestTableHash(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

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
