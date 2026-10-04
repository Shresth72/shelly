package shelly

import (
	"errors"
	"strings"
)

type quoteState int

var ErrIncomplete = errors.New("shell incomplete")

const (
	none quoteState = iota
	single
	double
)

func Tokenize(input string) ([]string, error) {
	var tokens []string
	var current strings.Builder

	var quote quoteState
	escaped := false
	tokenStarted := false

	flush := func() {
		if tokenStarted {
			tokens = append(tokens, current.String())
			current.Reset()
			tokenStarted = false
		}
	}

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			tokenStarted = true
			continue
		}

		if ch == '\\' && quote != single {
			escaped = true
			tokenStarted = true
			continue
		}

		if ch == '\'' || ch == '"' {
			if quote == none {
				if ch == '\'' {
					quote = single
				} else {
					quote = double
				}
				tokenStarted = true
			} else if (quote == single && ch == '\'') || (quote == double && ch == '"') {
				quote = none
			} else {
				current.WriteByte(ch)
			}
			continue
		}

		if quote == none {
			switch ch {
			case ' ', '\t', '\n':
				flush()
				continue

			case '|', '<', '>':
				flush()
				if ch == '>' && i+1 < len(input) && input[i+1] == '>' {
					tokens = append(tokens, ">>")
					i++
				} else {
					tokens = append(tokens, string(ch))
				}
				continue
			}
		}

		current.WriteByte(ch)
		tokenStarted = true
	}

	if escaped || quote != none {
		return nil, ErrIncomplete
	}

	flush()
	return tokens, nil
}
