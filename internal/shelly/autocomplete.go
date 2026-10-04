package shelly

import (
	"fmt"

	"github.com/chzyer/readline"
	// "github.com/codecrafters-io/shell-starter-go/internal/utils"
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
	items := make([]readline.PrefixCompleterInterface, 0)

	for _, command := range builtinCommands {
		items = append(items, readline.PcItem(command))
	}

	// for _, command := range utils.Executables() {
	// 	items = append(items, readline.PcItem(command))
	// }

	return &BellCompleter{
		PrefixCompleter: readline.NewPrefixCompleter(items...),
	}
}
