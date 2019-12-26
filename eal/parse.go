package eal

import (
	"bufio"
	"strings"

	"github.com/yerden/go-dpdk/common"
)

func parseCmd(input string) ([]string, error) {
	s := bufio.NewScanner(strings.NewReader(input))
	s.Split(common.SplitFunc(common.DefaultSplitter))

	var argv []string
	for s.Scan() {
		argv = append(argv, s.Text())
	}
	return argv, s.Err()
}

// InitCmd initializes EAL as in rte_eal_init. Options are specified
// in a unparsed command line string. This string is parsed and
// Init is then called upon.
func InitCmd(input string) (int, error) {
	argv, err := parseCmd(input)
	if err != nil {
		return 0, err
	}
	return Init(argv)
}
