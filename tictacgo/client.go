package tictacgo

import (
	"encoding/json"
	"errors"
	"net"
	"time"
)

type Client struct {
	room          string
	player        string
	conn          net.Conn
	game          *Game
	updateChannel chan Response
	errorChannel  chan error
}

func (c *Client) GetPlayer() string {
	return c.player
}

func (c *Client) GetGame() *Game {
	return c.game
}

func (c *Client) GetUpdateChannel() chan Response {
	return c.updateChannel
}

func (c *Client) GetErrorChannel() chan error {
	return c.errorChannel
}

func (c *Client) SetPlayer(player string) {
	c.player = player
}

func NewClient() *Client {
	g := NewGame()
	return &Client{
		room:          "",
		conn:          nil,
		player:        "",
		game:          g,
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

func (c *Client) receiveResponse() error {
	for {
		var response Response

		decoder := json.NewDecoder(c.conn)
		if err := decoder.Decode(&response); err != nil {
			c.errorChannel <- err

		}

		if response.Board != nil {
			for i := 0; i < 9; i++ {
				c.game.SetCell(i, response.Board[i])
			}
		}

		if response.Player != "" {
			c.SetPlayer(response.Player)
		}

		if response.Turn != "" {
			c.game.SetTurn(response.Turn)
		}

		if response.Winner != "" {
			c.game.SetWinner(response.Winner)
		}

		c.updateChannel <- response

		time.Sleep(50 * time.Millisecond)
	}
}

func (c *Client) MakeMove(move string) error {
	return json.NewEncoder(c.conn).Encode(Message{
		Request: MakeMove,
		Room:    c.room,
		Player:  c.player,
		Move:    move,
	})
}

func (c *Client) JoinRoom(room string) error {
	if c.conn == nil {
		return errors.New("JoinRoom() client connection is nil")
	}

	c.room = room

	return json.NewEncoder(c.conn).Encode(Message{
		Request: JoinRoom,
		Room:    room,
	})
}
