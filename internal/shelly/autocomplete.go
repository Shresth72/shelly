package shelly

import (
	"fmt"
	"sort"
	"sync"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/internal/utils"
)

type BellCompleter struct {
	mu sync.RWMutex

	completer *readline.PrefixCompleter
	ready     bool

	belled bool
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
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.belled {
		c.belled = true

		fmt.Print("\x07")
		return nil, 0
	}
	c.belled = false

	completer := c.completer

	if !c.ready {
		prefix := string(line[:pos])
		executables := utils.CompleteExecutables(prefix)
		completer = buildCommandCompleter(executables)
	}

	return completer.Do(line, pos)
}

func (c *BellCompleter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.belled = false
}

func buildCommandCompleter(executables []string) *readline.PrefixCompleter {
	commands := make(map[string]struct{}, len(builtinCommands)+len(executables))

	for _, command := range builtinCommands {
		commands[command] = struct{}{}
	}

	for _, command := range executables {
		commands[command] = struct{}{}
	}

	sortedCommands := make([]string, 0, len(commands))
	for command := range commands {
		sortedCommands = append(sortedCommands, command)
	}

	sort.Strings(sortedCommands)

	items := make([]readline.PrefixCompleterInterface, 0, len(sortedCommands))
	for _, command := range sortedCommands {
		items = append(items, readline.PcItem(command))
	}

	return readline.NewPrefixCompleter(items...)
}
