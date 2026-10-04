package shelly

import (
	"fmt"
	"os"
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

type InputListener struct {
	readline.Listener
}

func (l *InputListener) OnChange(line []rune, pos int, key rune) ([]rune, int, bool) {
	if key == readline.CharDelete {
		fmt.Println()
		os.Exit(0)
	}

	return line, pos, false
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

func (c *BellCompleter) InputFilter(r rune) (rune, bool) {
	return r, true
}

func (c *BellCompleter) InputListener(line []rune, pos int, key rune) ([]rune, int, bool) {
	if key == 4 {
		fmt.Println()
		return nil, 0, true
	}
	return line, pos, false
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
