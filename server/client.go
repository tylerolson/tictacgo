package server

import (
	"encoding/json"
	"errors"
	"net"
	"time"

	"github.com/tylerolson/tictacgo"
)

type Client struct {
	Player string
	Game   *tictacgo.Game

	roomName      string
	started       bool
	conn          net.Conn
	updateChannel chan Response
	errorChannel  chan error
}

func (c *Client) GetUpdateChannel() chan Response {
	return c.updateChannel
}

func (c *Client) GetErrorChannel() chan error {
	return c.errorChannel
}

func (c *Client) IsStarted() bool {
	return c.started
}

func NewClient() *Client {
	g := tictacgo.NewGame()
	return &Client{
		Player:        "",
		Game:          g,
		roomName:      "",
		started:       false,
		conn:          nil,
		updateChannel: make(chan Response),
		errorChannel:  make(chan error),
	}
}

func (c *Client) EstablishConnection(ip string) error {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return err
	}

	c.conn = conn

	go c.receiveResponse()

	return nil
}

func (c *Client) receiveResponse() {
	for {
		var rawContent json.RawMessage
		response := Response{
			Content: &rawContent,
		}

		if err := json.NewDecoder(c.conn).Decode(&response); err != nil {
			c.errorChannel <- err
			return
		}

		switch response.Type {
		case AssignMark:
			var content AssignMarkContent

			if err := json.Unmarshal(rawContent, &content); err != nil {
				c.errorChannel <- err
				return
			}

			c.Player = content.Player
		case UpdateGame:
			var content UpdateGameContent

			if err := json.Unmarshal(rawContent, &content); err != nil {
				c.errorChannel <- err
				return
			}

			c.started = content.Started
			c.Game.SetGame(content.Game)
		}

		c.updateChannel <- response

		time.Sleep(50 * time.Millisecond)
	}
}

func (c *Client) MakeMove(move string) error {
	return json.NewEncoder(c.conn).Encode(Request{
		Type: MakeMove,
		Content: MakeMoveContent{
			Room:   c.roomName,
			Player: c.Player,
			Move:   move,
		},
	})
}

func (c *Client) JoinRoom(roomName string) error {
	if c.conn == nil {
		return errors.New("JoinRoom() client connection is nil")
	}

	c.roomName = roomName

	return json.NewEncoder(c.conn).Encode(Request{
		Type: JoinRoom,
		Content: RoomContent{
			Room: roomName,
		},
	})
}
