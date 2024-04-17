package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	tictacgo "github.com/tylerolson/tictacgo"
)

type player struct {
	mark       string
	connection net.Conn
}

type Room struct {
	name    string
	game    *tictacgo.Game
	started bool
	players map[string]player
}

type Server struct {
	Rooms       map[string]*Room
	restRouter  *gin.Engine
	tcpListener net.Listener
}

func printRoom(room Room) {
	log.Info().
		Str("Name", room.name).
		Interface("game", room.game).
		Interface("Players", room.players).
		Msg("Printed room")
}

func newPlayer(mark string, connection net.Conn) player {
	return player{
		mark:       mark,
		connection: connection,
	}
}

func (s *Server) removePlayerFromAll(addr string) {
	for _, room := range s.Rooms {
		delete(room.players, addr)
	}
}

// rooms

func (s *Server) MakeRoom(name string) {
	g := tictacgo.NewGame()
	r := Room{
		name:    name,
		game:    g,
		started: false,
		players: make(map[string]player),
	}
	s.Rooms[name] = &r

	log.Info().Msg("Created room " + name)
}

// rest

func (s *Server) getRoomsRoute(c *gin.Context) {
	rooms := make([]RoomResponse, 0)

	for _, v := range s.Rooms {
		room := RoomResponse{
			Name: v.name,
			Size: len(v.players),
		}
		rooms = append(rooms, room)
	}

	response := Response{
		Type:    GetRoom,
		Content: rooms,
	}

	c.IndentedJSON(http.StatusOK, response)
}

func (s *Server) makeRoomRoute(c *gin.Context) {
	var rawContent json.RawMessage
	request := Request{
		Content: &rawContent,
	}

	if err := c.BindJSON(&request); err != nil {
		log.Fatal().Err(err).Msg("Failed to read POST")
	}

	var content RoomContent
	err := json.Unmarshal(rawContent, &content)
	if err != nil {
		fmt.Println("Error unmarshalling JoinRoomContent:", err)
		return
	}

	log.Info().Str("name", content.RoomName).Send()

	s.MakeRoom(content.RoomName)

	c.IndentedJSON(http.StatusCreated, content.RoomName)
}

func (s *Server) StartRESTServer() {
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

func (s *Server) handleConnection(conn net.Conn) {
	for {
		var rawContent json.RawMessage

		request := Request{
			Content: &rawContent,
		}
		address := conn.RemoteAddr().String()

		if err := json.NewDecoder(conn).Decode(&request); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, syscall.ECONNRESET) {
				log.Info().Msg("Client disconnected")
				s.removePlayerFromAll(address)
				break
			}
			log.Fatal().Err(err).Msg("Couldn't decode request")
		}

		log.Info().Msg("Got request of type " + string(request.Type))

		switch request.Type {
		case JoinRoom:
			var content RoomContent
			err := json.Unmarshal(rawContent, &content)
			if err != nil {
				fmt.Println("Error unmarshalling JoinRoomContent:", err)
				return
			}

			room, ok := s.Rooms[content.RoomName]
			if !ok {
				log.Warn().Msg("Room does not exist")
				return
			}
			log.Debug().Str("ADDR", address).Send()

			mark := "X"

			if len(room.players) == 1 {
				mark = "O"
				room.started = true
			} else if len(room.players) > 1 {
				log.Info().Msg("Room full!")
				return
			}

			room.players[conn.RemoteAddr().String()] = newPlayer(mark, conn)
			responseContent := AssignMarkContent{
				Room:   content.RoomName,
				Player: mark,
			}
			s.sendMessage(conn, responseContent, AssignMark)
			s.broadcastUpdates(content.RoomName)
		case MakeMove:
			var content MakeMoveContent
			err := json.Unmarshal(rawContent, &content)
			if err != nil {
				log.Err(err).Msg("Error unmarshalling MakeMoveContent:")
				return
			}

			room, ok := s.Rooms[content.Room]
			if !ok {
				log.Warn().Msg("room '" + content.Room + "' does not exist")
				return
			}

			if content.Player == "" {
				log.Warn().Msg("Content Player does not exist")
				return
			}

			log.Info().Interface("content", content.Room).Send()
			printRoom(*room)

			if len(room.players) >= 2 { // setting this greater or equal for now, should be equal
				if room.game.Turn == content.Player {
					room.game.Move(content.Move)
					log.Debug().Msg("Made move " + content.Move + "in room " + room.name)
				}
			}

			s.broadcastUpdates(content.Room)
		}
	}
}

func (s *Server) sendMessage(conn net.Conn, content any, responseType ResponseType) {
	res := Response{
		Type:    responseType,
		Content: content,
	}

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(res)
	if err != nil {
		return
	}
}

func (s *Server) broadcastUpdates(roomStr string) {
	rooms := make([]string, len(s.Rooms))

	i := 0
	for k := range s.Rooms {
		rooms[i] = k
		i++
	}

	for _, room := range s.Rooms {
		if room.name == roomStr {
			for _, player := range room.players {
				s.sendMessage(player.connection, UpdateGameContent{
					Game:    *room.game,
					Started: room.started,
				}, UpdateGame)
			}
		}
	}
}

func (s *Server) StartTCPServer() {
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

// func main() {
// 	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
// 	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// 	s := server{
// 		rooms: make(map[string]*room),
// 	}

// 	s.makeRoom("test13")
// 	s.makeRoom("yo")
// 	go s.startTCPServer()
// 	s.startRESTServer()
// }
