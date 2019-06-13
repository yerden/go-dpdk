package lcore

import (
	"bufio"
	"errors"
	"io"
)

// Default NUMA node value.
const (
	NumaNodeAny = -1
)

// Maximum number of logical CPU cores on the system.
var (
	MaxLcoreID = loadCPUMap()
)

var (
	hexMap = map[byte]int{
		'0': 0x0, '1': 0x1, '2': 0x2, '4': 0x4,
		'8': 0x8, '3': 0x3, '5': 0x5, '6': 0x6,
		'9': 0x9, 'a': 0xa, 'c': 0xc, '7': 0x7,
		'b': 0xb, 'd': 0xd, 'e': 0xe, 'f': 0xf,
	}
)

var (
	errInvalidMap = errors.New("invalid cpu map")
)

// read cpu map in hex format as a first line from b and return lcore
// ids and error if encountered.
func readCpuHexMap(b io.Reader) ([]int, error) {
	cores := []int{}
	scanner := bufio.NewScanner(b)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	cpumap := scanner.Bytes()
	for i := range cpumap {
		c := cpumap[len(cpumap)-i-1]
		n, ok := hexMap[c]
		if !ok {
			return nil, errInvalidMap
		}
		if n&0x1 != 0 {
			cores = append(cores, 4*i)
		}
		if n&0x2 != 0 {
			cores = append(cores, 4*i+1)
		}
		if n&0x4 != 0 {
			cores = append(cores, 4*i+2)
		}
		if n&0x8 != 0 {
			cores = append(cores, 4*i+3)
		}
	}

	return cores, nil
}

// NumaNode returns id of the NUMA node for specified logical CPU
// core id. if core id is invalid, NumaNodeAny is returned.
func NumaNode(id uint) int {
	return getNumaNode(id)
}
