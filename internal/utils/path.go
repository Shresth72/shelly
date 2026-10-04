package utils

import (
	"os"
	"strings"
	"sync"
)

func hasExecutePermission(info os.FileInfo) bool {
	return info.Mode().Perm()&0111 != 0
}

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

func executableNames(dir, prefix string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var commands []string

	for _, entry := range entries {
		name := entry.Name()

		if entry.IsDir() || !strings.HasPrefix(name, prefix) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if hasExecutePermission(info) {
			commands = append(commands, name)
		}
	}

	return commands
}

func CompleteExecutables(prefix string) []string {
	commands := make(map[string]struct{})

	for dir := range strings.SplitSeq(os.Getenv("PATH"), ":") {
		if dir == "" {
			dir = "."
		}

		for _, name := range executableNames(dir, prefix) {
			commands[name] = struct{}{}
		}
	}

	return mapKeys(commands)
}

func Executables() []string {
	var wg sync.WaitGroup
	var mu sync.Mutex

	commands := make(map[string]struct{})

	for dir := range strings.SplitSeq(os.Getenv("PATH"), ":") {
		if dir == "" {
			dir = "."
		}

		wg.Add(1)

		go func(dir string) {
			defer wg.Done()

			found := executableNames(dir, "")

			mu.Lock()
			for _, name := range found {
				commands[name] = struct{}{}
			}
			mu.Unlock()
		}(dir)
	}

	wg.Wait()

	return mapKeys(commands)
}

func mapKeys(set map[string]struct{}) []string {
	result := make([]string, 0, len(set))

	for name := range set {
		result = append(result, name)
	}

	return result
}
