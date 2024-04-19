package server

import (
	"github.com/rs/zerolog/log"
	"github.com/tylerolson/tictacgo"
)

type Room struct {
	name    string
	game    *tictacgo.Game
	started bool
	players map[string]*Player
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		game:    tictacgo.NewGame(),
		started: false,
		players: make(map[string]*Player),
	}
}

func (r Room) PrintRoom() {
	log.Info().
		Str("Name", r.name).
		Interface("game", r.game).
		Interface("Players", r.players).
		Msg("Printed room")
}
