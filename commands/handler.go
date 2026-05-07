package commands

import (
	"fmt"
	"io"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

type CommandContext struct {
	// Command
	Name   string
	Args   []string
	CmdStr string

	// Context
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Pipe - Type of Node with just Left/Right (Node type as well)

type CommandHandler func(ctx CommandContext) bool

var commands map[string]CommandHandler

func init() {
	commands = map[string]CommandHandler{
		"echo": handleEcho,
		"exit": handleExit,
		"type": handleType,
		"pwd":  handlePwd,
		"cd":   handleCd,
	}
}

// TODO: Pass Reader/Writer to handler
func HandleCommands(
	input string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) bool {
	input = strings.TrimSpace(input)

	parts := utils.Tokenize(input)
	if len(parts) == 0 {
		return false
	}

	redirects, err := parseRedirects(parts, stdout, stderr)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return false
	}
	defer redirects.Close()

	parts = redirects.Args

	ctx := CommandContext{
		Name:   parts[0],
		Args:   parts[1:],
		CmdStr: strings.Join(parts[1:], " "),

		Stdin:  stdin,
		Stdout: redirects.Stdout,
		Stderr: redirects.Stderr,
	}

	if handler, ok := commands[ctx.Name]; ok {
		return handler(ctx)
	}

	handleUnknown(ctx)
	return false
}
