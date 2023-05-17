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

	g := c.GetGame()
	g.Print()

	for !g.HasWinner() {
		fmt.Println("make a move")
		cell := ""
		_, err := fmt.Scanln(&cell)
		if err != nil {
			log.Fatalln(err)
		}
		c.MakeMove(cell)
		g = c.GetGame()

		time.Sleep(50 * time.Millisecond)
		g.Print()
	}
}
