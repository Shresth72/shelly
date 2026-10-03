package shelly

import "github.com/chzyer/readline"

func NewAutoCompleter() *readline.PrefixCompleter {
	items := make([]*readline.PrefixCompleter, 0, len(allCommands))

	for _, command := range allCommands {
		items = append(items, readline.PcItem(command))
	}

	return readline.NewPrefixCompleter(
		convertToInterfaces(items)...,
	)
}

func convertToInterfaces(
	items []*readline.PrefixCompleter,
) []readline.PrefixCompleterInterface {
	result := make([]readline.PrefixCompleterInterface, len(items))

	for i, item := range items {
		result[i] = item
	}

	return result
}
