/*
Copyright © 2024 Shoaibashk Shoaibashk.2000@gmail.com
*/
package keys

import "github.com/charmbracelet/bubbles/key"

// Map defines the key bindings
type Map struct {
	Quit       key.Binding
	Menu       key.Binding
	Enter      key.Binding
	Up         key.Binding
	Down       key.Binding
	Escape     key.Binding
	CtrlA      key.Binding
	Help       key.Binding
	Connect    key.Binding
	Disconnect key.Binding
}

// DefaultKeys returns the default key bindings
func DefaultKeys() Map {
	return Map{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Menu: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("^A", "menu"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("↵", "select"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓", "down"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		CtrlA: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("^A", "command"),
		),
		Help: key.NewBinding(
			key.WithKeys("z"),
			key.WithHelp("Z", "help"),
		),
		Connect: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("C", "connect"),
		),
		Disconnect: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("D", "disconnect"),
		),
	}
}
