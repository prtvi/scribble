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

	ee := e.Group("", routes.Middleware)

	ee.GET("/welcome", routes.Welcome)

	ee.GET("/app", routes.App)
	ee.POST("/app", routes.RegisterToPool)

	ee.GET("/create-pool", routes.CreatePool)
	ee.POST("/create-pool", routes.CreatePoolLink)

	ee.GET("/ws", routes.HandlerWsConnection)

	e.Logger.Fatal(e.Start(":1323"))
}
