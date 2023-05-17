package tictacgo

import (
	"fmt"
	"strconv"
)

type Game struct {
	board  []string
	turn   string
	winner string
	moves  int
}

func NewGame() Game {
	g := &Game{
		board:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		turn:   "X",
		winner: "",
		moves:  0,
	}

	return *g
}

func (g *Game) GetBoard() []string {
	return g.board
}

func (g *Game) GetTurn() string {
	return g.turn
}

func (g *Game) GetWinner() string {
	return g.winner
}

func (g *Game) GetMoves() int {
	return g.moves
}

func (g *Game) HasWinner() bool {
	return g.winner != ""
}

func (g *Game) SetCell(cell int, value string) {
	g.board[cell] = value
}

func (g *Game) CheckWinner() bool {
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

func (g *Game) Move(cell string) bool {
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

func (g *Game) Print() {
	for i := 0; i < 9; i++ {
		fmt.Print(g.board[i] + " ")
		if (i+1)%3 == 0 {
			fmt.Print("\n")
		}
	}
}
