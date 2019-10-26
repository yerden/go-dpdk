package eal

import (
	"fmt"
	"strings"
	"unicode"
)

// Parameter specifies a command line option-argument pair. If
// keyword Opt doesn't imply argument value, Arg should be "".
type Parameter struct {
	Opt, Arg string
}

// Set mutates Parameter setting new value to option in a form
// array of values.
func (p Parameter) Set(a ...interface{}) Parameter {
	if len(a) > 0 {
		if arg, ok := a[0].(string); ok {
			p.Arg = fmt.Sprintf(arg, a[1:]...)
		} else {
			p.Arg = fmt.Sprint(a...)
		}
	}
	return p
}

// Join creates command line out of parameters.
func Join(params []Parameter) []string {
	argv := make([]string, 0, 2*len(params))
	for _, p := range params {
		if argv = append(argv, p.Opt); p.Arg != "" {
			argv = append(argv, p.Arg)
		}
	}
	return argv
}

// NewParameter creates new Parameter with opt as a keyword and an
// optional array of values which constitute opt's value.
//
// Panic will be emitted if opt contains whitespace.
func NewParameter(opt string, a ...interface{}) Parameter {
	if strings.IndexFunc(opt, unicode.IsSpace) < 0 {
		return Parameter{Opt: opt}.Set(a...)
	}
	panic("invalid option '" + opt + "'")
}
