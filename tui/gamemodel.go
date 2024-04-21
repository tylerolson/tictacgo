package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerolson/tictacgo"
	"github.com/tylerolson/tictacgo/server"
)

type gameModel struct {
	game       *tictacgo.Game
	boardTable table.Model
	gameKeys   gameKeyMap
	room       string
	client     *server.Client
	err        error
}

func newGameModel(room string) gameModel {
	columns := []table.Column{{Title: "", Width: 1}, {Title: "", Width: 1}, {Title: "", Width: 1}}
	rows := []table.Row{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}}
	styles := table.Styles{
		Header:   lipgloss.NewStyle(),
		Cell:     lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Align(lipgloss.Center, lipgloss.Center),
		Selected: lipgloss.NewStyle(),
	}

	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithHeight(10), table.WithStyles(styles))

	gm := gameModel{
		game:       tictacgo.NewGame(),
		boardTable: t,
		gameKeys:   gameKeys,
		room:       room,
	}

	if room != "" {
		c := server.NewClient()
		if gm.err = c.EstablishConnection("localhost:8080"); gm.err == nil {
			c.JoinRoom(room)
			c.Player = "X"

			gm.client = c
		}

	}

	return gm
}

func receiveUpdate(channel chan server.Response) tea.Cmd {
	return func() tea.Msg {
		return <-channel
	}
}

func receiveError(channel chan error) tea.Cmd {
	return func() tea.Msg {
		return <-channel
	}
}

func (gm gameModel) getUpdatedTable(board []string) gameModel {
	r := gm.boardTable.Rows()
	for i := 0; i < 9; i++ {
		r[i/3][i%3] = board[i]
	}
	gm.boardTable.SetRows(r)
	return gm
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
	case server.Response:
		switch msg.Type {
		case server.UpdateGame:
			gm = gm.getUpdatedTable(gm.client.Game.Board)
		}
		return gm, receiveUpdate(gm.client.GetUpdateChannel())
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, gm.gameKeys.Quit):
			if gm.room == "" {
				return newMenuModel(), nil
			} else {
				gm.client.CloseConnection()
				rm := newRoomModel()
				return rm, rm.Init()
			}
		case key.Matches(msg, gm.gameKeys.Move):
			if gm.room != "" {
				gm.err = gm.client.MakeMove(msg.String())
				return gm, nil
			}

			if gm.game.Move(msg.String()) {
				gm = gm.getUpdatedTable(gm.game.Board)
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
		game = gm.client.Game

		if !gm.client.IsStarted() {
			s.WriteString("Waiting for other player...\n")
		}
	}

	if !game.HasWinner() {
		s.WriteString("It is " + game.Turn + "'s turn")
	} else if game.Winner == "tie" {
		s.WriteString("It is a tie!")
	} else {
		s.WriteString(game.Winner + " wins!")
	}

	if gm.room != "" {
		s.WriteString("\nYou are " + gm.client.Player)
	}

	s.WriteString("\n\n\n" + help.New().View(gm.gameKeys) + "\n\n")

	errorMsg := ""
	if gm.err != nil {
		errorMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("%+v", gm.err))
	}

	return lipgloss.NewStyle().Margin(2, 10).Render(s.String() + errorMsg)
}
