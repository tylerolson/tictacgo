package tictacgo

import (
	"fmt"
	"strconv"
)

type Game struct {
	Board  []string
	Turn   string
	Winner string
	Moves  int
}

func NewGame() *Game {
	return &Game{
		Board:  []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
		Turn:   "X",
		Winner: "",
		Moves:  0,
	}
}

func (g *Game) SetGame(game Game) {
	g.Board = game.Board
	g.Turn = game.Turn
	g.Winner = game.Winner
	g.Moves = game.Moves
}

func (g *Game) HasWinner() bool {
	return g.Winner != ""
}

func (g *Game) SetCell(cell int, value string) {
	g.Board[cell] = value
}

func (g *Game) CheckWinner() bool {
	// horizontal

	if g.Moves >= 9 {
		g.Winner = "tie"
		return true
	}

	if g.Board[0] == g.Board[1] && g.Board[1] == g.Board[2] {
		g.Winner = g.Board[0]
		return true
	}

	if g.Board[3] == g.Board[4] && g.Board[4] == g.Board[5] {
		g.Winner = g.Board[3]
		return true
	}

	if g.Board[6] == g.Board[7] && g.Board[7] == g.Board[8] {
		g.Winner = g.Board[6]
		return true
	}

	// vertical

	if g.Board[0] == g.Board[3] && g.Board[3] == g.Board[6] {
		g.Winner = g.Board[0]
		return true
	}

	if g.Board[1] == g.Board[4] && g.Board[4] == g.Board[7] {
		g.Winner = g.Board[1]
		return true
	}

	if g.Board[2] == g.Board[5] && g.Board[5] == g.Board[8] {
		g.Winner = g.Board[2]
		return true
	}

	// diagonal

	if g.Board[0] == g.Board[4] && g.Board[4] == g.Board[8] {
		g.Winner = g.Board[0]
		return true
	}

	if g.Board[2] == g.Board[4] && g.Board[4] == g.Board[6] {
		g.Winner = g.Board[2]
		return true
	}

	return false
}

func (g *Game) Move(cell string) bool {
	if g.CheckWinner() {
		return false // throw err maybe?
	}

	cellInt, _ := strconv.Atoi(cell)
	cellInt--

	if g.Board[cellInt] == cell {
		g.Board[cellInt] = g.Turn
	} else {
		return false
	}

	if g.Turn == "X" {
		g.Turn = "O"
	} else {
		g.Turn = "X"
	}

	g.Moves++

	g.CheckWinner()

	return true
}

func (g *Game) Print() {
	for i := 0; i < 9; i++ {
		fmt.Print(g.Board[i] + " ")
		if (i+1)%3 == 0 {
			fmt.Print("\n")
		}
	}
}
