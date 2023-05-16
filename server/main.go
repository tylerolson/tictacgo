package main

import (
	"encoding/gob"
	"github.com/tylerolson/tictacgo/game"
	_ "github.com/tylerolson/tictacgo/game"
	"log"
	"net"
)

var rooms = make(map[string]room)

func printMessage(message game.Message) {
	log.Println("REQUEST:", message.Request)
	log.Println("ROOM NAME:", message.RoomName)
}

func createRoom(name string) {
	g := game.NewGame()
	r := room{name, g}
	rooms[name] = r
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
		}
	}(listen)

	log.Println("Running on port 8080")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	m := new(game.Message)
	gobDecoder := gob.NewDecoder(conn)
	err := gobDecoder.Decode(m)
	if err != nil {
		return
	}

	if m.Request == game.CREATE_ROOM {
		createRoom(m.RoomName)
	}

	printMessage(*m)

	err = conn.Close()
	if err != nil {
		return
	}
}
