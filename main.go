package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerolson/tictacgo/tictacgo"
	"os"
	"strings"
)

type menuModel struct {
	choices  []string
	cursor   int
	menuKeys menuKeyMap
}

func newMenuModel() *menuModel {
	return &menuModel{
		choices:  []string{"Start Solo", "Create Room", "Join Room", "Exit"},
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
				return newGameModel(true), nil
			} else if m.cursor == 1 { //create room
				return newGameModel(false), nil
			} else if m.cursor == 2 { //join room
				return newGameModel(false), nil
			} else if m.cursor == 3 { //exit
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m menuModel) View() string {
	var s strings.Builder
	s.WriteString("\n")

	for i, choice := range m.choices {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = "> " // cursor!
		}
		s.WriteString(cursor)
		s.WriteString(choice)
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(help.New().View(m.menuKeys))
	s.WriteString("\n")

	var marginStyle = lipgloss.NewStyle().MarginLeft(10)

	return marginStyle.Render(s.String())
}

// game model

type gameModel struct {
	game       *tictacgo.Game
	boardTable table.Model
	gameKeys   gameKeyMap
	isLocal    bool
	client     *tictacgo.Client
}

func newGameModel(local bool) *gameModel {
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

	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle()
	s.Header = lipgloss.NewStyle()
	s.Cell = lipgloss.NewStyle()
	s.Cell = s.Cell.Border(lipgloss.NormalBorder()).Bold(false).Align(lipgloss.Center, lipgloss.Center)
	t.SetStyles(s)

	if local {
		return &gameModel{
			game:       tictacgo.NewGame(),
			boardTable: t,
			gameKeys:   gameKeys,
			isLocal:    local,
		}
	} else {
		c := tictacgo.NewClient()
		c.EstablishConnection("localhost:8080")
		c.CreateRoom("yo")
		c.SetPlayer("X")
		return &gameModel{
			game:       tictacgo.NewGame(),
			boardTable: t,
			gameKeys:   gameKeys,
			isLocal:    local,
			client:     c,
		}
	}
}

func (gm gameModel) Init() tea.Cmd {
	return nil
}

func (gm gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, gm.gameKeys.Quit):
			return newMenuModel(), nil
		case key.Matches(msg, gm.gameKeys.Move):
			r := gm.boardTable.Rows()
			if gm.isLocal {
				if !gm.game.Move(msg.String()) {
					return gm, nil
				}
				for i := 0; i < 9; i++ {
					r[i/3][i%3] = gm.game.GetBoard()[i]
				}
			} else {
				gm.client.MakeMove(msg.String())
				for i := 0; i < 9; i++ {
					r[i/3][i%3] = gm.client.GetGame().GetBoard()[i]
				}
			}
			gm.boardTable.SetRows(r)
			return gm, nil
		}
	}

	return gm, nil
}

func (gm gameModel) View() string {
	s := strings.Builder{}

	s.WriteString(gm.boardTable.View())
	s.WriteString("\n\n")

	game := gm.game
	if !gm.isLocal {
		game = gm.client.GetGame()
	}

	if game.CheckWinner() {
		if game.GetWinner() == "tie" {
			s.WriteString("It is a tie!")
		} else {
			s.WriteString(game.GetWinner())
			s.WriteString(" wins!")
		}
	} else {
		s.WriteString("It is ")
		s.WriteString(game.GetTurn())
		s.WriteString("'s turn")
	}

	s.WriteString("\n")
	s.WriteString("\n")
	s.WriteString(help.New().View(gm.gameKeys))
	s.WriteString("\n")

	return lipgloss.NewStyle().MarginLeft(10).Render(s.String())
}

func main() {
	p := tea.NewProgram(newMenuModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
