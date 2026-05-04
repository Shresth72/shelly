package main

import (
	"bufio"
	"fmt"
	"os"
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
		fmt.Printf("%s: command not found\n", command[:len(command)-1])
	}

}
