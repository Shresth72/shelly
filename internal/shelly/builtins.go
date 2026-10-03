package shelly

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

func writef(w io.Writer, format string, args ...any) error {
	_, err := fmt.Fprintf(w, format, args...)
	return err
}

func exitCmd(cmd *Command, sh *Shell) error {
	return ErrExit
}

func echoCmd(cmd *Command, sh *Shell) error {
	_, err := fmt.Fprintln(
		sh.Stdout,
		strings.Join(cmd.Args, " "),
	)
	return err
}

func pwdCmd(cmd *Command, sh *Shell) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	_, err = fmt.Fprintln(sh.Stdout, dir)
	return err
}

func typeCmd(cmd *Command, sh *Shell) error {
	for _, target := range cmd.Args {
		if isBuiltin(target) {
			if err := writef(sh.Stdout, "%s is a shell builtin\n", target); err != nil {
				return err
			}
			continue
		}

		path, _ := utils.FindExecutable(target)

		if path == "" {
			if err := writef(sh.Stdout, "%s: not found\n", target); err != nil {
				return err
			}
			continue
		}

		if err := writef(sh.Stdout, "%s is %s\n", target, path); err != nil {
			return err
		}
	}

	return nil
}

func cdCmd(cmd *Command, sh *Shell) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("cd: too many arguments")
	}

	var target string

	if len(cmd.Args) == 0 {
		target = os.Getenv("HOME")
		if target == "" {
			return fmt.Errorf("cd: HOME not set")
		}
	} else {
		target = cmd.Args[0]
	}

	if target == "~" || strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cd: unable to determine home directory")
		}

		target = strings.Replace(target, "~", home, 1)
	}

	err := os.Chdir(target)
	if err != nil {
		return fmt.Errorf("cd: %s: No such file or directory", target)
	}

	return nil
}
