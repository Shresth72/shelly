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
	var result []string
	var current strings.Builder

	state := none

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '\'':
			if state == none {
				state = single
			} else if state == single {
				state = none
			} else {
				current.WriteByte(ch)
			}

		case '"':
			if state == none {
				state = double
			} else if state == double {
				state = none
			} else {
				current.WriteByte(ch)
			}

		case ' ':
			if state == none {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(ch)
			}

		// case '\\':
		// 	if state == double {
		// 		if state != escape {
		// 			state = escape
		// 		} else if state
		// 	}

		default:
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
