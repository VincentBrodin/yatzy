package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Game struct {
	Clients map[*Client]bool

	Broadcast chan []byte

	Register chan *Client

	Unregister chan *Client
}

func (g *Game) Start() {
	for {
		select {
		case client := <-g.Register:
			log.Println("Client connected")
			g.Clients[client] = true
		case client := <-g.Unregister:
			if _, ok := g.Clients[client]; ok {
				delete(g.Clients, client)
				close(client.Buf)
			}
		case message := <-g.Broadcast:
			for client := range g.Clients {
				select {
				case client.Buf <- message:
				default:
					close(client.Buf)
					delete(g.Clients, client)
				}
			}
		}
	}
}

func (g *Game) Serve(c echo.Context) error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by default (good for development)
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	client := &Client{g, conn, make(chan []byte, 256)}
	client.Game.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.Write()
	go client.Read()
	return nil
}

func NewGame() *Game {
	return &Game{
		make(map[*Client]bool),
		make(chan []byte),
		make(chan *Client),
		make(chan *Client),
	}
}
