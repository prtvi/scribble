package socket

import (
	"fmt"
	"net/http"
	utils "scribble/utils"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// / middleware
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dt := utils.FormatTimeLong(time.Now())
		reqMethod := c.Request().Method

		utils.Cp("green", fmt.Sprintf("%s: %s  %s", reqMethod, fmt.Sprintf("%s", c.Request().URL), dt))

		return next(c)
	}
}

// GET /
func Welcome(c echo.Context) error {
	return c.Render(http.StatusOK, "welcome", map[string]any{
		"debug": debug,
	})
}

// GET /ws?poolId=234bkj&clientId=123123&clientName=joy&clientColor=2def45
func HandlerWsConnection(c echo.Context) error {
	// handle socket connections for the pools

	// get the query params
	poolId := c.QueryParam("poolId")
	clientId := c.QueryParam("clientId")
	clientName := c.QueryParam("clientName")
	clientColor := c.QueryParam("clientColor")

	// register the socket connection from client
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		fmt.Fprintf(c.Response().Writer, "%+v\n", err)
	}

	pool := hub[poolId]

	// create a new client to append to Pool.Clients map
	client := &Client{
		ID:            clientId,
		Name:          clientName,
		Color:         clientColor,
		DoneSketching: false,
		HasGuessed:    false,
		Score:         0,
		Conn:          conn,
		Pool:          pool,
	}

	// register and notify other clients
	pool.Register <- client
	client.read()

	return nil
}

// -----------------------------------------------------------------------------

// GET /create-room
func CreateRoom(c echo.Context) error {
	// render a form to create a new pool
	return c.Render(http.StatusOK, "createRoom", map[string]any{
		"RoomCreated": false,
		"Link":        "",

		"debug": debug,
	})
}

// POST /create-room
func CreateRoomLink(c echo.Context) error {
	// on post request to this route, create a new pool, start listening to connections on that pool, render the link to join this pool

	// get the pool capacity from form input
	capacity, _ := strconv.Atoi(c.FormValue("capacity"))
	utils.Cp("yellow", "Pool capacity:", utils.Cs("white", c.FormValue("capacity")))

	// create a new pool with an uuid
	poolId := utils.GenerateUUID()
	pool := newPool(poolId, capacity)

	// append to global Hub map, and start listening to pool connections
	hub[poolId] = pool
	go pool.start()

	printHubStatus()

	// generate link to join the pool
	link := "/app?join=" + poolId
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", link) // TODO

	// send the link for the same
	return c.Render(http.StatusOK, "createRoom", map[string]any{
		"RoomCreated": true,
		"Link":        link,

		// show on submit, room size in input field
		"Capacity": pool.Capacity,

		"debug": debug,
	})
}

// -----------------------------------------------------------------------------

// GET /app
func App(c echo.Context) error {
	// if /app?join=poolId, then render the playing areax
	// if /app          , then render message

	poolId := c.QueryParam("join")

	// if poolId is empty then do not render any forms, just display message
	if poolId == "" {
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Hi there, are you lost?!",

			"debug": debug,
		})
	}

	// check if pool exists, if is does not exist then render no form
	pool, ok := hub[poolId]
	if !ok {
		// if not then do not render both forms and display message
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Pool expired or non-existent!",

			"debug": debug,
		})
	}

	// if game has already started then do not render both forms and display message
	if pool.HasGameStarted {
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Sorry! The game has already started! 🥹",

			"debug": debug,
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

			"debug": debug,
		})
	}

	// else if every check, checks out then render "RegisterToPool" form
	return c.Render(http.StatusOK, "app", map[string]any{
		"RegisterToPool": true,
		"ConnectSocket":  false,
		"CurrentSize":    len(pool.Clients),
		"Message":        "",

		// hidden in form, added as hidden in "RegisterToPool" form to submit later when POST request is made to join the pool
		"PoolId": poolId,

		"debug": debug,
	})
}

// POST /app
func RegisterToPool(c echo.Context) error {
	// on post request made to this route to capture clientName from "RegisterToPool" post form

	poolId := c.FormValue("poolId")
	clientName := c.FormValue("clientName")

	// extra check to prevent user from joining any random pool which does not exist
	pool, ok := hub[poolId]
	if !ok {
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Pool expired or non-existent!",

			"debug": debug,
		})
	}

	// if client reloads after game has already started
	if pool.HasGameStarted {
		return c.Render(http.StatusOK, "app", map[string]any{
			"RegisterToPool": false,
			"ConnectSocket":  false,
			"Message":        "Sorry! The game has already started! 🥹",

			"debug": debug,
		})
	}

	// generate client id and color
	clientId := utils.GenerateUUID()[0:8]
	clientColor := pool.getColorForClient()

	// render ConnectSocket form to establish socket connection
	// socket connection will start only if "ConnectSocket" form is rendered
	return c.Render(http.StatusOK, "app", map[string]any{
		"RegisterToPool": false,
		"ConnectSocket":  true,
		"JoiningLink":    pool.JoiningLink,
		"Message":        "",

		// variables in DOM
		"GameStartDurationInSeconds": utils.GetSecondsLeftFrom(pool.GameStartTime),
		"TimeForEachWordInSeconds":   TimeForEachWordInSeconds,

		// for rendering title on browser
		"ClientNameExists": true,

		// init as js vars
		"PoolId":        poolId,
		"ClientId":      clientId,
		"ClientName":    clientName,
		"ClientColor":   clientColor,
		"GameStartTime": utils.FormatTimeLong(pool.GameStartTime),

		"debug": debug,
	})
}
