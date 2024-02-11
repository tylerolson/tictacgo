package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerolson/tictacgo/tictacgo"
)

type roomModel struct {
	table    table.Model
	cursor   int
	roomKeys roomKeyMap
}

func newRoomModel() *roomModel {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Players", Width: 10},
	}

	rows := []table.Row{
		{"Test", "1/2"},
		{"Test", "0/2"},
		{"Test", "2/2"},
		{"Test14134", "0/2"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &roomModel{
		table:    t,
		cursor:   0,
		roomKeys: roomKeys,
	}
}

func UpdateTable(t *table.Model, rooms []tictacgo.Room) {
	var rows []table.Row

	for _, v := range rooms {
		rows = append(rows, table.Row{v.Name, strconv.Itoa(v.Size)})
	}

	t.SetRows(rows)
}

func (m roomModel) Init() tea.Cmd {
	UpdateTable(&m.table, GetRooms())
	return nil
}

func (m roomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.roomKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.roomKeys.Refresh):
			UpdateTable(&m.table, GetRooms())
			return m, nil
		case key.Matches(msg, m.roomKeys.Enter):
			gm := newGameModel(m.table.SelectedRow()[0])
			return gm, gm.Init()
		}
	}

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m roomModel) View() string {
	var s strings.Builder

	s.WriteString("Rooms\n")
	s.WriteString(m.table.View())
	s.WriteString("\n" + help.New().View(m.roomKeys))

	return lipgloss.NewStyle().Margin(2, 10).Render(s.String())
}

func GetRooms() []tictacgo.Room {
	res, err := http.Get("http://127.0.0.1:8081/rooms")
	if err != nil {
		log.Fatal(err, "Failed to get")
	}

	var rooms []tictacgo.Room

	err = json.NewDecoder(res.Body).Decode(&rooms)
	if err != nil {
		log.Fatal("Couldn't decode response")
	}

	return rooms
}

//func CreateRoom(room string) {}
