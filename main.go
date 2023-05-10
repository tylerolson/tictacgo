package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
)

type model struct {
	state   string
	choices []string
	cursor  int

	board [][]string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	if m.state == "menu" {
		return m.UpdateMenu(msg)
	} else if m.state == "game" {
		return m.UpdateGame(msg)
	}

	return m, nil
}

func (m model) View() string {
	if m.state == "menu" {
		return m.ViewMenu()
	} else if m.state == "game" {
		return m.ViewGame()
	}

	return "oops"
}

// sub update

func (m model) UpdateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter", " ":
			if m.cursor == 0 {
				m.state = "game"
				return m, nil
			} else if m.cursor == 1 {
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m model) UpdateGame(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		}
	}

	return m, nil
}

//sub view

func (m model) ViewMenu() string {
	s := strings.Builder{}

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		s.WriteString(cursor)
		s.WriteString(choice)
		s.WriteString("\n")
	}

	return s.String()

}

func (m model) ViewGame() string {
	s := strings.Builder{}

	for i := 0; i < 3; i++ {
		s.WriteString(m.board[i][0] + " | " + m.board[i][1] + " | " + m.board[i][2])
		if i != 2 {
			s.WriteString("\n----------")
		}
		s.WriteString("\n")
	}

	return s.String()

}

func main() {
	initialModel := model{
		state:   "menu",
		choices: []string{"Start", "Exit"},
		cursor:  0,
		board:   [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}},
	}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
