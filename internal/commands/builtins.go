package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

// TODO: Pass Command and make these just return what they want to print
// And try to achieve only Command as the thing have access to
func handleExit(ctx CommandContext) bool {
	return true
}

func handleEcho(ctx CommandContext) bool {
	fmt.Fprintln(ctx.Stdout, ctx.CmdStr)
	return false
}

func handlePwd(ctx CommandContext) bool {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(ctx.Stderr, "Error:", err)
		return false
	}
	fmt.Fprintln(ctx.Stdout, dir)
	return false
}

func handleType(ctx CommandContext) bool {
	if len(ctx.Args) == 0 {
		return false
	}

	for _, target := range ctx.Args {
		if _, ok := commands[target]; ok {
			fmt.Fprintf(ctx.Stdout, "%s is a shell builtin\n", target)
			continue
		}

		path, _ := utils.FindExecutable(target)
		if path == "" {
			fmt.Fprintf(ctx.Stdout, "%s: not found\n", target)
			continue
		}

		fmt.Fprintf(ctx.Stdout, "%s is %s\n", target, path)
	}

	return false
}

func handleCd(ctx CommandContext) bool {
	if len(ctx.Args) > 1 {
		fmt.Fprintln(ctx.Stderr, "cd: too many arguments")
		return false
	}

	var target string

	if len(ctx.Args) == 0 {
		target = os.Getenv("HOME")
		if target == "" {
			fmt.Fprintln(ctx.Stderr, "cd: HOME not set")
			return false
		}
	} else {
		target = ctx.Args[0]
	}

	if target == "~" || strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(ctx.Stderr, "cd: unable to determine home directory")
			return false
		}

		target = strings.Replace(target, "~", home, 1)
	}

	err := os.Chdir(target)
	if err != nil {
		fmt.Fprintf(ctx.Stderr, "cd: %s: No such file or directory\n", target)
	}

	return false
}
