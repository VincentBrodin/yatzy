package main

import (
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
	go game.Start()
	e.GET("/ws", game.Serve)
	e.Logger.Fatal(e.Start(":3000"))
}
