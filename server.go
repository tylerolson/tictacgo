package main

import (
	"encoding/json"
	"github.com/tylerolson/tictacgo/tictacgo"
	"log"
	"net"
)

type room struct {
	name        string
	game        *tictacgo.Game
	connections map[string]net.Conn
}

type server struct {
	listener net.Listener
	rooms    map[string]*room
}

func (s *server) createRoom(name string) {
	g := tictacgo.NewGame()
	r := room{
		name:        name,
		game:        g,
		connections: make(map[string]net.Conn),
	}
	s.rooms[name] = &r

	log.Println("Created room " + name)
}

func main() {
	s := server{
		rooms: make(map[string]*room),
	}

	var err error
	s.listener, err = net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(s.listener)

	log.Println("Running on port 8080")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
	for {
		var message tictacgo.Message
		decoder := json.NewDecoder(conn)
		err := decoder.Decode(&message)
		if err != nil {
			return
		}

		if message.Request == tictacgo.CREATE_ROOM {
			s.createRoom(message.Room)
			s.rooms[message.Room].connections[conn.RemoteAddr().String()] = conn
		} else if message.Request == tictacgo.JOIN_ROOM {
			s.rooms[message.Room].connections[conn.RemoteAddr().String()] = conn
		} else if message.Request == tictacgo.MAKE_MOVE { //MAKE_MOVE ROOM PLAYER MOVE
			room, ok := s.rooms[message.Room]
			if !ok {
				log.Println("room '" + message.Room + "' does not exist")
				return
			}

			if room.game.GetTurn() == message.Player {
				room.game.Move(message.Move)
			}

			s.broadcastUpdates(message.Room)
		}

	}
}

func (s *server) broadcastUpdates(roomStr string) {
	mess := tictacgo.Message{
		Request: tictacgo.UPDATE,
		Board:   s.rooms[roomStr].game.GetBoard(),
		Turn:    s.rooms[roomStr].game.GetTurn(),
		Winner:  s.rooms[roomStr].game.GetWinner(),
	}
	for _, room := range s.rooms {
		if room.name == roomStr {
			for _, conn := range room.connections {
				encoder := json.NewEncoder(conn)
				err := encoder.Encode(mess)
				if err != nil {
					return
				}
			}
		}
	}
}
