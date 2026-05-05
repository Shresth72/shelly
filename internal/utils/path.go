package utils

import (
	"os"
	"strings"
)

func FindExecutable(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", false
	}

	for dir := range strings.SplitSeq(pathEnv, ":") {
		if dir == "" {
			dir = "."
		}

		fullPath := dir + "/" + command
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}

		if hasExecutePermission(info) {
			return fullPath, true
		}
	}

	return "", false
}

func hasExecutePermission(info os.FileInfo) bool {
	return info.Mode().Perm()&0111 != 0
}
