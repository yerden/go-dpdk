package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/yerden/go-dpdk/eal"
	"github.com/yerden/go-dpdk/ethdev"
	"github.com/yerden/go-dpdk/util"
)

var burstSize = flag.Int("burst", 256, "Specify RX burst size")
var printMetadata = flag.Bool("print", false, "Specify to print each packet's metadata")
var dryRun = flag.Bool("dryRun", false, "If true traffic will not be processed")

// PortQueue describes port and rx queue id.
type PortQueue struct {
	Pid ethdev.Port
	Qid uint16
}

// dissect all given lcores and store them into map hashed by affine
// socket id.
func dissectLcores(lcores []uint) map[uint][]uint {
	table := map[uint][]uint{}

	for _, lcore := range lcores {
		socket := eal.LcoreToSocket(lcore)

		if affine, ok := table[socket]; !ok {
			table[socket] = []uint{lcore}
		} else {
			table[socket] = append(affine, lcore)
		}
	}

	return table
}

// DistributeQueues assigns all RX queues for each port in ports to
// lcores. Assignment is NUMA-aware.
//
// Returns os.ErrInvalid if port id is invalid.
// Returns os.ErrNotExist if no lcores are available by NUMA
// constraints.
func DistributeQueues(ports []ethdev.Port, lcores []uint) (map[uint]PortQueue, error) {
	table := map[uint]PortQueue{}
	lcoreMap := dissectLcores(lcores)

	for _, pid := range ports {
		if err := distributeQueuesPort(pid, lcoreMap, table); err != nil {
			return nil, err
		}
	}

	return table, nil
}

func distributeQueuesPort(pid ethdev.Port, lcoreMap map[uint][]uint, table map[uint]PortQueue) error {
	var info ethdev.DevInfo

	if err := pid.InfoGet(&info); err != nil {
		return err
	}

	socket := pid.SocketID()
	if socket < 0 {
		return os.ErrInvalid
	}

	lcores, ok := lcoreMap[uint(socket)]
	if !ok {
		fmt.Println("no lcores for socket:", socket)
		return os.ErrNotExist
	}

	nrx := info.NbRxQueues()
	if nrx == 0 {
		return os.ErrClosed
	}

	if int(nrx) > len(lcores) {
		return fmt.Errorf("pid=%d nrx=%d cannot run on %d lcores", pid, nrx, len(lcores))
	}

	var lcore uint
	var acquired util.LcoresList
	for i := uint16(0); i < nrx; i++ {
		lcore, lcores = lcores[0], lcores[1:]
		acquired = append(acquired, lcore)
		lcoreMap[uint(socket)] = lcores
		table[lcore] = PortQueue{Pid: pid, Qid: i}
	}

	fmt.Printf("pid=%d runs on socket=%d, lcores=%v\n", pid, socket, util.LcoresList(acquired))

	return nil
}

func LcoreFunc(pq PortQueue, qcr *QueueCounterReporter) func(*eal.LcoreCtx) {
	return func(ctx *eal.LcoreCtx) {
		defer log.Println("lcore", eal.LcoreID(), "exited")

		if *dryRun {
			return
		}

		// parser
		var (
			eth  layers.Ethernet
			vlan layers.Dot1Q
			ip4  layers.IPv4
			ip6  layers.IPv6
			gtpu layers.GTPv1U
		)

		var dlc gopacket.DecodingLayerContainer
		dlc = gopacket.DecodingLayerSparse{}
		dlc = dlc.Put(&eth)
		dlc = dlc.Put(&vlan)
		dlc = dlc.Put(&ip4)
		dlc = dlc.Put(&ip6)
		dlc = dlc.Put(&gtpu)

		parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet)
		parser.SetDecodingLayerContainer(dlc)
		parser.IgnoreUnsupported = true

		// eal
		pid := pq.Pid
		qid := pq.Qid
		qc := qcr.Register(pid, qid)

		src := util.NewEthdevMbufArray(pid, qid, int(eal.SocketID()), uint16(*burstSize))
		defer src.Free()

		buf := src.Buffer()

		decoded := make([]gopacket.LayerType, 10)

		log.Printf("processing pid=%d, qid=%d, lcore=%d\n", pid, qid, eal.LcoreID())
		for {
			n := src.Recharge()

			for i := 0; i < n; i++ {
				data := buf[i].Data()

				if err := parser.DecodeLayers(data, &decoded); err != nil {
					log.Println("parsing error:", err)
					continue
				}

				if *printMetadata {
					fmt.Printf("packet: %d bytes\n", len(data))
				}
			}

			qc.Incr(buf[:n])
		}

	}
}
