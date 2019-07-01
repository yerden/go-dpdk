package common

import (
	"bufio"
	"bytes"
	"errors"
	"unicode"
	"unicode/utf8"
)

var (
	ErrUnprintable  = errors.New("unprintable char")
	ErrOpenQuote    = errors.New("no closing quote")
	DefaultSplitter = &Splitter{
		unicode.IsSpace,
		func(r rune) (rune, bool) {
			if r == '"' {
				return '"', true
			}
			if r == '\'' {
				return '\'', true
			}
			return ' ', false
		},
		false,
	}
)

type Splitter struct {
	// True if rune is a white space.
	IsSpace func(rune) bool

	// True if rune is quote. Any rune embraced by the one of these
	// pairs is considered a part of a token even if IsSpace returns
	// true.  A pairs must not contradict white space and another
	// pair.
	//
	// If true, return closing quote rune.
	IsQuote func(rune) (rune, bool)

	// If true, final token is allowed not to contain closing quote.
	// If false, ErrOpenQuote error will be returned if no closing
	// quote found.
	AllowOpenQuote bool
}

func SplitFunc(s *Splitter) bufio.SplitFunc {
	isSpaceOrQuote := func(r rune) bool {
		_, ok := s.IsQuote(r)
		return ok || s.IsSpace(r)
	}

	isNotSpace := func(r rune) bool {
		return !s.IsSpace(r)
	}

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// skip whitespace
		n := bytes.IndexFunc(data, isNotSpace)
		if n < 0 {
			return len(data), nil, nil
		}

		// start of a token
		quoted := false
		advance = n
		data = data[n:]
		n = 0

		for {
			k := bytes.IndexFunc(data[n:], isSpaceOrQuote)
			if k < 0 {
				// unterminated token, not enough data
				break
			}

			// rune 'r' is either space or quote
			r, wid := utf8.DecodeRune(data[n+k:])

			// if quote, then look for closing quote
			if q, ok := s.IsQuote(r); ok {
				n += k + wid
				if k = bytes.IndexRune(data[n:], q); k < 0 {
					quoted = true
					break
				}
				n += k + wid
				continue
			}

			// 'r' is white space
			return advance + n + k + wid, data[:n+k], nil
		}

		if !atEOF {
			// unterminated token, need more data
			return advance, nil, nil
		} else if !quoted || s.AllowOpenQuote {
			// unterminated token, no more data
			return advance + len(data), data, nil
		} else {
			// unterminated quote is not allowed
			return advance + len(data), data, ErrOpenQuote
		}
	}
}
