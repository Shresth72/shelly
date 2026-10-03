package shelly

import (
	"errors"
	"os"
)

func ExecuteAST(
	node Node,
	sh *Shell,
) (bool, error) {
	switch n := node.(type) {
	case *Command:
		return executeCommand(n, sh)
	case *Pipe:
		return executePipe(n, sh)
	default:
		return false, errors.New("Unknown AST node")
	}
}

func executePipe(pipe *Pipe, sh *Shell) (bool, error) {
	reader, writer, err := os.Pipe()
	if err != nil {
		return false, err
	}

	leftShell := *sh
	leftShell.Stdout = writer

	rightShell := *sh
	rightShell.Stdout = reader

	type result struct {
		exit bool
		err  error
	}

	leftResult := make(chan result, 1)
	rightResult := make(chan result, 1)

	go func() {
		exit, err := ExecuteAST(pipe.Left, &leftShell)
		_ = writer.Close()
		leftResult <- result{exit: exit, err: err}
	}()

	go func() {
		exit, err := ExecuteAST(pipe.Left, &rightShell)
		_ = writer.Close()
		rightResult <- result{exit: exit, err: err}
	}()

	left := <-leftResult
	right := <-rightResult

	if left.err != nil && !errors.Is(left.err, ErrExit) {
		return false, left.err
	}
	if right.err != nil && !errors.Is(right.err, ErrExit) {
		return false, right.err
	}

	return false, nil
}
