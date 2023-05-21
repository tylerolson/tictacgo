package client

import (
	"encoding/json"
	"github.com/tylerolson/tictacgo/tictacgo"
	"log"
	"net"
	"time"
)

type Client struct {
	room   string
	player string
	conn   net.Conn
	game   tictacgo.Game
}

func (c *Client) GetGame() *tictacgo.Game {
	return &c.game
}

func NewClient(player string) Client {
	g := tictacgo.NewGame()
	return Client{
		room:   "",
		conn:   nil,
		player: player,
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

func (c *Client) sendRequest(message interface{}) {
	encoder := json.NewEncoder(c.conn)
	err := encoder.Encode(message)
	if err != nil {
		return
	}
}

func (c *Client) receiveResponse() {
	for {
		var message map[string]interface{}

		decoder := json.NewDecoder(c.conn)
		if err := decoder.Decode(&message); err != nil {
			log.Println(err)
			return
		}

		for i := 0; i < 9; i++ {
			c.game.SetCell(i, message["Board"].([]interface{})[i].(string))
		}

		c.game.SetTurn(message["Turn"].(string))

		time.Sleep(50 * time.Millisecond)
	}
}

func (c *Client) MakeMove(move string) {
	type message struct {
		Request string
		Room    string
		Player  string
		Move    string
	}

	c.sendRequest(message{
		Request: "MAKE_MOVE",
		Room:    c.room,
		Player:  c.player,
		Move:    move,
	})
}

func (c *Client) CreateRoom(room string) {
	c.room = room
	type message struct {
		Request string
		Room    string
	}

	c.sendRequest(message{
		Request: "CREATE_ROOM",
		Room:    room,
	})
}
