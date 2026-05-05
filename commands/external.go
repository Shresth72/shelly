package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

func handleUnknown(ctx CommandContext) {
	path, executable := utils.FindExecutable(ctx.Name)
	if path == "" {
		fmt.Printf("%s: command not found\n", ctx.Name)
		return
	}

	if !executable {
		fmt.Println("Permission Denied")
		return
	}
	runExecutable(ctx)
}

func runExecutable(ctx CommandContext) {
	cmd := exec.Command(ctx.Name, ctx.Args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
