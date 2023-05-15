package game

import "strconv"

type TicTacGo struct {
	board  []string
	turn   string
	winner string
	moves  int
}

func NewGame() TicTacGo {
	g := &TicTacGo{
		board:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		turn:   "X",
		winner: "",
		moves:  0,
	}

	return *g
}

func (g *TicTacGo) GetBoard() []string {
	return g.board
}

func (g *TicTacGo) GetTurn() string {
	return g.turn
}

func (g *TicTacGo) GetWinner() string {
	return g.winner
}

func (g *TicTacGo) GetMoves() int {
	return g.moves
}

func (g *TicTacGo) CheckWinner() bool {
	// horizontal

	if g.moves >= 9 {
		g.winner = "tie"
		return true
	}

	if g.board[0] == g.board[1] && g.board[1] == g.board[2] {
		g.winner = g.board[0]
		return true
	}

	if g.board[3] == g.board[4] && g.board[4] == g.board[5] {
		g.winner = g.board[3]
		return true
	}

	if g.board[6] == g.board[7] && g.board[7] == g.board[8] {
		g.winner = g.board[6]
		return true
	}

	//vertical

	if g.board[0] == g.board[3] && g.board[3] == g.board[6] {
		g.winner = g.board[0]
		return true
	}

	if g.board[1] == g.board[4] && g.board[4] == g.board[7] {
		g.winner = g.board[1]
		return true
	}

	if g.board[2] == g.board[5] && g.board[5] == g.board[8] {
		g.winner = g.board[2]
		return true
	}

	//diagonal

	if g.board[0] == g.board[4] && g.board[4] == g.board[8] {
		g.winner = g.board[0]
		return true
	}

	if g.board[2] == g.board[4] && g.board[4] == g.board[6] {
		g.winner = g.board[2]
		return true
	}

	return false
}

func (g *TicTacGo) Move(cell string) bool {
	if g.CheckWinner() {
		return false //throw err maybe?
	}

	cellInt, _ := strconv.Atoi(cell)
	cellInt--

	if g.board[cellInt] == cell {
		g.board[cellInt] = g.turn
	} else {
		return false
	}

	if g.turn == "X" {
		g.turn = "O"
	} else {
		g.turn = "X"
	}

	g.moves++

	return true
}
