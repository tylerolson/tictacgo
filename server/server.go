package server

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server struct {
	Rooms       map[string]*Room
	restRouter  *gin.Engine
	tcpListener net.Listener
}

func NewServer() *Server {
	return &Server{
		Rooms: make(map[string]*Room),
	}
}
func (s *Server) removePlayerFromAll(addr string) {
	for _, room := range s.Rooms {
		delete(room.players, addr)
	}
}

func (s *Server) MakeRoom(name string) {
	s.Rooms[name] = NewRoom(name)

	log.Info().Str("name", name).Msg("Created room")
}

func (s *Server) GetRoom(name string) *Room {
	room, ok := s.Rooms[name]
	if !ok {
		log.Warn().Str("name", name).Msg("Room does not exist")
		return nil
	}
	return room
}

// rest

func (s *Server) getRooms(c *gin.Context) {
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

func (s *Server) postRooms(c *gin.Context) {
	var rawContent json.RawMessage
	request := Request{
		Content: &rawContent,
	}

	if err := c.BindJSON(&request); err != nil {
		log.Err(err).Msg("Failed to read JoinRoomRequest")
		return
	}

	var content RoomContent
	if err := json.Unmarshal(rawContent, &content); err != nil {
		log.Err(err).Msg("Failed to unmarshall JoinRoomContent")
		return
	}

	s.MakeRoom(content.Room)
	c.IndentedJSON(http.StatusCreated, content.Room)
}

func (s *Server) StartRESTServer() {
	s.restRouter = gin.New()
	s.restRouter.Use(gin.Logger())
	s.restRouter.GET("/rooms", s.getRooms)
	s.restRouter.POST("/rooms", s.postRooms)
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
				log.Info().Str("address", address).Msg("Client disconnected")
				s.removePlayerFromAll(address)
				break
			}
			log.Err(err).Msg("Failed to read request")
			return
		}

		log.Info().Str("type", string(request.Type)).Str("address", address).Msg("Got request")

		switch request.Type {
		case JoinRoom:
			var content RoomContent
			err := json.Unmarshal(rawContent, &content)
			if err != nil {
				log.Err(err).Msg("Failed to unmarshall JoinRoomContent")
				return
			}

			room := s.GetRoom(content.Room)
			if room == nil {
				return
			}

			mark := "X"

			if len(room.players) == 1 {
				mark = "O"
				room.started = true
			} else if len(room.players) > 1 {
				log.Info().Str("name", content.Room).Msg("Room is full")
				return
			}

			room.players[address] = NewPlayer(mark, conn)
			responseContent := AssignMarkContent{
				Room:   content.Room,
				Player: mark,
			}
			s.sendMessage(conn, AssignMark, responseContent)
			s.broadcastUpdates(content.Room)

			log.Info().Str("address", address).Str("mark", mark).Str("room", content.Room).Msg("Player joined room")
		case MakeMove:
			var content MakeMoveContent
			err := json.Unmarshal(rawContent, &content)
			if err != nil {
				log.Err(err).Msg("Failed to unmarshall MakeMoveContent")
				return
			}

			room := s.GetRoom(content.Room)
			if room == nil {
				return
			}

			if len(room.players) >= 2 { // setting this greater or equal for now, should be equal
				if room.game.Turn == content.Player {
					room.game.Move(content.Move)
					log.Info().Str("move", content.Move).Str("room", room.name).Msg("Made move")
				}
			}

			s.broadcastUpdates(content.Room)
		}
	}
}

func (s *Server) sendMessage(conn net.Conn, responseType ResponseType, content any) {
	res := Response{
		Type:    responseType,
		Content: content,
	}

	encoder := json.NewEncoder(conn)
	err := encoder.Encode(res)
	if err != nil {
		log.Err(err).Msg("Failed to encode message")
		return
	}
}

func (s *Server) broadcastUpdates(roomName string) {
	for _, room := range s.Rooms {
		if room.name == roomName {
			for _, player := range room.players {
				s.sendMessage(player.connection, UpdateGame, UpdateGameContent{
					Game:    *room.game,
					Started: room.started,
				})
			}
		}
	}
}

func (s *Server) StartTCPServer() {
	var err error
	if s.tcpListener, err = net.Listen("tcp", ":8080"); err != nil {
		log.Error().Err(err).Msg("Failed to start listener")
	}

	defer func(listen net.Listener) {
		if err := listen.Close(); err != nil {
			log.Fatal().Err(err).Msg("deferred listener")
		}
	}(s.tcpListener)

	log.Info().Msg("Started tcp server on port 8080")
	func() {
		for {
			conn, err := s.tcpListener.Accept()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to accept connection")
			}

			go s.handleConnection(conn)
		}
	}()
}
