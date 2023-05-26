package main

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tylerolson/tictacgo/tictacgo"
	"net"
	"os"
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

	log.Info().Msg("Created room " + name)
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	s := server{
		rooms: make(map[string]*room),
	}

	var err error
	s.listener, err = net.Listen("tcp", ":8080")
	if err != nil {
		log.Error().Err(err).Msg("listener failed to start")
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Error().Err(err).Msg("deferred listener")
		}
	}(s.listener)

	log.Info().Msg("Running on port 8080")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("failed to accepted connection")
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

		log.Debug().
			Str("Request", message.Request).
			Str("Room", message.Room).
			Str("Player", message.Player).
			Str("Move", message.Move).
			Str("Turn", message.Turn).
			Str("Winner", message.Winner).
			Strs("Board", message.Board).
			Msg("Received message")

		if message.Request == tictacgo.CreateRoom {
			s.createRoom(message.Room)
			s.rooms[message.Room].connections[conn.RemoteAddr().String()] = conn
		} else if message.Request == tictacgo.JoinRoom {
			r, ok := s.rooms[message.Room]
			if !ok {
				log.Warn().Msg("Room does not exist")
				return
			}
			r.connections[conn.RemoteAddr().String()] = conn
		} else if message.Request == tictacgo.MakeMove { //MakeMove ROOM PLAYER MOVE
			room, ok := s.rooms[message.Room]
			if !ok {
				log.Warn().Msg("room '" + message.Room + "' does not exist")
				return
			}

			if message.Player == "" {
				log.Warn().Msg("Message Player does not exist")
				return
			}

			if room.game.GetTurn() == message.Player {
				room.game.Move(message.Move)
				log.Debug().Msg("Made move " + message.Move + "in room " + room.name)
			}

			s.broadcastUpdates(message.Room)
		}

	}
}

func (s *server) broadcastUpdates(roomStr string) {
	mess := tictacgo.Message{
		Request: tictacgo.Update,
		Room:    roomStr,
		Move:    s.rooms[roomStr].game.GetTurn(),
		Board:   s.rooms[roomStr].game.GetBoard(),
		Turn:    s.rooms[roomStr].game.GetTurn(),
		Winner:  s.rooms[roomStr].game.GetWinner(),
	}
	res := tictacgo.Response{
		Code:    tictacgo.Success,
		Message: mess,
	}
	for _, room := range s.rooms {
		if room.name == roomStr {
			for _, conn := range room.connections {
				encoder := json.NewEncoder(conn)
				err := encoder.Encode(res)
				if err != nil {
					return
				}
			}
		}
	}
}
