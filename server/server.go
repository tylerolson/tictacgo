package main

import (
	"encoding/json"
	"github.com/tylerolson/tictacgo/tictacgo"
	"log"
	"net"
)

type room struct {
	name string
	game tictacgo.Game
}

var rooms = make(map[string]*room)

func createRoom(name string) {
	g := tictacgo.NewGame()
	r := room{name, g}
	rooms[name] = &r

	log.Println("Created room " + name)
}

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Fatal(err)
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
	for {
		var message map[string]interface{}
		decoder := json.NewDecoder(conn)
		err := decoder.Decode(&message)
		if err != nil {
			return
		}

		if message["Request"] == "CREATE_ROOM" {
			createRoom(message["Room"].(string))
		} else if message["Request"] == "JOIN_ROOM" {
			// something
		} else if message["Request"] == "MAKE_MOVE" { //MAKE_MOVE ROOM PLAYER MOVE
			room, ok := rooms[message["Room"].(string)]
			if !ok {
				log.Println("room '" + message["Room"].(string) + "' does not exist")
				return
			}

			if room.game.GetTurn() == message["Player"].(string) {
				room.game.Move(message["Move"].(string))
			}

			type message struct {
				Request string
				Board   []string
				Turn    string
			}

			mess := message{
				Request: "UPDATE",
				Board:   room.game.GetBoard(),
				Turn:    room.game.GetTurn(),
			}

			encoder := json.NewEncoder(conn)
			err := encoder.Encode(mess)
			if err != nil {
				return
			}
		}

	}
}
