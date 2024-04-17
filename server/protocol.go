package server

import "github.com/tylerolson/tictacgo"

type RequestType string

const (
	MakeRoom RequestType = "MakeRoom"
	JoinRoom RequestType = "JoinRoom"
	MakeMove RequestType = "MakeMove"
)

type Request struct {
	Type    RequestType `json:"requesttype"`
	Content interface{}
}

type RoomContent struct {
	RoomName string `json:"room"`
}

type MakeMoveContent struct {
	Room   string `json:"room"`
	Move   string `json:"move"`
	Player string `json:"player"`
}

type ResponseType string

const (
	GetRoom    ResponseType = "GetRoom"
	AssignMark ResponseType = "AssignMark"
	UpdateGame ResponseType = "UpdateGame"
)

type Response struct {
	Type    ResponseType `json:"type"`
	Content any          `json:"content"`
}

type RoomResponse struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

type AssignMarkContent struct {
	Room   string `json:"room"`
	Player string `json:"player"`
}

type UpdateGameContent struct {
	Game    tictacgo.Game `json:"game"`
	Started bool          `json:"started"`
}
