package pipeline

import (
	"testing"

	"golang.org/x/sys/unix"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
	"github.com/yerden/go-dpdk/port"
	"github.com/yerden/go-dpdk/ring"
	"github.com/yerden/go-dpdk/table"
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

func TestPortRingRx(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		pl := Create(&Params{
			Name:         "test_pipeline",
			SocketID:     0,
			OffsetPortID: 0,
		})
		assert(pl != nil)

		mp, err := mempool.CreateMbufPool(
			"hello",
			1024,
			2048,
		)
		assert(err == nil, err)

		pSource1, err := pl.PortInCreate(&PortInParams{
			Params: &port.Source{
				Mempool: mp,
			},
			BurstSize: 32,
		})

		assert(err == nil, err)

		table1, err := pl.TableCreate(&TableParams{
			Params: &table.ArrayParams{
				Entries: 16,
			},
		})
		assert(err == nil, err)

		err = pl.ConnectToTable(pSource1, table1)
		assert(err == nil, err)

		pSink1, err := pl.PortOutCreate(&PortOutParams{
			Params: &port.Sink{
				Filename:   "/dev/null",
				MaxPackets: 32,
			},
		})

		assert(pSink1 == 0)
		assert(nil == err, err)

		// sample entry
		entry := &TableEntry{}
		entry.SetAction(ActionPort)
		entry.SetPortID(pSink1)

		dfltEntry, err := pl.TableDefaultEntryAdd(table1, entry)
		assert(err == nil, err)
		assert(*dfltEntry == *entry)

		assert(nil == pl.Check())
		assert(nil == pl.Disable(pSource1))
		assert(nil == pl.Flush())

		// destroy the tables, ports and pipeline
		assert(nil == pl.Free())
	})
	assert(err == nil, err)
}

func TestPipelineStub(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	assert(initEAL() == nil)

	err := eal.ExecOnMain(func(*eal.LcoreCtx) {
		pl := Create(&Params{
			Name:         "test_pipeline",
			SocketID:     0,
			OffsetPortID: 0,
		})
		assert(pl != nil)

		r, err := ring.Create("test_ring", 1024, ring.OptSC)
		assert(err == nil, err)
		assert(r != nil)
		defer r.Free()

		// create pipeline ports
		pRing, err := pl.PortInCreate(&PortInParams{
			Params:    &port.RingRx{Ring: r, Multi: false},
			BurstSize: 32,
		})
		assert(err == nil, err)

		pSink1, err := pl.PortOutCreate(&PortOutParams{
			Params: &port.Sink{
				Filename:   "/dev/null",
				MaxPackets: 32,
			},
		})

		assert(pSink1 == 0)
		assert(nil == err, err)

		// create pipeline tables
		tStub, err := pl.TableCreate(&TableParams{
			Params: &table.StubParams{},
		})
		assert(err == nil, err)

		entry := NewTableEntry(0)
		entry.SetAction(2) // C.RTE_PIPELINE_ACTION_PORT_META
		defaultEntry, err := pl.TableDefaultEntryAdd(tStub, entry)
		assert(err == nil, err)
		assert(*defaultEntry == *entry)

		// connect input port.
		err = pl.ConnectToTable(pRing, tStub)
		assert(err == nil, err)

		// check the pipeline
		assert(nil == pl.Check())
		assert(nil == pl.Flush())

		// destroy the tables, ports and pipeline
		assert(nil == pl.Free())
	})
	assert(err == nil, err)
}
