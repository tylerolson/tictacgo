package client

import (
	"github.com/tylerolson/tictacgo/tictacgo"
	"log"
	"net"
	"strings"
	"time"
)

type client struct {
	room   string
	player string
	conn   net.Conn
	game   tictacgo.Game
}

func (c *client) GetGame() tictacgo.Game {
	return c.game
}

func NewClient(player string) *client {
	g := tictacgo.NewGame()
	return &client{
		room:   "",
		conn:   nil,
		player: player,
		game:   g,
	}
}

func (c *client) EstablishConnection(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(err)
	}

	c.conn = conn

	go c.receiveResponse()
}

func (c *client) sendRequest(message string) {
	_, err := c.conn.Write([]byte(message))
	if err != nil {
		return
	}
}

func (c *client) receiveResponse() {
	for {
		bytes := make([]byte, 1024)
		_, err := c.conn.Read(bytes)
		if err != nil {
			return
		}

		message := strings.Split(string(bytes), ",")

		if message[0] == "BOARD" {
			for i := 0; i < 9; i++ {
				c.game.SetCell(i, message[i+1])
			}
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func (c *client) MakeMove(move string) {
	m := "MAKE_MOVE,"
	m += c.room + ","
	m += c.player + ","
	m += move + ","

	c.sendRequest(m)
}

func (c *client) CreateRoom(room string) {
	c.room = room
	m := "CREATE_ROOM,"
	m += room + ","

	c.sendRequest(m)
}
