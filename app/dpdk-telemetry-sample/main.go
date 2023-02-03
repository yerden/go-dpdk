package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"syscall"

	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/memzone"
	"github.com/yerden/go-dpdk/ring"
	"github.com/yerden/go-dpdk/telemetry"
)

func main() {
	fmt.Println("vim-go")

	_, err := eal.Init(os.Args)
	if err != nil {
		panic(err)
	}

	telemetry.RegisterCmd("/ring/list", "Show list of Create()-d rings",
		func(cmd, params string, d *telemetry.Data) int {
			d.StartArray(telemetry.StringVal)
			memzone.Walk(func(mz *memzone.Memzone) {
				if name := strings.TrimPrefix(mz.Name(), "RG_"); name != mz.Name() {
					d.AddArrayString(name)
				}
			})
			return 0
		})

	telemetry.RegisterCmd("/ring/info", "Show ring information. Param: ring name",
		func(cmd, params string, d *telemetry.Data) int {
			r, err := ring.Lookup(params)
			if err != nil {
				return -int(syscall.ENOENT)
			}

			d.StartDict()
			d.AddDictString("ring_name", r.Name())
			d.AddDictInt("ring_count", int(r.Count()))
			d.AddDictInt("ring_capacity", int(r.Cap()))
			return 0
		})

	var wg sync.WaitGroup
	wg.Add(1)
	err = eal.ExecOnMain(func(*eal.LcoreCtx) {
		// defer wg.Done() // switched off Done
		_, e := ring.Create("test_ring", 1024)
		if e != nil {
			panic(e)
		}
	})

	if err != nil {
		panic(err)
	}

	wg.Wait()
}
