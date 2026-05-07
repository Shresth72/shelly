package commands

import (
	"fmt"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

func handleUnknown(ctx CommandContext) {
	path, executable := utils.FindExecutable(ctx.Name)
	if path == "" {
		fmt.Fprintf(ctx.Stderr, "%s: command not found\n", ctx.Name)
		return
	}

	if !executable {
		fmt.Fprintln(ctx.Stderr, "Permission Denied")
		return
	}
	runExecutable(ctx)
}

func runExecutable(ctx CommandContext) {
	cmd := exec.Command(ctx.Name, ctx.Args...)

	cmd.Stdin = ctx.Stdin
	cmd.Stdout = ctx.Stdout
	cmd.Stderr = ctx.Stderr

	err := cmd.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return
		}
		fmt.Fprintln(ctx.Stderr, "Error:", err)
	}
}
