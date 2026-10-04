package shelly

import (
	"errors"
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
) (bool, error) {
	input = strings.TrimSpace(input)

	parts, err := Tokenize(input)
	if err != nil {
		return false, err
	}
	if len(parts) == 0 {
		return false, nil
	}

	ast, err := ParseAST(parts)
	if err != nil {
		fmt.Fprintln(sh.Stderr, err)
		return false, err
	}

	exit, err := ExecuteAST(ast, sh)
	if errors.Is(err, ErrExit) {
		return true, nil
	}
	if err != nil {
		fmt.Fprintln(sh.Stderr, err)
		return false, err
	}

	return exit, nil
}
