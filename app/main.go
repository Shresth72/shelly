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

		handleCommands(command)
	}

}

func handleCommands(command string) {
	command = strings.TrimSpace(command)

	switch command {
	case "exit":
		os.Exit(1)
	default:
		fmt.Printf("%s: command not found\n", command)
	}

}
