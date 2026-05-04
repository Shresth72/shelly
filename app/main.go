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
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading line:", err)
			os.Exit(1)
		}

		if handleCommands(command) {
			break
		}
	}

}

func handleCommands(command string) bool {
	command = strings.TrimSpace(command)

	switch command {
	case "exit":
		return true
	default:
		fmt.Printf("%s: command not found\n", command)
		return false
	}
}
