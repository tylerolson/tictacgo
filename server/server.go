package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tylerolson/tictacgo/tictacgo"
	"io"
	"net"
	"net/http"
	"os"
)

type player struct {
	mark       string
	connection net.Conn
}

type room struct {
	name    string
	game    *tictacgo.Game
	players map[string]player
}

type server struct {
	rooms       map[string]*room
	restRouter  *gin.Engine
	tcpListener net.Listener
}

func newPlayer(mark string, connection net.Conn) player {
	return player{
		mark:       mark,
		connection: connection,
	}
}

//rooms

func (s *server) makeRoom(name string) {
	g := tictacgo.NewGame()
	r := room{
		name:    name,
		game:    g,
		players: make(map[string]player),
	}
	s.rooms[name] = &r

	log.Info().Msg("Created room " + name)
}

//rest

func (s *server) getRoomsRoute(c *gin.Context) {
	var rooms []tictacgo.Room

	for _, v := range s.rooms {
		room := tictacgo.Room{
			Name: v.name,
			Size: len(v.players),
		}
		rooms = append(rooms, room)
	}

	c.IndentedJSON(http.StatusOK, rooms)
}

func (s *server) makeRoomRoute(c *gin.Context) {
	newBody := struct {
		Name string `json:"name"`
	}{}

	if err := c.BindJSON(&newBody); err != nil {
		log.Fatal().Err(err).Msg("Failed to read POST")
	}

	log.Info().Str("name", newBody.Name)

	s.makeRoom(newBody.Name)

	c.IndentedJSON(http.StatusCreated, newBody)
}

func (s *server) startRESTServer() {
	s.restRouter = gin.New()
	s.restRouter.Use(gin.Logger())
	s.restRouter.GET("/rooms", s.getRoomsRoute)
	s.restRouter.POST("/rooms", s.makeRoomRoute)
	err := s.restRouter.Run("localhost:8081")
	if err != nil {
		return
	}
}

// tcp

func (s *server) handleConnection(conn net.Conn) {
	for {
		var message tictacgo.Message

		if err := json.NewDecoder(conn).Decode(&message); err != nil {
			if err == io.EOF {
				log.Info().Msg("Client disconnected")
				break
				//todo remove client
			} else {
				log.Fatal().Err(err).Msg("Couldn't decode message")
			}
		}

		address := conn.RemoteAddr().String()

		log.Debug().
			Str("Request", message.Request).
			Str("Room", message.Room).
			Str("Player", message.Player).
			Str("Move", message.Move).
			Str("Turn", message.Turn).
			Str("Winner", message.Winner).
			Strs("Board", message.Board).
			Msg("Received message")

		switch message.Request {
		case tictacgo.JoinRoom:
			r, ok := s.rooms[message.Room]
			if !ok {
				log.Warn().Msg("Room does not exist")
				return
			}
			log.Debug().Str("ADDR", address)

			mark := "X"

			if len(r.players) == 1 {
				mark = "O"
			} else if len(r.players) > 1 {
				log.Info().Msg("Room full!")
				return
			}

			r.players[conn.RemoteAddr().String()] = newPlayer(mark, conn)
			mess := tictacgo.Message{
				Request: tictacgo.Update,
				Room:    message.Room,
				Player:  mark,
			}
			s.sendMessage(conn, mess)
			s.broadcastUpdates(mess.Room)
		case tictacgo.MakeMove: //MakeMove ROOM PLAYER MOVE
			room, ok := s.rooms[message.Room]
			if !ok {
				log.Warn().Msg("room '" + message.Room + "' does not exist")
				return
			}

			if message.Player == "" {
				log.Debug().Msg("Message Player does not exist")
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

func (s *server) sendMessage(conn net.Conn, message tictacgo.Message) {
	res := tictacgo.Response{
		Code:    tictacgo.Success,
		Message: message,
	}

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(res)
	if err != nil {
		return
	}
}

func (s *server) broadcastUpdates(roomStr string) {
	rooms := make([]string, len(s.rooms))

	i := 0
	for k := range s.rooms {
		rooms[i] = k
		i++
	}
	message := tictacgo.Message{
		Request: tictacgo.Update,
		Room:    roomStr,
		Move:    s.rooms[roomStr].game.GetTurn(),
		Turn:    s.rooms[roomStr].game.GetTurn(),
		Winner:  s.rooms[roomStr].game.GetWinner(),
		Board:   s.rooms[roomStr].game.GetBoard(),
	}

	for _, room := range s.rooms {
		if room.name == roomStr {
			for _, player := range room.players {
				s.sendMessage(player.connection, message)
			}
		}
	}
}

func (s *server) startTCPServer() {
	var err error
	if s.tcpListener, err = net.Listen("tcp", ":8080"); err != nil {
		log.Error().Err(err).Msg("listener failed to start")
	}

	defer func(listen net.Listener) {
		if err := listen.Close(); err != nil {
			log.Fatal().Err(err).Msg("deferred listener")
		}
	}(s.tcpListener)

	log.Info().Msg("Start tcp server on port 8080")
	func() {
		for {
			conn, err := s.tcpListener.Accept()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to accepted connection")
			}

			go s.handleConnection(conn)
		}
	}()
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	s := server{
		rooms: make(map[string]*room),
	}

	s.makeRoom("test13")
	s.makeRoom("yo")
	go s.startTCPServer()
	s.startRESTServer()
}
