package utils

import (
	"os"
	"strings"
	"sync"
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

func Executables() []string {
	pathEnv := os.Getenv("PATH")

	var wg sync.WaitGroup
	var mu sync.Mutex

	var commands []string

	for dir := range strings.SplitSeq(pathEnv, ":") {
		if dir == "" {
			dir = "."
		}

		wg.Add(1)

		go func(dir string) {
			defer wg.Done()

			entries, err := os.ReadDir(dir)
			if err != nil {
				return
			}

			var found []string

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				info, err := entry.Info()
				if err != nil {
					continue
				}

				if hasExecutePermission(info) {
					found = append(found, entry.Name())
				}
			}

			mu.Lock()
			commands = append(commands, found...)
			mu.Unlock()
		}(dir)
	}

	wg.Wait()

	return commands
}

func hasExecutePermission(info os.FileInfo) bool {
	return info.Mode().Perm()&0111 != 0
}
