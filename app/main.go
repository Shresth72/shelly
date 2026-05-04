package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
		handleEchoCmd(ctx)
	case CmdType:
		handleTypeCmd(ctx)
	default:
		handleUnknownCmd(ctx)
	}

	return false
}

func handleEchoCmd(ctx CommandContext) {
	fmt.Println(ctx.CmdStr)
}

func handleTypeCmd(ctx CommandContext) {
	if len(ctx.Args) == 0 {
		return
	}

	target := ctx.Args[0]

	typeOf := parseCommand(target)
	switch typeOf {
	case CmdUnknown:
		handleUnknownType(target)
	default:
		fmt.Printf("%s is a shell builtin\n", target)
	}
}

func handleUnknownCmd(ctx CommandContext) {
	path, executable := validateExecutable(ctx.Name)

	if path == "" {
		fmt.Printf("%s: command not found\n", ctx.Name)
		return
	}

	if !executable {
		fmt.Println("Permission Denied")
		return
	}
	runExecutable(path, ctx)
}

func handleUnknownType(cmd string) {
	path, _ := validateExecutable(cmd)

	if path == "" {
		fmt.Printf("%s: not found\n", cmd)
		return
	}

	fmt.Printf("%s is %s\n", cmd, path)
}

func validateExecutable(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", false
	}

	for dir := range strings.SplitSeq(pathEnv, ":") {
		if dir == "" {
			dir = "."
		}

		fullPath := dir + "/" + command

		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if hasExecutePermission(info) {
			return fullPath, true
		}
	}

	return "", false
}

// func validateExecutable(command string) (string, bool) {
// 	path, err := exec.LookPath(command)
// 	if err != nil {
// 		return "", false
// 	}
//
// 	info, err := os.Stat(path)
// 	if err != nil {
// 		return "", false
// 	}
//
// 	return path, hasExecutePermission(info)
// }

func hasExecutePermission(info os.FileInfo) bool {
	return info.Mode().Perm()&0111 != 0
}

func runExecutable(exe string, ctx CommandContext) {
	cmd := exec.Command(exe, ctx.Args...)

	cmd.Args = append([]string{ctx.Name}, ctx.Args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
