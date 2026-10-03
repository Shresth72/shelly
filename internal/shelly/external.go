package shelly

import (
	"fmt"
	"os/exec"

	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

func executeExternal(cmd *Command, sh *Shell) error {
	path, executable := utils.FindExecutable(cmd.Name)
	if path == "" {
		return fmt.Errorf("%s: Command not found", cmd.Name)
	}

	if !executable {
		return fmt.Errorf("%s: Permission denied", cmd.Name)
	}

	proc := exec.Command(path, cmd.Args...)

	proc.Stdin = sh.Stdin
	proc.Stdout = sh.Stdout
	proc.Stderr = sh.Stderr

	err := proc.Run()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil
		}
		return err
	}

	return nil
}
