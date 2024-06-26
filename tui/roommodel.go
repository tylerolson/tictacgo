package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tylerolson/tictacgo/server"
)

type roomModel struct {
	table     table.Model
	cursor    int
	roomKeys  roomKeyMap
	err       error
	textInput textinput.Model
}

func newRoomModel() roomModel {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Players", Width: 10},
	}

	rows := []table.Row{
		{"Rooms not avaliable", "?/2"},
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

	textInput := textinput.New()
	textInput.Blur()
	textInput.Placeholder = "Enter room name"
	textInput.Width = 50

	return roomModel{
		table:     t,
		cursor:    0,
		roomKeys:  roomKeys,
		textInput: textInput,
	}
}

func (m roomModel) Init() tea.Cmd {
	return updateTable(m.table)
}

func (m roomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case error:
		m.err = msg
	case table.Model:
		m.table = msg
	case tea.KeyMsg:
		switch { // TODO add escape to cancel room create
		case key.Matches(msg, m.roomKeys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.roomKeys.Refresh):
			return m, updateTable(m.table)
		case key.Matches(msg, m.roomKeys.Create):
			m.table.Blur()
			return m, m.textInput.Focus()
		case key.Matches(msg, m.roomKeys.Enter):
			if m.textInput.Focused() {
				m.err = createRoom(m.textInput.Value())
				m.textInput.Blur()
				m.textInput.Reset()
				m.table.Focus()
				return m, updateTable(m.table)
			}
			if !strings.Contains(m.table.SelectedRow()[1], "?") {
				gm := newGameModel(m.table.SelectedRow()[0])
				return gm, gm.Init()
			}

		}
	}

	var textCmd, tableCmd tea.Cmd
	m.textInput, textCmd = m.textInput.Update(msg)
	m.table, tableCmd = m.table.Update(msg)
	return m, tea.Batch(cmd, textCmd, tableCmd)
}

func (m roomModel) View() string {
	var s strings.Builder

	s.WriteString(" Rooms:\n\n")
	s.WriteString(m.table.View() + "\n\n")

	if m.textInput.Focused() {
		s.WriteString(m.textInput.View() + "\n")
	} else {
		s.WriteString("\n")
	}

	s.WriteString("\n" + help.New().View(m.roomKeys) + "\n\n")

	errorMsg := ""
	if m.err != nil {
		errorMsg = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("%+v", m.err))
	}

	return lipgloss.NewStyle().Margin(2, 10).Render(s.String() + errorMsg)
}

func updateTable(t table.Model) tea.Cmd { // tea.Cmd
	return func() tea.Msg {
		rooms, err := getRooms()
		if err != nil {
			return err
		}

		var rows []table.Row

		for _, v := range rooms {
			rows = append(rows, table.Row{v.Name, strconv.Itoa(v.Size)})
		}

		t.SetRows(rows)

		return t
	}
}

func getRooms() ([]server.RoomResponse, error) {
	res, err := http.Get("http://127.0.0.1:8081/rooms")
	if err != nil {
		return nil, fmt.Errorf("REST server is not running\n%w", err)
	}
	defer res.Body.Close()

	var rawContent json.RawMessage
	response := server.Response{
		Content: &rawContent,
	}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode REST response\n%w", err)
	}

	if response.Type == server.GetRoom {
		var content []server.RoomResponse

		if err = json.Unmarshal(rawContent, &content); err != nil {
			return nil, fmt.Errorf("failed to unmarshall REST RoomResponse\n%w", err)
		}

		return content, err
	}

	return nil, err
}

func createRoom(roomName string) error {
	request := server.Request{
		Type: server.MakeRoom,
		Content: server.RoomContent{
			Room: roomName,
		},
	}

	buff := bytes.NewBuffer(make([]byte, 0))

	if err := json.NewEncoder(buff).Encode(request); err != nil {
		return fmt.Errorf("failed to encode REST request\n%w", err)
	}

	if _, err := http.Post("http://127.0.0.1:8081/rooms", "application/json", buff); err != nil {
		return fmt.Errorf("failed to POST CreateRoom\n%w", err)
	}

	return nil
}
