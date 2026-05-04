package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var _ = fmt.Print

func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		commands, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading line:", err)
			os.Exit(1)
		}

		if handleCommands(commands) {
			break
		}
	}

}

func handleCommands(commands string) bool {
	commands = strings.TrimSpace(commands)
	commandsArr := strings.Split(commands, " ")

	command := commandsArr[0]

	switch command {
	case "exit":
		return true
	case "echo":
		handleEcho(commandsArr)
	default:
		fmt.Printf("%s: command not found\n", command)
		return false
	}

	return false
}

func handleEcho(commands []string) {
	cmdStr := strings.Join(commands[1:], " ")
	fmt.Printf("%s\n", cmdStr)
}
