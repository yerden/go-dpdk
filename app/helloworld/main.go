package main

import (
	"golang.org/x/sys/unix"
	"log"
	"os"
	"sync"

	"github.com/yerden/go-dpdk/eal"
)

func main() {
	var set unix.CPUSet

	if err := unix.SchedGetaffinity(0, &set); err != nil {
		panic(err)
	}

	err := eal.InitWithParams(os.Args[0],
		eal.NewParameter("-c", eal.NewMap(&set)),
		eal.NewParameter("--no-huge"),
		eal.NewParameter("--no-pci"),
		eal.NewParameter("--master-lcore", 0),
	)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, id := range eal.Lcores(false) {
		wg.Add(1)
		eal.ExecuteOnLcore(id, func(lcore *eal.Lcore) {
			defer wg.Done()
			log.Println("Executed on lcore", lcore.ID())
		})
	}
	wg.Wait()

	if err := eal.Cleanup(); err != nil {
		panic(err)
	}
}
