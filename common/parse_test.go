package common_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/yerden/go-dpdk/common"
)

func TestParseSplitDoubleQuotes(t *testing.T) {
	assert := common.Assert(t, true)

	input := "-c f    -n 4 --socket-mem=\"1024, 1024\" --no-pci -d /path/to/so"
	b := bytes.NewBufferString(input)
	s := bufio.NewScanner(b)
	s.Split(common.SplitFunc(common.DefaultSplitter))
	var argv []string

	for s.Scan() {
		argv = append(argv, s.Text())
	}

	assert(len(argv) == 8, argv, len(argv))
	assert(argv[0] == "-c", argv)
	assert(argv[1] == "f", argv)
	assert(argv[2] == "-n", argv)
	assert(argv[3] == "4", argv)
	assert(argv[4] == "--socket-mem=\"1024, 1024\"", argv)
	assert(argv[5] == "--no-pci", argv)
	assert(argv[6] == "-d", argv)
	assert(argv[7] == "/path/to/so", argv)
}
