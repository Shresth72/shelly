package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/commands"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading line:", err)
			os.Exit(1)
		}

		if commands.HandleCommands(input, os.Stdin, os.Stdout, os.Stderr) {
			break
		}
	}
}
