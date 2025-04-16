package main

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/VincentBrodin/yatzy/backend/ws"
	"github.com/google/uuid"
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
	game.Register(0, Join)
	game.Register(1, Roll)
	game.Register(2, Select)
	go game.Start()

	e.GET("/ws", game.Serve)
	e.Logger.Fatal(e.Start(":3000"))
}

func Join(game *ws.Game, in *ws.Packet) error {
	input := struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	}{}
	if err := json.Unmarshal(in.Message, &input); err != nil {
		return err
	}

	id, err := uuid.Parse(input.Id)
	if err == nil {
		// We got a bad id
		if _, ok := game.State.Players[id]; !ok {
			fmt.Println("Invalid id")
			json, err := json.Marshal(struct {
				Error string `json:"error"`
			}{Error: "Invalid id"})
			if err != nil {
				return err
			}

			in.Client.Send(&ws.Packet{
				CallId:  0,
				Client:  nil,
				Message: json,
			})
			return nil
		}
	} else {
		id = uuid.New()
		json, err := json.Marshal(struct {
			Id string `json:"id"`
		}{Id: id.String()})

		if err != nil {
			return err
		}
		game.State.Players[id] = input.Username

		// Send uuid to player
		in.Client.Send(&ws.Packet{
			CallId:  0,
			Client:  nil,
			Message: json})
	}
	in.Client.Id = id
	game.BroadcastState(1)

	return nil
}

func Roll(game *ws.Game, in *ws.Packet) error {
	for _, die := range game.State.Dice {
		if !die.Selected {
			die.Value = rand.Intn(6) + 1
		}
	}

	game.BroadcastState(2)
	return nil
}

func Select(game *ws.Game, in *ws.Packet) error {
	newState := struct {
		Index    int  `json:"index"`
		Selected bool `json:"selected"`
	}{}
	if err := json.Unmarshal(in.Message, &newState); err != nil {
		return err
	}

	game.State.Dice[newState.Index].Selected = newState.Selected

	game.BroadcastState(1)

	return nil
}
