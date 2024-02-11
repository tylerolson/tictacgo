package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type menuModel struct {
	choices  []string
	cursor   int
	menuKeys menuKeyMap
}

func newMenuModel() *menuModel {
	return &menuModel{
		choices:  []string{"Start Solo", "Multiplayer", "Exit"},
		cursor:   0,
		menuKeys: menuKeys,
	}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.menuKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.menuKeys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.menuKeys.Down):
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.menuKeys.Enter):
			if m.cursor == 0 { //local
				gm := newGameModel("")
				return gm, gm.Init()
			} else if m.cursor == 1 { //create room
				rm := newRoomModel()
				return rm, rm.Init()
			} else if m.cursor == 2 { //exit
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m menuModel) View() string {
	var s strings.Builder

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = "> "
		}
		s.WriteString(cursor + choice + "\n")
	}
	s.WriteString("\n\n" + help.New().View(m.menuKeys))

	return lipgloss.NewStyle().Margin(2, 10).Render(s.String())
}
