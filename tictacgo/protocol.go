package tictacgo

const (
	CREATE_ROOM = "CREATE_ROOM"
	JOIN_ROOM   = "JOIN_ROOM"
	MAKE_MOVE   = "MAKE_MOVE"
	UPDATE      = "UPDATE"
)

type Message struct {
	Request string   `json:"request"`
	Room    string   `json:"room"`
	Player  string   `json:"player"`
	Move    string   `json:"move"`
	Turn    string   `json:"turn"`
	Winner  string   `json:"winner"`
	Board   []string `json:"board"`
}
