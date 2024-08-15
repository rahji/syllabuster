package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Tab      key.Binding
	Generate key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
// It's part of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Generate, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
// It's part of the key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return nil
}

var keys = keyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next field"),
	),
	Generate: key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("ctr+g", "generate files"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}
