package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CommandType int

var _ = fmt.Print

const (
	CmdUnknown CommandType = iota
	CmdExit
	CmdEcho
	CmdType
)

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

	command := parts[0]
	cmdType := parseCommand(command)

	switch cmdType {
	case CmdExit:
		return true
	case CmdEcho:
		handleEcho(parts[1:])
	case CmdType:
		handleType(parts[1:])
	default:
		fmt.Printf("%s: command not found\n", command)
		return false
	}
	return false
}

func handleType(args []string) {
	cmdStr := strings.Join(args, " ")

	typeOf := parseCommand(cmdStr)
	switch typeOf {
	case CmdUnknown:
		fmt.Printf("%s: not found\n", cmdStr)
	default:
		fmt.Printf("%s is a shell builtin\n", cmdStr)
	}
}

func handleEcho(args []string) {
	output := strings.Join(args, " ")
	fmt.Printf("%s\n", output)
}
