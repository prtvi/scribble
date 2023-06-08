package main

import (
	socket "scribble/socket"
	utils "scribble/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	utils.LoadEnv()
	socket.DebugMode()

	go socket.Maintainer()

	e := echo.New()
	e.Static("/public", "public")
	e.Renderer = utils.InitTemplates()

	ee := e.Group("", socket.Logger)

	ee.GET("/", socket.Welcome)

	ee.GET("/create-room", socket.CreateRoom)
	ee.POST("/create-room", socket.CreateRoomLink)

	ee.GET("/app", socket.App)
	ee.POST("/app", socket.RegisterToPool)

	ee.GET("/ws", socket.HandlerWsConnection)

	e.Logger.Fatal(e.Start(":1323"))
}
