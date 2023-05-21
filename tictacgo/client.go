package tictacgo

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type Client struct {
	room   string
	player string
	conn   net.Conn
	game   *Game
}

func (c *Client) GetGame() *Game {
	return c.game
}

func (c *Client) SetPlayer(player string) {
	c.player = player
}

func NewClient() *Client {
	g := NewGame()
	return &Client{
		room:   "",
		conn:   nil,
		player: "",
		game:   g,
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
		var message Message

		decoder := json.NewDecoder(c.conn)
		if err := decoder.Decode(&message); err != nil {
			log.Println(err)
			return
		}

		for i := 0; i < 9; i++ {
			c.game.SetCell(i, message.Board[i])
		}

		c.game.SetTurn(message.Turn)
		c.game.SetWinner(message.Winner)

		time.Sleep(50 * time.Millisecond)
	}
}

func (c *Client) MakeMove(move string) {

	err := json.NewEncoder(c.conn).Encode(Message{
		Request: MAKE_MOVE,
		Room:    c.room,
		Player:  c.player,
		Move:    move,
	})
	if err != nil {
		return
	}
}

func (c *Client) CreateRoom(room string) {
	c.room = room

	err := json.NewEncoder(c.conn).Encode(Message{
		Request: CREATE_ROOM,
		Room:    room,
	})
	if err != nil {
		return
	}
}

func (c *Client) JoinRoom(room string) {
	c.room = room

	err := json.NewEncoder(c.conn).Encode(Message{
		Request: JOIN_ROOM,
		Room:    room,
	})
	if err != nil {
		return
	}
}
