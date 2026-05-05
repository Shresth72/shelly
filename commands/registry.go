package commands

import (
	"strings"
)

type CommandContext struct {
	Name   string
	Args   []string
	CmdStr string
}

type CommandHandler func(ctx CommandContext) bool

var commands map[string]CommandHandler

func init() {
	commands = map[string]CommandHandler{
		"echo": handleEcho,
		"exit": handleExit,
		"type": handleType,
		"pwd":  handlePwd,
	}
}

// TODO: Pass Reader/Writer to handler
func HandleCommands(input string) bool {
	input = strings.TrimSpace(input)

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	ctx := CommandContext{
		Name:   parts[0],
		Args:   parts[1:],
		CmdStr: strings.Join(parts[1:], " "),
	}

	if handler, ok := commands[ctx.Name]; ok {
		return handler(ctx)
	}

	handleUnknown(ctx)
	return false
}
