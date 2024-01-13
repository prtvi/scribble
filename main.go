package main

import (
	socket "scribble/socket"
	utils "scribble/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	defer func() {
		// TODO - gotta test recover
		if recover() == nil {
			return
		}

		utils.Cp("red", "panic occurred! recovering ...")
	}()

	isDebugEnv := utils.LoadAndGetEnv()
	socket.InitDebugEnv(isDebugEnv)

	e := echo.New()
	e.Static("/public", "public")
	e.Renderer = utils.InitTemplates()
	ee := e.Group("", socket.Logger)

	ee.GET("/", socket.Index)

	ee.GET("/create-room", socket.CreateRoomForm)
	ee.POST("/create-room", socket.CreateRoom)

	ee.GET("/app", socket.JoinPool)
	ee.POST("/app", socket.EnterPool)

	ee.GET("/ws", socket.WsConnect)

	e.Logger.Fatal(e.Start(":1323"))
}
