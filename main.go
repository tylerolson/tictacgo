package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerolson/tictacgo/game"
	"os"
	"strings"
)

type model struct {
	state   string
	choices []string
	cursor  int

	game       game.TicTacGo
	boardTable table.Model
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

// sub updates

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
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if !m.game.Move(msg.String()) {
				return m, nil
			}

			r := m.boardTable.Rows()

			for i := 0; i < 9; i++ {
				r[i/3][i%3] = m.game.GetBoard()[i]
			}

			m.boardTable.SetRows(r)

			return m, nil
		}
	}

	return m, nil
}

//sub views

func (m model) ViewMenu() string {
	s := strings.Builder{}

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		s.WriteString("      ")
		s.WriteString(cursor)
		s.WriteString(choice)
		s.WriteString("\n")
	}

	return s.String()

}

func (m model) ViewGame() string {
	s := strings.Builder{}

	s.WriteString(
		lipgloss.NewStyle().PaddingLeft(7).Render(m.boardTable.View()),
	)

	if m.game.CheckWinner() {
		if m.game.GetWinner() == "tie" {
			s.WriteString("\n\nIt is a tie!")
		} else {
			s.WriteString("\n\n")
			s.WriteString(m.game.GetWinner())
			s.WriteString(" wins!")
		}
	} else {
		s.WriteString("\n\n It is ")
		s.WriteString(m.game.GetWinner())
		s.WriteString("'s turn")
	}

	return s.String()

}

func main() {
	columns := []table.Column{
		{Title: "", Width: 3},
		{Title: "", Width: 3},
		{Title: "", Width: 3},
	}
	rows := []table.Row{
		{"1", "2", "3"},
		{"4", "5", "6"},
		{"7", "8", "9"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(10),
	)

	g := game.NewGame()

	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle()
	s.Header = lipgloss.NewStyle()
	s.Cell = lipgloss.NewStyle()
	s.Cell = s.Cell.Border(lipgloss.NormalBorder()).Bold(false).Align(lipgloss.Center, lipgloss.Center)
	t.SetStyles(s)

	initialModel := model{
		state:      "menu",
		choices:    []string{"Start", "Exit"},
		cursor:     0,
		game:       g,
		boardTable: t,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
