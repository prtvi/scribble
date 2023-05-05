package routes

import (
	"net/http"
	utils "scribble/utils"
	"time"

	"github.com/labstack/echo/v4"
)

var dataForAppRoute = map[string]any{
	"RegisterToPool": false,
	"ConnectSocket":  false,
	"Message":        "",
}

// GET /app
func App(c echo.Context) error {
	// if /app?join=poolId, then render the playing areax
	// if /app          , then render message

	poolId := c.QueryParam("join")

	// if poolId is empty then do not render any forms, just display message
	if poolId == "" {
		dataForAppRoute["Message"] = "Hi there, are you lost?!"
		return c.Render(http.StatusOK, "app", dataForAppRoute)
	}

	// check if pool exists, if is does not exist then render no form
	pool, ok := HUB[poolId]
	if !ok {
		// if not then do not render both forms and display message
		dataForAppRoute["Message"] = "Pool expired or non-existent!"
		return c.Render(http.StatusOK, "app", dataForAppRoute)
	}

	// if pool exists, get its capacity and curr size
	poolCap := pool.Capacity
	poolCurrSizePlus1 := len(pool.Clients) + 1

	if poolCurrSizePlus1 > poolCap {
		// if poolCurrSizePlus1 is greater than capacity then do not render both forms and display message
		dataForAppRoute["Message"] = "Your party is full!"
		return c.Render(http.StatusOK, "app", dataForAppRoute)
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

	// extra check to prevent user from joining any random pool which does not exist
	pool, ok := HUB[poolId]
	if !ok {
		dataForAppRoute["Message"] = "Pool expired or non-existent!"
		return c.Render(http.StatusOK, "app", dataForAppRoute)
	}

	// generate client id and color
	clientId := utils.GenerateUUID()[0:8]
	clientColor := utils.COLORS[pool.ColorAssignmentIndex]

	// check if game has started
	currTime := time.Now()
	diff := pool.GameStartTime.Sub(currTime)
	var hasGameStarted string = ""
	if diff <= 0 {
		hasGameStarted = "started"
	}

	// render ConnectSocket form to establish socket connection
	// socket connection will start only if "ConnectSocket" form is rendered
	return c.Render(http.StatusOK, "app", map[string]any{
		"RegisterToPool": false,
		"ConnectSocket":  true,

		// init as js vars
		"PoolId":              poolId,
		"ClientId":            clientId,
		"ClientName":          clientName,
		"ClientColor":         clientColor,
		"GameStartsInSeconds": diff.Seconds(),
		"HasGameStarted":      hasGameStarted,
	})
}
