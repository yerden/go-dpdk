package pipeline

import (
	"io"
	"testing"
	"time"

	"github.com/yerden/go-dpdk/common"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/mempool"
	"github.com/yerden/go-dpdk/port"
	"github.com/yerden/go-dpdk/table"
)

func TestPipelineLoop(t *testing.T) {
	assert := common.Assert(t, true)

	// Initialize EAL on all cores
	eal.InitOnceSafe("test", 4)

	ctrl := NewSimpleController()
	errCh := make(chan error, 2)

	eal.ExecOnMainAsync(errCh, func(*eal.LcoreCtx) {
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
		defer mp.Free()

		pSource1, err := pl.PortInCreate(&PortInParams{
			Params: &port.Source{
				Mempool: mp,
			},
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

		// connect input port.
		err = pl.ConnectToTable(pSource1, tStub)
		assert(err == nil, err)

		// check the pipeline
		assert(nil == pl.Check())
		assert(nil == pl.Flush())

		errCh <- io.EOF // just some error to signal

		// run infinite loop on the pipeline, waiting for control
		err = pl.RunLoop(1024, ctrl)
		assert(err == nil, err)

		// destroy the tables, ports and pipeline
		assert(nil == pl.Free())
	})

	assert(io.EOF == <-errCh) // we're about to start pipeline loop

	// waiting for pipeline to soar
	time.Sleep(100 * time.Millisecond)

	// stop the pipeline loop
	ctrl.Stop()

	assert(nil == <-errCh)
}
