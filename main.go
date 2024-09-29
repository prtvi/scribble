package main

import (
	"fmt"
	socket "scribble/socket"
	utils "scribble/utils"

	"github.com/labstack/echo/v4"
)

func panicRecover() {
	if err := recover(); err == nil {
		return
	}

	utils.Cp("red", "restarting service ...")
	main()
}

func main() {
	defer panicRecover()
	go socket.Maintainer()

	e := echo.New()

	e.Static("/public", "public")
	e.Static("/scribble/public", "public")
	e.Renderer = utils.InitTemplates()

	ee := e.Group("/scribble", socket.Logger)

	ee.GET("", socket.Index)
	ee.GET("/", socket.Index)

	ee.GET("/create-room", socket.CreateRoomForm)
	ee.POST("/create-room", socket.CreateRoom)

	ee.GET("/app", socket.JoinPool)
	ee.POST("/app", socket.EnterPool)

	ee.GET("/ws", socket.WsConnect)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", utils.GetPort())))
}
