package utils

import "strings"

type quoteState int

const (
	none quoteState = iota
	single
	double
	escape
)

func Tokenize(input string) []string {
	var tokens []string
	var current strings.Builder

	var quote byte
	escaped := false

	flush := func() {
		if current.Len() > 0 {
			tokens = append(tokens, current.String())
			current.Reset()
		}
	}

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}

		switch ch {
		case '\\':
			if quote == '\'' {
				current.WriteByte(ch)
			} else {
				escaped = true
			}

		case '"', '\'':
			if quote == 0 {
				quote = ch
			} else if quote == ch {
				quote = 0
			} else {
				current.WriteByte(ch)
			}

		case ' ', '\t', '\n':
			if quote != 0 {
				current.WriteByte(ch)
			} else {
				flush()
			}

		default:
			current.WriteByte(ch)
		}
	}

	flush()

	return tokens
}
