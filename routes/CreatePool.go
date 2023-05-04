package routes

import (
	"fmt"
	"net/http"
	socket "scribble/socket"
	utils "scribble/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

var dataForCreatePoolRoute = map[string]any{
	"Link": "",
}

// GET /create-pool
func CreatePool(c echo.Context) error {
	// render a form to create a new pool
	return c.Render(http.StatusOK, "createPool", dataForCreatePoolRoute)
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
	dataForCreatePoolRoute["Link"] = link
	return c.Render(http.StatusOK, "createPool", dataForCreatePoolRoute)
}
