package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CommandType int

const (
	CmdUnknown CommandType = iota
	CmdExit
	CmdEcho
	CmdType
)

type CommandContext struct {
	Name   string
	Args   []string
	CmdStr string
}

func (c CommandContext) isBuiltin() bool {
	return parseCommand(c.Name) != CmdUnknown
}

func parseCommand(cmd string) CommandType {
	switch cmd {
	case "exit":
		return CmdExit
	case "echo":
		return CmdEcho
	case "type":
		return CmdType
	default:
		return CmdUnknown
	}
}

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading line:", err)
			os.Exit(1)
		}

		if handleCommands(input) {
			break
		}
	}

}

func handleCommands(commands string) bool {
	input := strings.TrimSpace(commands)

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	ctx := CommandContext{
		Name:   parts[0],
		Args:   parts[1:],
		CmdStr: strings.Join(parts[1:], " "),
	}
	cmdType := parseCommand(ctx.Name)

	switch cmdType {
	case CmdExit:
		return true
	case CmdEcho:
		handleEcho(ctx)
	case CmdType:
		handleType(ctx)
	default:
		fmt.Printf("%s: command not found\n", ctx.Name)
		return false
	}
	return false
}

func handleEcho(ctx CommandContext) {
	fmt.Println(ctx.CmdStr)
}

func handleType(ctx CommandContext) {
	target := ctx.CmdStr

	typeOf := parseCommand(target)
	switch typeOf {
	case CmdUnknown:
		handleUnknownTypeOf(ctx)
	default:
		fmt.Printf("%s is a shell builtin\n", target)
	}
}

func handleUnknownTypeOf(ctx CommandContext) {
	if validatePathExecutable(ctx) {
		return
	}
	fmt.Printf("%s: not found\n", ctx.CmdStr)
}

func validatePathExecutable(ctx CommandContext) bool {
	path := os.Getenv("PATH")
	if len(path) == 0 {
		return false
	}

	for dir := range strings.SplitSeq(path, ":") {
		fullPath := dir + "/" + ctx.CmdStr
		info, err := os.Stat(fullPath)
		if err == nil && info.Mode().Perm()&0111 != 0 {
			fmt.Printf("%s is %s\n", ctx.CmdStr, fullPath)
			return true
		}
	}
	return false
}
