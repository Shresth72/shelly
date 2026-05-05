package utils

import "strings"

func Tokenize(input string) []string {
	var result []string
	var current strings.Builder

	inSingle := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch ch {
		case '\'':
			inSingle = !inSingle

		case ' ':
			if inSingle {
				current.WriteByte(ch)
			} else {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
			}

		default:
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
