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
				return newGameModel(true, "X"), nil
			} else if m.cursor == 1 { //create room
				gm := newGameModel(false, "X")
				return gm, gm.Init()
			} else if m.cursor == 2 { //join room
				gm := newGameModel(false, "O")
				return gm, gm.Init()
			} else if m.cursor == 3 { //exit
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

// room model

type roomModel struct {
	choices  []string
	cursor   int
	menuKeys menuKeyMap
}

func newRoomModel() *roomModel {
	return &roomModel{
		choices:  []string{"First", "Second"},
		cursor:   0,
		menuKeys: menuKeys,
	}
}

func (m roomModel) Init() tea.Cmd {
	return nil
}

func (m roomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				return newGameModel(true, "X"), nil
			} else if m.cursor == 1 { //create room
				gm := newGameModel(false, "X")
				return gm, gm.Init()
			} else if m.cursor == 2 { //join room
				gm := newGameModel(false, "O")
				return gm, gm.Init()
			} else if m.cursor == 3 { //exit
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m roomModel) View() string {
	var s strings.Builder

	s.WriteString("Rooms\n")

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

// game model

type gameModel struct {
	game       *tictacgo.Game
	boardTable table.Model
	gameKeys   gameKeyMap
	isLocal    bool
	client     *tictacgo.Client
}

func newGameModel(local bool, player string) *gameModel {
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

	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(10))

	s := table.DefaultStyles()
	s.Selected = lipgloss.NewStyle()
	s.Header = lipgloss.NewStyle()
	s.Cell = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Bold(false).Align(lipgloss.Center, lipgloss.Center)
	t.SetStyles(s)

	gm := gameModel{
		game:       tictacgo.NewGame(),
		boardTable: t,
		gameKeys:   gameKeys,
		isLocal:    local,
	}

	if !local {
		c := tictacgo.NewClient()
		c.EstablishConnection("localhost:8080")
		c.SetPlayer(player)
		if player == "X" {
			c.CreateRoom("yo")
		} else if player == "O" {
			c.JoinRoom("yo")
		}
		gm.client = c

		receiveUpdate(c.GetUpdateChannel())
	}

	return &gm
}

func receiveUpdate(channel chan tictacgo.Response) tea.Cmd {
	return func() tea.Msg {
		return <-channel
	}
}

func (gm gameModel) Init() tea.Cmd {
	return receiveUpdate(gm.client.GetUpdateChannel())
}

func (gm gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tictacgo.Response:
		g := gm.client.GetGame()
		g.SetTurn(msg.Turn)
		g.SetWinner(msg.Winner)
		g.SetBoard(msg.Board)
		r := gm.boardTable.Rows()
		for i := 0; i < 9; i++ {
			r[i/3][i%3] = gm.client.GetGame().GetBoard()[i]
		}
		gm.boardTable.SetRows(r)
		return gm, receiveUpdate(gm.client.GetUpdateChannel())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, gm.gameKeys.Quit):
			return newMenuModel(), nil
		case key.Matches(msg, gm.gameKeys.Move):
			if !gm.isLocal {
				gm.client.MakeMove(msg.String())
				return gm, nil
			}

			if gm.game.Move(msg.String()) {
				r := gm.boardTable.Rows()
				for i := 0; i < 9; i++ {
					r[i/3][i%3] = gm.game.GetBoard()[i]
				}
				gm.boardTable.SetRows(r)
				return gm, nil
			}
		}
	}

	return gm, nil
}

func (gm gameModel) View() string {
	s := strings.Builder{}

	s.WriteString(gm.boardTable.View() + "\n")

	game := gm.game
	if !gm.isLocal {
		game = gm.client.GetGame()
	}

	if game.GetWinner() == "" {
		s.WriteString("It is " + game.GetTurn() + "'s turn")
	} else if game.GetWinner() == "tie" {
		s.WriteString("It is a tie!")
	} else {
		s.WriteString(game.GetWinner() + " wins!")
	}

	if !gm.isLocal {
		s.WriteString("\nYou are " + gm.client.GetPlayer())
	}

	s.WriteString("\n\n\n" + help.New().View(gm.gameKeys))

	return lipgloss.NewStyle().Margin(2, 10).Render(s.String())
}

func main() {
	p := tea.NewProgram(newMenuModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
