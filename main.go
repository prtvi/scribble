package main

import (
	routes "scribble/routes"
	utils "scribble/utils"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Static("/public", "public")
	e.Renderer = utils.InitTemplates()

	ee := e.Group("", routes.Logger)

	ee.GET("/", routes.Welcome)

	ee.GET("/create-pool", routes.CreatePool)
	ee.POST("/create-pool", routes.CreatePoolLink)

	ee.GET("/app", routes.App)
	ee.POST("/app", routes.RegisterToPool)

	ee.GET("/ws", routes.HandlerWsConnection)

	e.GET("/api/get-clients-in-pool", routes.GetAllClientsInPool)
	e.GET("/api/start-game", routes.StartGame)

	e.Logger.Fatal(e.Start(":1323"))
}
