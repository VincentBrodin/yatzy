package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler func(g *Game, p *Packet) error

type Game struct {
	State *State

	clients    map[*Client]bool
	packets    chan *Packet
	register   chan *Client
	unregister chan *Client
	handlers   map[uint32]Handler
}

func (g *Game) Start() {
	for {
		select {
		case client := <-g.register:
			log.Println("Client connected")
			g.clients[client] = true
		case client := <-g.unregister:
			if _, ok := g.clients[client]; ok {
				log.Println("Client disconnected")
				delete(g.clients, client)
				close(client.packets)
			}
		case packet := <-g.packets:
			handler, ok := g.handlers[packet.CallId]
			if ok {
				go handler(g, packet)
			} else {
				log.Printf("No handler for %d\n", packet.CallId)
			}
		}
	}
}

func (g *Game) Register(callId uint32, handler Handler) {
	g.handlers[callId] = handler
}

func (g *Game) Broadcast(packet *Packet) {
	for client := range g.clients {
		client.Send(packet)
	}
}

func (g *Game) BroadcastState(callId uint32) {
	json, err := json.Marshal(g.State)
	if err != nil {
		log.Println(err)
		panic(err)
		return
	}

	g.Broadcast(&Packet{
		CallId:  callId,
		Client:  nil,
		Message: json,
	})
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
	client := &Client{uuid.Max, g, conn, make(chan *Packet, 256)}
	client.Game.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.write()
	go client.read()
	json, err := json.Marshal(g.State)
	if err != nil {
		return err
	}
	client.Send(&Packet{
		CallId:  1,
		Client:  nil,
		Message: json,
	})
	return nil
}

func NewGame() *Game {
	return &Game{
		NewState(),
		make(map[*Client]bool),
		make(chan *Packet),
		make(chan *Client),
		make(chan *Client),
		make(map[uint32]Handler),
	}
}
