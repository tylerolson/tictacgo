package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tylerolson/tictacgo/server"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	s := server.Server{
		Rooms: make(map[string]*server.Room),
	}

	s.MakeRoom("test13")
	s.MakeRoom("yo")
	go s.StartTCPServer()
	s.StartRESTServer()

}
