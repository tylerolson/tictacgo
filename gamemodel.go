package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerolson/tictacgo/tictacgo"
)

type gameModel struct {
	game       *tictacgo.Game
	boardTable table.Model
	gameKeys   gameKeyMap
	room       string
	client     *tictacgo.Client
	err        error
}

func newGameModel(room string) *gameModel {
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
		room:       room,
	}

	if room != "" {
		c := tictacgo.NewClient()
		if gm.err = c.EstablishConnection("localhost:8080"); gm.err == nil {
			c.JoinRoom(room)
			c.SetPlayer("X")

			gm.client = c
		}

	}

	return &gm
}

func receiveUpdate(channel chan tictacgo.Response) tea.Cmd {
	return func() tea.Msg {
		return <-channel
	}
}

func receiveError(channel chan error) tea.Cmd {
	return func() tea.Msg {
		return <-channel
	}
}

func (gm gameModel) Init() tea.Cmd {
	upCmd := receiveUpdate(gm.client.GetUpdateChannel())
	errCmd := receiveError(gm.client.GetErrorChannel())
	return tea.Batch(upCmd, errCmd)
}

func (gm gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		gm.err = msg
	case tictacgo.Response:
		g := gm.client.GetGame()
		g.SetTurn(msg.Turn)
		g.SetWinner(msg.Winner)
		if msg.Board != nil {
			g.SetBoard(msg.Board) //possibly can be null
			r := gm.boardTable.Rows()
			for i := 0; i < 9; i++ {
				r[i/3][i%3] = gm.client.GetGame().GetBoard()[i]
			}
			gm.boardTable.SetRows(r)
		}

		return gm, receiveUpdate(gm.client.GetUpdateChannel())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, gm.gameKeys.Quit):
			return newMenuModel(), nil
		case key.Matches(msg, gm.gameKeys.Move):
			if gm.room != "" {
				gm.err = gm.client.MakeMove(msg.String())
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
	if gm.room != "" {
		game = gm.client.GetGame()
	}

	if game.GetWinner() == "" {
		if game.GetTurn() == "" {
			s.WriteString("Waiting for other player...")
		} else {
			s.WriteString("It is " + game.GetTurn() + "'s turn")
		}
	} else if game.GetWinner() == "tie" {
		s.WriteString("It is a tie!")
	} else {
		s.WriteString(game.GetWinner() + " wins!")
	}

	if gm.room != "" {
		s.WriteString("\nYou are " + gm.client.GetPlayer())
	}

	s.WriteString("\n\n\n" + help.New().View(gm.gameKeys) + "\n\n")

	errorMsg := ""
	if gm.err != nil {
		errorMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(gm.err.Error())
	}

	return lipgloss.NewStyle().Margin(2, 10).Render(s.String() + errorMsg)
}
