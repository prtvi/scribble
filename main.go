package main

import (
	socket "scribble/socket"
	utils "scribble/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/public", "public")
	e.Renderer = utils.InitTemplates()

	go socket.Maintainer()

	socket.DebugMode()

	ee := e.Group("", socket.Logger)

	ee.GET("/", socket.Welcome)

	ee.GET("/create-pool", socket.CreatePool)
	ee.POST("/create-pool", socket.CreatePoolLink)

	ee.GET("/app", socket.App)
	ee.POST("/app", socket.RegisterToPool)

	ee.GET("/ws", socket.HandlerWsConnection)

	e.Logger.Fatal(e.Start(":1323"))
}
