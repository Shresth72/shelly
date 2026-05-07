package shelly

import (
	"fmt"
	"io"
	"strings"
)

type Context struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func HandleInput(
	input string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) bool {
	input = strings.TrimSpace(input)

	parts := Tokenize(input)
	if len(parts) == 0 {
		return false
	}

	ast, err := ParseAST(parts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return false
	}

	exit, err := ExecuteAST(ast)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return false
	}

	return exit
}
