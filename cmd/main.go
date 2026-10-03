package main

import (
	"fmt"
	"io"
	"os"

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

	for {
		input, err := rl.Readline()

		if err == readline.ErrInterrupt {
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading line:", err)
			break
		}
		if input == "" {
			continue
		}

		if shelly.HandleInput(input, &shell) {
			break
		}
	}
}
