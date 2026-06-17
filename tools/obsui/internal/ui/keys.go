package ui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	NextTab  key.Binding
	PrevTab  key.Binding
	Restart  key.Binding
	Quit     key.Binding
	Up       key.Binding
	Down     key.Binding
	PageUp   key.Binding
	PageDown key.Binding
}

var keys = keyMap{
	NextTab:  key.NewBinding(key.WithKeys("tab", "right"), key.WithHelp("tab/→", "next tab")),
	PrevTab:  key.NewBinding(key.WithKeys("shift+tab", "left"), key.WithHelp("shift+tab/←", "prev tab")),
	Restart:  key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "restart process")),
	Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
	Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
	Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
	PageUp:   key.NewBinding(key.WithKeys("pgup")),
	PageDown: key.NewBinding(key.WithKeys("pgdown")),
}
