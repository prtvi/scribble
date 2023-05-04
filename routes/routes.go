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

var errorMapForAppRoute = map[string]any{
	"RegisterToPool": false,
	"ConnectSocket":  false,
	"Message":        "",
}

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

// GET /create-pool
func CreatePool(c echo.Context) error {
	// render a form to create a new pool
	return c.Render(http.StatusOK, "createPool", map[string]any{
		"Link": "",
	})
}

// POST /create-pool
func CreatePoolLink(c echo.Context) error {
	// on post request to this route, create a new pool, start listening to connections on that pool, render the link to join this pool

	// get the pool capacity from form input
	capacity, _ := strconv.Atoi(c.FormValue("capacity"))
	utils.Cp("yellow", "Pool capacity:", utils.Cs("white", c.FormValue("capacity")))

	// create a new pool with an uuid
	poolId := utils.GenerateUUID()
	pool := socket.NewPool(poolId, capacity)

	// append to global Hub map, and start listening to pool connections
	Hub[poolId] = pool
	go pool.Start()

	utils.Cp("blue", "Hub size (number of pools):", utils.Cs("white", fmt.Sprintf("%d", len(Hub))))

	// generate link to join the pool
	link := "/app?join=" + poolId
	utils.Cp("yellow", "Pool link:", utils.Cs("whiteU", "http://localhost:1323"+link))

	// send the link for the same
	return c.Render(http.StatusOK, "createPool", map[string]any{
		"Link": link,
	})
}

// GET /app
func App(c echo.Context) error {
	// if /app?join=poolId, then render the playing areax
	// if /app          , then render message

	poolId := c.QueryParam("join")

	// if poolId is empty then do not render any forms, just display message
	if poolId == "" {
		errorMapForAppRoute["Message"] = "Hi there, are you lost?!"
		return c.Render(http.StatusOK, "app", errorMapForAppRoute)
	}

	// check if pool exists, if is does not exist then render no form
	pool, ok := Hub[poolId]
	if !ok {
		// if not then do not render both forms and display message
		errorMapForAppRoute["Message"] = "Pool expired or non-existent!"
		return c.Render(http.StatusOK, "app", errorMapForAppRoute)
	}

	// if pool exists, get its capacity and curr size
	poolCap := pool.Capacity
	poolCurrSizePlus1 := len(pool.Clients) + 1

	if poolCurrSizePlus1 > poolCap {
		// if poolCurrSizePlus1 is greater than capacity then do not render both forms and display message
		errorMapForAppRoute["Message"] = "Your party is full!"
		return c.Render(http.StatusOK, "app", errorMapForAppRoute)
	}

	// else if every check, checks out then render "RegisterToPool" form
	return c.Render(http.StatusOK, "app", map[string]any{
		"RegisterToPool": true,
		"ConnectSocket":  false,

		// hidden in form, added as hidden in "RegisterToPool" form to submit later when POST request is made to join the pool
		"PoolId": poolId,
	})
}

// POST /app
func RegisterToPool(c echo.Context) error {
	// on post request made to this route to capture clientName from "RegisterToPool" post form

	poolId := c.FormValue("poolId")
	clientName := c.FormValue("clientName")

	// generate client id
	clientId := utils.GenerateUUID()[0:8]

	// extra check to prevent user from joining any random pool which does not exist
	if _, ok := Hub[poolId]; !ok {
		errorMapForAppRoute["Message"] = "Pool expired or non-existent!"
		return c.Render(http.StatusOK, "app", errorMapForAppRoute)
	}

	// render ConnectSocket form to establish socket connection
	// socket connection will start only if "ConnectSocket" form is rendered
	return c.Render(http.StatusOK, "app", map[string]any{
		"RegisterToPool": false,
		"ConnectSocket":  true,

		// init as js vars
		"PoolId":     poolId,
		"ClientId":   clientId,
		"ClientName": clientName,
	})
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
