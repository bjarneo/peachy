package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap defines all key bindings
type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Enter       key.Binding
	Escape      key.Binding
	Open        key.Binding
	Extract     key.Binding
	Save        key.Binding
	Load        key.Binding
	Reset       key.Binding
	Help        key.Binding
	Quit        key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/confirm"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel/back"),
		),
		Open: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open image"),
		),
		Extract: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "extract colors"),
		),
		Save: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "save"),
		),
		Load: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "load"),
		),
		Reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}
