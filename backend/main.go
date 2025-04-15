package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"

	"github.com/VincentBrodin/yatzy/backend/ws"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	game := ws.NewGame()
	game.Register(1, message)
	go game.Start()

	e.GET("/ws", game.Serve)
	e.Logger.Fatal(e.Start(":3000"))
}

func message(game *ws.Game, in *ws.Packet) error {
	log.Printf("msg: %s\n", in.Message)
	dice := make([]byte, 0) // 32bit int * 5
	for range 5 {
		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.BigEndian, uint32(rand.Intn(6)+1)); err != nil {
			return err
		}
		dice = append(dice, buf.Bytes()...)
	}
	out := &ws.Packet{
		CallId:  1,
		Client:  in.Client,
		Message: dice,
	}
	in.Client.Send(out)
	return nil
}
