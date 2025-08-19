package ui

import (
	key "github.com/charmbracelet/bubbles/v2/key"
)

// keymap defines the app-level key bindings.
type keymap struct {
	MoveUp     key.Binding
	MoveDown   key.Binding
	OpenAction key.Binding
	Pick       key.Binding
	Squash     key.Binding
	Fixup      key.Binding
	Edit       key.Binding
	Drop       key.Binding
	Rebase     key.Binding
	Quit       key.Binding
}

var keys = keymap{
	MoveUp:     key.NewBinding(key.WithKeys("ctrl+up"), key.WithHelp("ctrl+↑", "move up")),
	MoveDown:   key.NewBinding(key.WithKeys("ctrl+down"), key.WithHelp("ctrl+↓", "move down")),
	OpenAction: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "set action")),
	Pick:       key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "pick")),
	Squash:     key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "squash")),
	Fixup:      key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "fixup")),
	Edit:       key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	Drop:       key.NewBinding(key.WithKeys("x", "d"), key.WithHelp("x/d", "drop")),
	Rebase:     key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "start rebase")),
	Quit:       key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("ctrl+c/q", "quit")),
}
