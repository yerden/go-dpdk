package main

import (
	"log"
	"os"
	"sync"
	"unsafe"

	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ring"
)

func main() {
	if _, err := eal.Init(os.Args); err != nil {
		log.Fatalln("EAL init failed:", err)
	}
	defer eal.Cleanup()

	// create a ring
	var r *ring.Ring
	var err error
	e := eal.ExecOnMaster(func(*eal.LcoreCtx) {
		r, err = ring.Create("test_ring", 1024)
	})
	if e != nil {
		panic(e)
	} else if err != nil {
		panic(err)
	}
	defer eal.ExecOnMaster(func(*eal.LcoreCtx) { r.Free() })

	// Pick first slave, panic if none.
	slave := eal.LcoresSlave()[0]

	// start sending and receiving messages
	var wg sync.WaitGroup
	wg.Add(1)
	n := 1000000 // 1M messages
	go func() {
		defer wg.Done()
		eal.ExecOnMaster(func(ctx *eal.LcoreCtx) {
			for i := 0; i < n; {
				if r.Enqueue(unsafe.Pointer(r)) {
					i++
				}
			}
			log.Println("sent", n, "messages")
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		eal.ExecOnLcore(slave, func(ctx *eal.LcoreCtx) {
			for i := 0; i < n; {
				ptr, ok := r.Dequeue()
				ok = ok && unsafe.Pointer(r) == ptr
				if ok {
					i++
				}
			}
			log.Println("received", n, "messages")
		})
	}()
	wg.Wait()
}
