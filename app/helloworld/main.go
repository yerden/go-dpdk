package main

import (
	"log"
	"os"

	"github.com/yerden/go-dpdk/eal"
)

func main() {
	if _, err := eal.Init(os.Args); err != nil {
		log.Fatalln("EAL init failed:", err)
	}
	defer eal.Cleanup()
	defer eal.StopLcores()

	for _, id := range eal.Lcores() {
		eal.ExecOnLcore(id, func(ctx *eal.LcoreCtx) {
			log.Println("hello from core", eal.LcoreID())
		})
	}
}
