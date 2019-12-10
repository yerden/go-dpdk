package eal

import (
	"fmt"
	"strings"
	"unicode"
)

// Parameter specifies a command line option-argument pair. If
// keyword Opt doesn't imply argument value, Arg should be nil.
type Parameter struct {
	Opt string
	Arg []interface{}
}

// Set mutates Parameter setting new value to option in a form of
// array of values.
func (p Parameter) Set(a ...interface{}) Parameter {
	p.Arg = append(make([]interface{}, 0, len(a)), a...)
	return p
}

// Join creates command line out of parameters.
func Join(params []Parameter) []string {
	argv := make([]string, 0, 2*len(params))
	for _, p := range params {
		if argv = append(argv, p.Opt); len(p.Arg) == 0 {
			continue
		} else if arg, ok := p.Arg[0].(string); ok {
			argv = append(argv, fmt.Sprintf(arg, p.Arg[1:]...))
		} else {
			argv = append(argv, fmt.Sprint(p.Arg...))
		}
	}
	return argv
}

func isBadChar(r rune) bool {
	return unicode.IsSpace(r) || !unicode.IsGraphic(r)
}

// NewParameter creates new Parameter with opt as a keyword and an
// optional array of values which constitute opt's value.
//
// Panic will be emitted if opt contains whitespace.
func NewParameter(opt string, a ...interface{}) Parameter {
	if strings.IndexFunc(opt, isBadChar) < 0 {
		return Parameter{Opt: opt}.Set(a...)
	}
	panic("invalid option '" + opt + "'")
}
