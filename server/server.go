package server

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"syscall"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Rooms       map[string]*Room
	restRouter  *http.ServeMux
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

func (s *Server) getRooms(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET /rooms request")
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) postRooms(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST /rooms request")

	var rawContent json.RawMessage
	request := Request{
		Content: &rawContent,
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Err(err).Msg("Failed to read JoinRoomRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var content RoomContent
	if err := json.Unmarshal(rawContent, &content); err != nil {
		log.Err(err).Msg("Failed to unmarshall JoinRoomContent")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.MakeRoom(content.Room)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) StartRESTServer() {
	s.restRouter = http.NewServeMux()
	s.restRouter.HandleFunc("GET /rooms", s.getRooms)
	s.restRouter.HandleFunc("POST /rooms", s.postRooms)
	err := http.ListenAndServe(":8081", s.restRouter)
	log.Info().Msg("yo")
	if err != nil {
		log.Err(err).Msg("Failed to start REST server")
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
