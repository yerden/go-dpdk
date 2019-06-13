// +build linux

package lcore

import (
	"fmt"
	"os"
	"syscall"
)

var (
	lcore2node []int
)

// linux implementation, should return maximum lcore id + 1.
func loadCPUMap() int {
	nodeMap := make(map[uint]int)
	maxLcoreID := 0
	for n := 0; ; n++ {
		path := fmt.Sprintf("/sys/devices/system/node/node%d/cpumap", n)
		f, err := os.Open(path)
		if err != nil {
			if e, ok := err.(*os.PathError); !ok || e.Err != syscall.ENOENT {
				panic(e)
			}
			break
		}
		defer f.Close()

		cores, err := readCPUHexMap(f)
		if err != nil {
			panic(err)
		}

		for _, c := range cores {
			if nodeMap[uint(c)] = n; c > maxLcoreID {
				maxLcoreID = c
			}
		}
	}

	lcore2node = make([]int, maxLcoreID+1)
	for lcore, node := range nodeMap {
		lcore2node[lcore] = node
	}

	return len(lcore2node)
}

func getNumaNode(id uint) int {
	if id < uint(len(lcore2node)) {
		return lcore2node[id]
	}
	return NumaNodeAny
}
