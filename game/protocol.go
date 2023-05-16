package game

const (
	CREATE_ROOM = "CREATE_ROOM"
	JOIN_ROOM   = "CREATE_ROOM"
)

type Message struct {
	Request  string
	RoomName string
}
