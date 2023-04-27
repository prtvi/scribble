package routes

import (
	"fmt"
	"net/http"
	socket "scribble/socket"
	utils "scribble/utils"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// map of {poolId: pool}
var Hub = map[string]*socket.Pool{}

// middleware
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dt := time.Now().String()[0:19]
		utils.Cp("green", fmt.Sprintf("%s: %s  %s", c.Request().Method, utils.Cs("white", fmt.Sprintf("%s", c.Request().URL)), utils.Cs("green", dt)))
		return next(c)
	}
}

// GET /welcome
// render welcome page
func Welcome(c echo.Context) error {
	return c.Render(http.StatusOK, "welcome", nil)
}

// GET /app
// if /app?join=sfds, then render the playing area
// if /app          , then render message
func App(c echo.Context) error {
	poolId := c.QueryParam("join")

	// if poolId is empty then do not render any forms, just display message
	if poolId == "" {
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Hi there, are you lost?!",
		})
	}

	// check if pool exists, verify if it exists then render the forms/message accordingly
	pool, ok := Hub[poolId]
	if !ok {
		// if not then do not render both forms and display message
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Pool expired or non-existent!",
		})
	}

	// if pool exists, get its capacity and curr size
	poolCap := pool.Capacity
	poolCurrSizePlus1 := len(pool.Clients) + 1

	if poolCurrSizePlus1 > poolCap {
		// if poolCurrSizePlus1 is greater than capacity then do not render both forms and display message
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Your party is full!",
		})
	}

	// else render "RegisterToPool" form
	return c.Render(http.StatusOK, "app", map[string]any{
		"RegisterToPool": true,
		"ConnectSocket":  false,
		"PoolId":         poolId, // hidden in form
	})
}

// POST /app
// on post request made to this route to capture clientName from "RegisterToPool" post form
func RegisterToPool(c echo.Context) error {
	poolId := c.FormValue("poolId")
	clientName := c.FormValue("clientName")

	// extra check to prevent user from joining any random pool which does not exist
	if _, ok := Hub[poolId]; !ok {
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Pool expired or non-existent!",
		})
	}

	// render ConnectSocket form to establish socket connection
	// socket connection will start only if "ConnectSocket" form is rendered
	return c.Render(http.StatusOK, "app", map[string]any{
		"RegisterToPool": false,
		"ConnectSocket":  true,
		"Message":        "",
		"PoolId":         poolId,     // hidden in form
		"ClientName":     clientName, // hidden in form
	})
}

// GET /create-pool
// render a form to create a new pool
func CreatePool(c echo.Context) error {
	return c.Render(http.StatusOK, "createPool", map[string]any{
		"Link": "",
	})
}

// POST /create-pool
// on post request to this route, create a new pool, start listening to connections on that pool, render the link to join this pool
func CreatePoolLink(c echo.Context) error {
	// get the pool capacity from form input
	capacity, _ := strconv.Atoi(c.FormValue("capacity"))
	utils.Cp("yellow", "Pool capacity:", utils.Cs("white", c.FormValue("capacity")))

	// create a new pool with an uuid
	poolId := utils.GenerateUUID()
	pool := socket.NewPool(poolId, capacity)

	// append to global Hub map, and start listening to pool connections
	Hub[poolId] = pool
	go pool.Start()

	// generate link to join the pool
	link := "/app?join=" + poolId
	utils.Cp("yellow", "Pool link:", utils.Cs("whiteU", "http://localhost:1323"+link))

	// send the link for the same
	return c.Render(http.StatusOK, "createPool", map[string]any{
		"Link": link,
	})
}

// GET /ws?poolId=234bkj&clientId=123123&clientName=joy
// handle socket connections for the pools
func HandlerWsConnection(c echo.Context) error {
	// get the poolId from query params
	poolId := c.QueryParam("poolId")

	// register connection
	return socket.ServeWs(Hub[poolId], c.Response().Writer, c.Request())
}
