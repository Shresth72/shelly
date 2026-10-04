package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/internal/shelly"
)

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		AutoComplete:    shelly.NewAutoCompleter(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		HistoryFile:     "/tmp/shelly_history",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating readline:", err)
		os.Exit(1)
	}
	defer rl.Close()

	shell := shelly.Shell{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	var input strings.Builder

	for {
		if input.Len() > 0 {
			rl.SetPrompt("> ")
		} else {
			rl.SetPrompt("$ ")
		}

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			input.Reset()
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(shell.Stderr, "Error reading line:", err)
			break
		}

		input.WriteString(line)

		exit, err := shelly.HandleInput(input.String(), &shell)
		if errors.Is(err, shelly.ErrIncomplete) {
			continue
		}

		input.Reset()

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if exit {
			break
		}
	}
}
