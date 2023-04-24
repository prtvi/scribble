package main

import (
	"fmt"
	"net/http"
	socket "scribble/socket"
	utils "scribble/utils"

	"github.com/labstack/echo/v4"
)

func serveWs(pool *socket.Pool, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("WebSocket Endpoint Hit by javascript")

	conn, err := socket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &socket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()

	return nil
}

func setupRoutes(e *echo.Echo) {
	pool := socket.NewPool()
	go pool.Start()

	e.GET("/app", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	e.GET("/ws", func(c echo.Context) error {
		return serveWs(pool, c.Response().Writer, c.Request())
	})
}

func main() {
	e := echo.New()
	e.Static("/public", "public")
	e.Renderer = utils.T

	setupRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
