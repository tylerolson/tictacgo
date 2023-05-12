package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strconv"
	"strings"
)

type Game struct {
	board  []string
	turn   string
	winner string
}

func (g *Game) checkWinner() bool {
	// horizontal

	if g.board[0] == g.board[1] && g.board[1] == g.board[2] {
		g.winner = g.board[0]
		return true
	}

	if g.board[3] == g.board[4] && g.board[4] == g.board[5] {
		g.winner = g.board[3]
		return true
	}

	if g.board[6] == g.board[7] && g.board[7] == g.board[8] {
		g.winner = g.board[6]
		return true
	}

	//vertical

	if g.board[0] == g.board[3] && g.board[3] == g.board[6] {
		g.winner = g.board[0]
		return true
	}

	if g.board[1] == g.board[4] && g.board[4] == g.board[7] {
		g.winner = g.board[1]
		return true
	}

	if g.board[2] == g.board[5] && g.board[5] == g.board[8] {
		g.winner = g.board[2]
		return true
	}

	//diagonal

	if g.board[0] == g.board[4] && g.board[5] == g.board[8] {
		g.winner = g.board[0]
		return true
	}

	if g.board[2] == g.board[4] && g.board[5] == g.board[6] {
		g.winner = g.board[2]
		return true
	}

	return false
}

func (g *Game) move(cell string) bool {
	if g.checkWinner() {
		return false //throw err maybe?
	}

	cellInt, _ := strconv.Atoi(cell)
	cellInt--

	if g.board[cellInt] == cell {
		g.board[cellInt] = g.turn
	} else {
		return false
	}

	return true
}

func newGame() Game {
	game := &Game{
		board:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		turn:   "X",
		winner: "",
	}

	return *game
}

type model struct {
	state   string
	choices []string
	cursor  int

	game  Game
	board table.Model
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
			if !m.game.move(msg.String()) {
				return m, nil
			}

			r := m.board.Rows()
			num, _ := strconv.Atoi(msg.String())
			num--

			if r[num/3][num%3] == msg.String() {
				r[num/3][num%3] = m.game.turn
			}

			m.board.SetRows(r)

			if m.game.turn == "X" {
				m.game.turn = "O"
			} else {
				m.game.turn = "X"
			}

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
		lipgloss.NewStyle().PaddingLeft(7).Render(m.board.View()),
	)

	if m.game.checkWinner() {
		s.WriteString("\n\n")
		s.WriteString(m.game.winner)
		s.WriteString(" wins!")
	} else {
		s.WriteString("\n\n It is ")
		s.WriteString(m.game.turn)
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

	g := newGame()

	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle()
	s.Header = lipgloss.NewStyle()
	s.Cell = lipgloss.NewStyle()
	s.Cell = s.Cell.Border(lipgloss.NormalBorder()).Bold(false).Align(lipgloss.Center, lipgloss.Center)
	t.SetStyles(s)

	initialModel := model{
		state:   "menu",
		choices: []string{"Start", "Exit"},
		cursor:  0,
		game:    g,
		board:   t,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
