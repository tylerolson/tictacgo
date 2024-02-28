# TicTacGo

A terminal application to play local or online tic-tac-toe games.
Uses Bubble Tea for view/state management, Gin for REST, and Go std for TCP server.

### Getting Started

Clone the repository. 

To run the client navigate to the root directory and run:
```bash
go run .
```

To start the server navigate to `/server` and run:
```bash
go run .
```

### TODO

* Join random room button? Join first available room?
* ~~Better server side cleanup after game is over~~
* ~~Change create/join rooms to join to a lobby screen and wait for others~~

### Libraries

* [Bubble Tea](https://github.com/charmbracelet/bubbletea)
* [Bubbles](https://github.com/charmbracelet/bubbles)
* [Lip Gloss](https://github.com/charmbracelet/lipgloss)
* [Gin](https://github.com/gin-gonic/gin)
* [zerolog](https://github.com/rs/zerolog)
