package main

import (
	"fmt"
	"github.com/tylerolson/tictacgo/tictacgo"
	"log"
	"time"
)

func _main() {
	fmt.Println("Pick a player to be (X/O)")
	playerChoice := ""
	_, err := fmt.Scanln(&playerChoice)
	if err != nil {
		return
	}

	c := tictacgo.NewClient()
	c.SetPlayer(playerChoice)
	c.EstablishConnection("localhost:8080")
	c.CreateRoom("yo")

	game := c.GetGame()

	for !game.HasWinner() {
		game.Print()
		fmt.Println("make a move")
		fmt.Println("current move is " + game.GetTurn())
		cell := ""
		_, err := fmt.Scanln(&cell)
		if err != nil {
			log.Fatalln(err)
		}
		c.MakeMove(cell)

		time.Sleep(50 * time.Millisecond)
	}
}
