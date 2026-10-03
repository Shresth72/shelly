package shelly

import (
	"fmt"
	"io"
	"strings"
)

type Shell struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func HandleInput(
	input string,
	sh *Shell,
) bool {
	input = strings.TrimSpace(input)

	parts := Tokenize(input)
	if len(parts) == 0 {
		return false
	}

	ast, err := ParseAST(parts)
	if err != nil {
		fmt.Fprintln(sh.Stderr, err)
		return false
	}

	exit, err := ExecuteAST(ast, sh)
	if err != nil {
		fmt.Fprintln(sh.Stderr, err)
		return false
	}

	return exit
}
