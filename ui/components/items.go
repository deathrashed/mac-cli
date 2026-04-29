package components

import "github.com/deathrashed/mac-cli/internal/commands"

type CategoryItem string

func (i CategoryItem) FilterValue() string { return string(i) }
func (i CategoryItem) Title() string       { return string(i) }
func (i CategoryItem) Description() string { return "" }

type CommandItem struct {
	Command commands.Command
}

func (i CommandItem) FilterValue() string {
	return i.Command.Name + " " + i.Command.Description + " " + i.Command.Category
}
func (i CommandItem) Title() string { return i.Command.Name }
func (i CommandItem) Description() string {
	if i.Command.Native {
		return i.Command.Description + " · native"
	}
	if i.Command.Dangerous {
		return i.Command.Description + " · confirm"
	}
	return i.Command.Description
}
