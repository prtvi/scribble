package socket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scribble/model"
	utils "scribble/utils"
	"strconv"
	"strings"
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
func Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", map[string]any{
		"RenderTemplateName": "home",
	})
}

// GET /create-room
func CreateRoomForm(c echo.Context) error {
	// render a form to create a new pool
	return c.Render(http.StatusOK, "index", map[string]any{
		"RenderTemplateName": "createRoom",
		"FormParams":         utils.FormParams,
		"RoomCreated":        false,
	})
}

// POST /create-room
func CreateRoom(c echo.Context) error {
	// on post request to this route, create a new pool, start listening to connections on that pool, render the link to join this pool

	players, _ := strconv.Atoi(c.FormValue("players"))
	drawTime, _ := strconv.Atoi(c.FormValue("drawTime"))
	rounds, _ := strconv.Atoi(c.FormValue("rounds"))
	wordCount, _ := strconv.Atoi(c.FormValue("wordCount"))
	hints, _ := strconv.Atoi(c.FormValue("hints"))
	wordMode := c.FormValue("wordMode")
	customWords := utils.SplitIntoWords(c.FormValue("customWords"))
	useCustomWordsOnly := c.FormValue("useCustomWordsOnly") == "on"

	pool, poolId := newPool(players, drawTime, rounds, wordCount, hints, wordMode, customWords, useCustomWordsOnly)

	// append to global Hub map, and start listening to pool connections
	hub[poolId] = pool
	go pool.start()

	// generate link to join the pool
	link := "/app?join=" + poolId
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", link) // TODO
	link += "&isOwner=true"

	// send the link for the same
	return c.Render(http.StatusOK, "index", map[string]any{
		"RenderTemplateName": "createRoom",
		"FormParams":         utils.FormParams,
		"RoomCreated":        true,
		"Link":               link,

		// show on submit value submitted on form
		"Players":            pool.Capacity,
		"DrawTime":           utils.DurationToSeconds(pool.DrawTime),
		"Rounds":             pool.Rounds,
		"WordCount":          pool.WordCount,
		"Hints":              pool.Hints,
		"WordMode":           pool.WordMode,
		"CustomWords":        strings.Join(pool.CustomWords, ","),
		"UseCustomWordsOnly": pool.UseCustomWordsOnly,
	})
}

// GET /app
func JoinPool(c echo.Context) error {
	// if /app?join=poolId, then render the playing areax
	// if /app          , then render message

	poolId := c.QueryParam("join")

	// if poolId is empty then do not render any forms, just display message
	if poolId == "" || len(poolId) == 0 {
		return c.Render(http.StatusOK, "index", map[string]any{
			"RenderTemplateName": "error",
			"Message":            "Hi there, are you lost?! The link seems to be broken. Make sure you copied the link properly! 😃",
		})
	}

	// check if pool exists, if is does not exist then render no form
	pool, ok := hub[poolId]
	if !ok {
		// if not then do not render both forms and display message
		return c.Render(http.StatusOK, "index", map[string]any{
			"RenderTemplateName": "error",
			"Message":            "Pool expired or non-existent! Make sure you have the correct link! 😃",
		})
	}

	// if pool exists, get its capacity and curr size
	poolCap := pool.Capacity
	poolCurrSizePlus1 := len(pool.Clients) + 1

	if poolCurrSizePlus1 > poolCap {
		// if poolCurrSizePlus1 is greater than capacity then do not render both forms and display message
		return c.Render(http.StatusOK, "index", map[string]any{
			"RenderTemplateName": "error",
			"Message":            "Your party is full! Maximum room capacity reached! 😃",
		})
	}

	// else if every check, checks out then render "RegisterToPool" form
	return c.Render(http.StatusOK, "index", map[string]any{
		"RenderTemplateName": "join",

		// hidden in form, added as hidden in "RegisterToPool" form to submit later when POST request is made to join the pool
		"PoolId": poolId,

		"currentSize": len(pool.Clients), // used in getting client num, used in debug mode
	})
}

// POST /app
func EnterPool(c echo.Context) error {
	// on post request made to this route to capture clientName from "RegisterToPool" post form

	poolId := c.FormValue("poolId")
	clientName := c.FormValue("clientName")

	// extra check to prevent user from joining any random pool which does not exist
	pool, ok := hub[poolId]
	if !ok {
		return c.Render(http.StatusOK, "index", map[string]any{
			"RenderTemplateName": "error",
			"Message":            "Pool expired or non-existent! Make sure you have the correct link! 😃",
		})
	}

	return c.Render(http.StatusOK, "app", map[string]any{
		"Rounds": pool.Rounds,
		"Colors": utils.COLORS_FOR_DRAWING,

		// init as js vars
		"PoolId":      poolId,
		"ClientId":    utils.GenerateUUID(),
		"ClientName":  clientName,
		"JoiningLink": pool.JoiningLink,
	})
}

// GET /ws?poolId=234bkj&clientId=123123&clientName=joy&avatarConfig='{}'
func WsConnect(c echo.Context) error {
	// get query params
	poolId := c.QueryParam("poolId")
	clientId := c.QueryParam("clientId")
	clientName := c.QueryParam("clientName")
	avatarConfig := c.QueryParam("avatarConfig")

	var avatarConfigObj model.AvatarConfig
	json.Unmarshal([]byte(avatarConfig), &avatarConfigObj)

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
		AvatarConfig:  avatarConfigObj,
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
