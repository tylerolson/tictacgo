package main

import (
	"github.com/tylerolson/tictacgo/tictacgo"
	"io"
	"log"
	"net"
	"strings"
)

var rooms = make(map[string]room)

func createRoom(name string) {
	g := tictacgo.NewGame()
	r := room{name, g}
	rooms[name] = r

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
		bytes := make([]byte, 1024)
		_, err := conn.Read(bytes)
		if err != nil && err == io.EOF {
			conn.Close()
			return
		} else if err != nil {
			log.Fatal(err)
		}

		log.Println(string(bytes))

		message := strings.Split(string(bytes), ",")

		if message[0] == "CREATE_ROOM" {
			createRoom(message[1])
		} else if message[0] == "JOIN_ROOM" {
			// something
		} else if message[0] == "MAKE_MOVE" { //MAKE_MOVE ROOM PLAYER MOVE
			g, ok := rooms[message[1]]
			if !ok {
				log.Println("room '" + message[1] + "' does not exist")
				return
			}

			if g.game.GetTurn() == message[2] {
				g.game.Move(message[3])
			}

			board := "BOARD,"

			for i := 0; i < 9; i++ {
				board += g.game.GetBoard()[i] + ","
			}

			_, err := conn.Write([]byte(board))
			if err != nil {
				return
			}

		}
	}
}
