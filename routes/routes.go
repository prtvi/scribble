package routes

import (
	"fmt"
	"net/http"
	socket "scribble/socket"
	utils "scribble/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// map of {poolId: pool}
var Hub = map[string]*socket.Pool{}

// / middleware
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dt := time.Now().String()[0:19]
		reqMethod := c.Request().Method

		var color string
		if reqMethod == "GET" {
			color = "green"
		} else if reqMethod == "POST" {
			color = "cyan"
		}

		utils.Cp(color, fmt.Sprintf("%s: %s  %s", reqMethod, utils.Cs("white", fmt.Sprintf("%s", c.Request().URL)), utils.Cs(color, dt)))

		return next(c)
	}
}

// GET /
func Welcome(c echo.Context) error {
	return c.Render(http.StatusOK, "welcome", nil)
}

// GET /ws?poolId=234bkj&clientId=123123&clientName=joy
func HandlerWsConnection(c echo.Context) error {
	// handle socket connections for the pools

	// get the poolId from query params
	poolId := c.QueryParam("poolId")

	// register connection
	return socket.ServeWs(Hub[poolId], c.Response().Writer, c.Request())
}

// GET /api/get-clients-in-pool?poolId=123jisd
func GetAllClientsInPool(c echo.Context) error {
	// returns all the clients (name and color properties) in the pool
	type clientInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	poolId := c.QueryParam("poolId")
	clientInfoList := make([]clientInfo, 0)

	pool, ok := Hub[poolId]
	if !ok {
		return c.JSON(http.StatusOK, clientInfoList)
	}

	for client := range pool.Clients {
		clientInfoList = append(clientInfoList, clientInfo{
			ID:    client.ID,
			Name:  client.Name,
			Color: client.Color,
		})
	}

	return c.JSON(http.StatusOK, clientInfoList)
}
