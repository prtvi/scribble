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

var Pools = map[string]*socket.Pool{}

// middleware
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		dt := time.Now().String()[0:19]
		fmt.Println(fmt.Sprintf("\033[32m%s: %s at %s \033[0m\n", c.Request().Method, c.Request().URL, dt))
		return next(c)
	}
}

// GET /app
// if /app?join=sfds, then render the playing area
// if /app          , then render welcome page to display "begin game"
func App(c echo.Context) error {
	poolId := c.QueryParam("join")

	ForPlaying := false
	if poolId != "" {
		ForPlaying = true
	}

	return c.Render(http.StatusOK, "app", map[string]any{
		"ForPlaying": ForPlaying,
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
	nMembers, _ := strconv.Atoi(c.FormValue("nMembers"))

	// create a new pool with an id
	poolId := utils.GenerateUUID()

	pool := socket.NewPool(poolId, nMembers)
	Pools[poolId] = pool
	go pool.Start()

	// generate link to join the pool
	link := "/app?join=" + poolId
	fmt.Println("Pool link:", link)

	// send the link for the same
	return c.Render(http.StatusOK, "createPool", map[string]any{
		"Link": link,
	})
}

type ResponseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func CheckPool(c echo.Context) error {
	poolId := c.QueryParam("poolId")

	pool, ok := Pools[poolId]
	if !ok {
		return c.JSON(http.StatusOK, ResponseMessage{
			Code:    http.StatusBadRequest,
			Message: "Pool expired or non-existent",
		})
	}

	poolCap := pool.Capacity
	poolCurrSizePlus1 := len(pool.Clients) + 1

	if poolCurrSizePlus1 > poolCap {
		return c.JSON(http.StatusOK, ResponseMessage{
			Code:    http.StatusExpectationFailed,
			Message: "Too many client connection requests",
		})
	}

	return c.JSON(http.StatusOK, ResponseMessage{
		Code:    http.StatusOK,
		Message: "Ready to make socket connection",
	})
}

// GET /ws?poolId=234bkj&clientId=123123
// handle socket connections for the pools
func HandlerWsConnection(c echo.Context) error {
	poolId := c.QueryParam("poolId")
	clientId := c.QueryParam("clientId")
	fmt.Println(poolId, clientId)

	return socket.ServeWs(Pools[poolId], c.Response().Writer, c.Request())
}
