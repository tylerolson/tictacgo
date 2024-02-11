package tictacgo

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type Client struct {
	room          string
	player        string
	conn          net.Conn
	game          *Game
	updateChannel chan Response
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
	}
}

func (c *Client) EstablishConnection(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(err)
	}

	c.conn = conn

	go c.receiveResponse()
}

func (c *Client) receiveResponse() {
	for {
		var response Response

		decoder := json.NewDecoder(c.conn)
		if err := decoder.Decode(&response); err != nil {
			log.Println("Couldn't decode response", err)
			return
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

func (c *Client) MakeMove(move string) {
	err := json.NewEncoder(c.conn).Encode(Message{
		Request: MakeMove,
		Room:    c.room,
		Player:  c.player,
		Move:    move,
	})
	if err != nil {
		log.Fatal("MakeMove failed to send", err)
	}
}

func (c *Client) JoinRoom(room string) {
	c.room = room

	err := json.NewEncoder(c.conn).Encode(Message{
		Request: JoinRoom,
		Room:    room,
	})
	if err != nil {
		log.Fatal("JoinRoom failed to send", err)
	}
}
