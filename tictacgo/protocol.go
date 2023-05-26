package tictacgo

const (
	CreateRoom = "CreateRoom"
	JoinRoom   = "JoinRoom"
	MakeMove   = "MakeMove"
	Update     = "Update"
	Success    = "Success"
	Fail       = "Fail"
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

type Response struct {
	Code string `json:"code"`
	Message
}
