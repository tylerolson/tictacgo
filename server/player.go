package server

import "net"

type Player struct {
	mark       string
	connection net.Conn
}

func NewPlayer(mark string, connection net.Conn) Player {
	return Player{
		mark:       mark,
		connection: connection,
	}
}
