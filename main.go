package main

import (
	"fmt"
	"net/http"
	socket "scribble/socket"
	utils "scribble/utils"

	"github.com/labstack/echo/v4"
)

var Pools = map[string]*socket.Pool{}

func serveWs(pool *socket.Pool, w http.ResponseWriter, r *http.Request) error {
	clientId := r.URL.Query().Get("clientId")

	conn, err := socket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &socket.Client{
		ID:   clientId,
		Conn: conn,
		Pool: pool,
	}

	fmt.Println("New client:", client.ID)

	pool.Register <- client
	client.Read()

	return nil
}

func main() {
	e := echo.New()
	e.Static("/public", "public")
	e.Renderer = utils.InitTemplates()

	e.GET("/create-pool", func(c echo.Context) error {
		return c.Render(http.StatusOK, "createPool", map[string]any{
			"Link": "",
		})
	})

	e.POST("/create-pool", func(c echo.Context) error {
		// create a new pool with an id
		poolId := utils.GenerateUUID()

		pool := socket.NewPool(poolId)
		Pools[poolId] = pool
		go pool.Start()

		// generate link to join the pool
		link := "/app?join=" + poolId

		fmt.Println("Pool link:", link)

		// send the link for the same
		return c.Render(http.StatusOK, "createPool", map[string]any{
			"Link": link,
		})
	})

	e.GET("/ws/:poolId", func(c echo.Context) error {
		poolId := c.Param("poolId")
		clientId := c.QueryParam("clientId")
		fmt.Println(poolId, clientId)

		if pool, ok := Pools[poolId]; ok {
			return serveWs(pool, c.Response().Writer, c.Request())
		}

		return c.JSON(http.StatusInternalServerError, `"msg":"some error"`)
	})

	e.GET("/app", func(c echo.Context) error {
		poolId := c.QueryParam("join")
		fmt.Println(poolId)

		return c.Render(http.StatusOK, "app", map[string]any{
			"Word": poolId,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
