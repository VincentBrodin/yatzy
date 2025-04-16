package main

import (
	"encoding/json"
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
	game.Register(1, Roll)
	game.Register(2, Select)
	go game.Start()

	e.GET("/ws", game.Serve)
	e.Logger.Fatal(e.Start(":3000"))
}

func Roll(game *ws.Game, in *ws.Packet) error {
	log.Printf("msg: %s\n", in.Message)
	for _, die := range game.State.Dice {
		if !die.Selected {
			die.Value = rand.Intn(6) + 1
		}
	}

	json, err := json.Marshal(game.State)
	if err != nil {
		return err
	}
	game.Broadcast(&ws.Packet{
		CallId:  1,
		Client:  nil,
		Message: json,
	})

	return nil
}

func Select(game *ws.Game, in *ws.Packet) error {
	log.Printf("msg: %s\n", in.Message)
	newState := struct {
		Index    int  `json:"index"`
		Selected bool `json:"selected"`
	}{}
	if err := json.Unmarshal(in.Message, &newState); err != nil {
		return err
	}

	game.State.Dice[newState.Index].Selected = newState.Selected

	json, err := json.Marshal(game.State)
	if err != nil {
		return err
	}

	game.Broadcast(&ws.Packet{
		CallId:  2,
		Client:  nil,
		Message: json,
	})
	return nil
}
