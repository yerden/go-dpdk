package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/mempool"
	"github.com/yerden/go-dpdk/port"
	"github.com/yerden/go-dpdk/util"
)

var queues = flag.Int("nq", 1, "Number of RX queues")

func main() {
	n, err := eal.Init(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	defer eal.Cleanup()
	defer eal.StopLcores()

	os.Args[n], os.Args = os.Args[0], os.Args[n:]
	flag.Parse()

	// mempool
	var mp *mempool.Mempool

	rxqueue := []port.EthdevRx{}

	eal.ExecOnMain(func(ctx *eal.LcoreCtx) {
		for i := 0; i < ethdev.CountTotal(); i++ {
			var err error
			pid := ethdev.Port(i)
			if !pid.IsValid() {
				continue
			}

			name := fmt.Sprintf("test_mbuf_pool_%d", i)

			// create mbuf mempool
			mp, err = mempool.CreateMbufPool(name,
				20000, // elements count
				2048,  // size of element
				mempool.OptSocket(int(eal.SocketID())),
				mempool.OptCacheSize(512),
				mempool.OptOpsName("stack"),
				mempool.OptPrivateDataSize(64), // for each mbuf
			)
			if err != nil {
				log.Fatal("error creating mempool:", err)
			}

			// configure Ethernet device
			if err = pid.DevConfigure(uint16(*queues), 0); err != nil {
				log.Fatal(err)
			}

			log.Printf("port %d/%d configured, driver='%s'\n", pid, ethdev.CountTotal(), driverName(pid))

			qs := queuesSeq(*queues)

			log.Printf("port %d: configuring RX queues...\n", pid)
			// configure RX queues
			for _, q := range qs {
				if err = pid.RxqSetup(q, 256, mp); err != nil {
					log.Fatalf("unable to setup rxq #%d: %v\n", q, err)
				}

				rxqueue = append(rxqueue, port.EthdevRx{PortID: uint16(pid), QueueID: q})
			}

			// starting device
			log.Printf("port %d: starting...\n", pid)
			if err := pid.Start(); err != nil {
				log.Fatalf("unable to start port %d: %v\n", pid, err)
			}

			// setup RSS
			log.Printf("port %d: setup RSS...\n", pid)
			switch driverName(pid) {
			case "net_ice":
				if _, err := rssEthVlanIPv4(pid); err != nil {
					log.Fatalf("unable to setup rss on port %d: %v\n", pid, err)
				}
			case "net_mlx5":
				// TODO
			default:
				log.Printf("no RSS configured\n")
			}
		}

		log.Println("lcores:", eal.Lcores())
	})

	var wg sync.WaitGroup
	for _, lcore := range eal.LcoresWorker() {
		var rxf port.EthdevRx
		if len(rxqueue) == 0 {
			break
		}

		rxf, rxqueue = rxqueue[0], rxqueue[1:]
		wg.Add(1)

		go func(id uint, rxf port.EthdevRx) {
			defer wg.Done()

			log.Println("executing on lcore", id)
			failOnErr(eal.ExecOnLcore(id, func(ctx *eal.LcoreCtx) {
				socket := eal.LcoreToSocket(id)
				pid := ethdev.Port(rxf.PortID)

				rx := util.NewRxBuffer(pid, rxf.QueueID, int(socket), 512)
				defer rx.Free()

				log.Printf("start reading port %d/rxq %d on lcore %d", rxf.PortID, rxf.QueueID, id)
				xstatNames, err := pid.XstatNames()
				if err != nil {
					log.Fatal(err)
				}

				xstats := make([]ethdev.Xstat, len(xstatNames))

				for n := 0; ; n++ {
					if n%1e6 == 0 {
						n, err := pid.XstatsGet(&xstats)
						if err != nil {
							log.Println(err)
						} else {
							for _, xstat := range xstats[:n] {
								log.Println(&xstatNames[xstat.Index], xstat.Value)
							}
						}
					}

					_, ci, err := rx.ZeroCopyReadPacketData()
					if err != nil {
						continue
					}

					log.Println(&ci)

				}
			}))
		}(lcore, rxf)
	}

	if len(rxqueue) > 0 {
		log.Println("queues left:", len(rxqueue))
	}

	log.Println("waiting...")
	wg.Wait()
}
