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

// GET /
func Welcome(c echo.Context) error {
	aboutText := []string{"scribble is a free online multiplayer drawing and guessing pictionary game.", "A normal game consists of a few rounds, where every round a player has to draw their chosen word and others have to guess it to gain points!", "The person with the most points at the end of the game, will then be crowned as the winner!"}

	howToSlides := []string{"When it's your turn, choose a word you want to draw!",
		"Try to draw your choosen word! No spelling!",
		"Let other players try to guess your drawn word!",
		"When it's not your turn, try to guess what other players are drawing!",
		"Score the most points and be crowned the winner at the end!"}

	return c.Render(http.StatusOK, "welcome", map[string]any{
		"StyleSheets": []string{"global"},
		"AboutText":   aboutText,
		"HowToSlides": howToSlides,
		"debug":       debug,
	})
}

// -----------------------------------------------------------------------------

// GET /create-room
func CreateRoom(c echo.Context) error {
	// render a form to create a new pool
	return c.Render(http.StatusOK, "createRoom", map[string]any{
		"StyleSheets": []string{"global", "createRoom"},
		"FormParams":  FormParams,
		"RoomCreated": false,
		"debug":       debug,
	})
}

// POST /create-room
func CreateRoomLink(c echo.Context) error {
	// on post request to this route, create a new pool, start listening to connections on that pool, render the link to join this pool

	// get the pool capacity from form input
	capacity, _ := strconv.Atoi(c.FormValue("players"))
	utils.Cp("yellow", "Pool capacity:", utils.Cs("white", c.FormValue("players")))

	// create a new pool with an uuid
	poolId := utils.GenerateUUID()
	pool := newPool(poolId, capacity)

	// append to global Hub map, and start listening to pool connections
	hub[poolId] = pool
	go pool.start()

	// generate link to join the pool
	link := "/app?join=" + poolId
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", link) // TODO

	// send the link for the same
	return c.Render(http.StatusOK, "createRoom", map[string]any{
		"StyleSheets": []string{"global", "createRoom"},
		"RoomCreated": true,
		"Link":        link,
		"Capacity":    pool.Capacity, // show on submit, room size in input field

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
	if poolId == "" || len(poolId) == 0 {
		return c.Render(http.StatusOK, "error", map[string]any{
			"StyleSheets": []string{"global"},
			"Message":     "Hi there, are you lost?! The link seems to be broken. Make sure you copied the link properly! 😃",

			"debug": debug,
		})
	}

	// check if pool exists, if is does not exist then render no form
	pool, ok := hub[poolId]
	if !ok {
		// if not then do not render both forms and display message
		return c.Render(http.StatusOK, "error", map[string]any{
			"StyleSheets": []string{"global"},
			"Message":     "Pool expired or non-existent! Make sure you have the correct link! 😃",

			"debug": debug,
		})
	}

	// if game has already started then do not render both forms and display message
	if pool.HasGameStarted {
		return c.Render(http.StatusOK, "error", map[string]any{
			"StyleSheets": []string{"global"},
			"Message":     "Oops! The game has already started! 🥹",

			"debug": debug,
		})
	}

	// if pool exists, get its capacity and curr size
	poolCap := pool.Capacity
	poolCurrSizePlus1 := len(pool.Clients) + 1

	if poolCurrSizePlus1 > poolCap {
		// if poolCurrSizePlus1 is greater than capacity then do not render both forms and display message
		return c.Render(http.StatusOK, "error", map[string]any{
			"StyleSheets": []string{"global"},
			"Message":     "Your party is full! Maximum room capacity reached! 😃",

			"debug": debug,
		})
	}

	// else if every check, checks out then render "RegisterToPool" form
	return c.Render(http.StatusOK, "join", map[string]any{
		"StyleSheets": []string{"global", "createRoom"},

		// hidden in form, added as hidden in "RegisterToPool" form to submit later when POST request is made to join the pool
		"PoolId": poolId,

		"debug":       debug,
		"currentSize": len(pool.Clients), // used in getting client num, used in debug mode
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
		return c.Render(http.StatusOK, "error", map[string]any{
			"StyleSheets": []string{"global"},
			"Message":     "Pool expired or non-existent! Make sure you have the correct link! 😃",

			"debug": debug,
		})
	}

	// if client reloads after game has already started
	if pool.HasGameStarted {
		return c.Render(http.StatusOK, "error", map[string]any{
			"StyleSheets": []string{"global"},
			"Message":     "Oops! The game has already started! 🥹",

			"debug": debug,
		})
	}

	// generate client id and color
	clientId := utils.GenerateUUID()[0:8]
	clientColor := pool.getColorForClient()

	// render ConnectSocket form to establish socket connection
	// socket connection will start only if "ConnectSocket" form is rendered
	return c.Render(http.StatusOK, "app", map[string]any{
		"StyleSheets": []string{"global", "app"},
		"JoiningLink": pool.JoiningLink,

		// variables in DOM
		"GameStartDurationInSeconds": utils.GetSecondsLeftFrom(pool.GameStartTime),

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
