package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

func handleExit(ctx CommandContext) bool {
	return true
}

func handleEcho(ctx CommandContext) bool {
	fmt.Println(ctx.CmdStr)
	return false
}

func handlePwd(ctx CommandContext) bool {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}
	fmt.Println(dir)
	return false
}

func handleType(ctx CommandContext) bool {
	if len(ctx.Args) == 0 {
		return false
	}

	for _, target := range ctx.Args {
		// Builtin
		if _, ok := commands[target]; ok {
			fmt.Printf("%s is a shell builtin\n", target)
			continue
		}

		// External
		path, _ := utils.FindExecutable(target)
		if path == "" {
			fmt.Printf("%s: not found\n", target)
			continue
		}

		fmt.Printf("%s is %s\n", target, path)
	}

	return false
}

func handleCd(ctx CommandContext) bool {
	if len(ctx.Args) > 1 {
		fmt.Println("cd: too many arguments")
		return false
	}

	var target string

	if len(ctx.Args) == 0 {
		target = os.Getenv("HOME")
		if target == "" {
			fmt.Println("cd: HOME not set")
			return false
		}
	} else {
		target = ctx.Args[0]
	}

	if target == "~" || strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("cd: unable to determine home directory")
			return false
		}

		target = strings.Replace(target, "~", home, 1)
	}

	err := os.Chdir(target)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", target)
	}

	return false
}
