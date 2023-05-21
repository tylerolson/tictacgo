package main

import (
	"fmt"
	"github.com/tylerolson/tictacgo/client"
	"log"
	"time"
)

func main() {
	c := client.NewClient("X")
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
