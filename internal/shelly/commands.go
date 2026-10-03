package shelly

import (
	"errors"
)

var ErrExit = errors.New("shell exit")

var builtinCommands = []string{
	"echo",
	"pwd",
	"type",
	"cd",
	"exit",
}

func isBuiltin(name string) bool {
	for _, cmd := range builtinCommands {
		if cmd == name {
			return true
		}
	}
	return false
}

func executeCommand(
	cmd *Command,
	sh *Shell,
) (bool, error) {
	ctx, cleanup, err := applyRedirects(cmd.Redirs, sh)
	if err != nil {
		return false, err
	}
	defer cleanup()

	switch cmd.Name {
	case "echo":
		err = echoCmd(cmd, ctx)
	case "pwd":
		err = pwdCmd(cmd, ctx)
	case "type":
		err = typeCmd(cmd, ctx)
	case "cd":
		err = cdCmd(cmd, ctx)
	case "exit":
		err = exitCmd(cmd, ctx)
	default:
		err = executeExternal(cmd, ctx)
	}

	if errors.Is(err, ErrExit) {
		return true, ErrExit
	}

	return false, err
}
