package shelly

import (
	"fmt"
	"sync"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

type BellCompleter struct {
	mu sync.RWMutex

	completer *readline.PrefixCompleter
	ready     bool
}

func NewAutoCompleter() *BellCompleter {
	c := &BellCompleter{}

	c.setCompleter(buildCommandCompleter(nil))

	go func() {
		executables := utils.Executables()

		c.setCompleter(buildCommandCompleter(executables))

		c.mu.Lock()
		c.ready = true
		c.mu.Unlock()
	}()

	return c
}

func (c *BellCompleter) setCompleter(completer *readline.PrefixCompleter) {
	c.mu.Lock()
	c.completer = completer
	c.mu.Unlock()
}

func (c *BellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	c.mu.RLock()
	completer := c.completer
	ready := c.ready
	c.mu.RUnlock()

	if !ready {
		prefix := string(line[:pos])
		executables := utils.CompleteExecutables(prefix)
		completer = buildCommandCompleter(executables)
	}

	newLine, length := completer.Do(line, pos)

	if len(newLine) == 0 {
		fmt.Print("\x07")
	}

	return newLine, length
}

func buildCommandCompleter(executables []string) *readline.PrefixCompleter {
	commands := make(map[string]struct{}, len(builtinCommands)+len(executables))
	for _, command := range builtinCommands {
		commands[command] = struct{}{}
	}
	for _, command := range executables {
		commands[command] = struct{}{}
	}

	items := make([]readline.PrefixCompleterInterface, 0, len(commands))
	for command := range commands {
		items = append(items, readline.PcItem(command))
	}

	return readline.NewPrefixCompleter(items...)
}

func buildCompleter(commands []string) *readline.PrefixCompleter {
	items := make([]readline.PrefixCompleterInterface, 0, len(commands))
	for _, command := range commands {
		items = append(items, readline.PcItem(command))
	}

	return readline.NewPrefixCompleter(items...)
}
