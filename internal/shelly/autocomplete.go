package shelly

import (
	"fmt"

	"github.com/chzyer/readline"
)

type BellCompleter struct {
	*readline.PrefixCompleter
}

func (c *BellCompleter) Do(line []rune, pos int) ([][]rune, int) {
	newLine, length := c.PrefixCompleter.Do(line, pos)

	if len(newLine) == 0 {
		fmt.Print("\a")
	}

	return newLine, length
}

func NewAutoCompleter() *BellCompleter {
	items := make([]readline.PrefixCompleterInterface, 0, len(allCommands))

	for _, command := range allCommands {
		items = append(items, readline.PcItem(command))
	}

	return &BellCompleter{
		PrefixCompleter: readline.NewPrefixCompleter(items...),
	}
}
