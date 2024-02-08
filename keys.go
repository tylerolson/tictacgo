package main

import "github.com/charmbracelet/bubbles/key"

type menuKeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Quit  key.Binding
	Enter key.Binding
}

type roomKeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Refresh key.Binding
	Quit    key.Binding
	Enter   key.Binding
}

type gameKeyMap struct {
	Move key.Binding
	Quit key.Binding
}

func (k menuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k roomKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Refresh, k.Enter, k.Quit}
}

func (k gameKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Move, k.Quit}
}

func (k menuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (k roomKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

func (k gameKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var menuKeys = menuKeyMap{
	Up: key.NewBinding(
		key.WithKeys("w", "up"),
		key.WithHelp("↑/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("s", "down"),
		key.WithHelp("↓/s", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "make selection"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var roomKeys = roomKeyMap{
	Up: key.NewBinding(
		key.WithKeys("w", "up"),
		key.WithHelp("↑/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("s", "down"),
		key.WithHelp("↓/s", "move down"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh rooms"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "make selection"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

var gameKeys = gameKeyMap{
	Move: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("1-9", "make move"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
