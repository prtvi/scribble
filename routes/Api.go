package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// GET /api/start-game?poolId=123dedf
func StartGame(c echo.Context) error {
	type startGameAck struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}

	poolId := c.QueryParam("poolId")

	pool, ok := HUB[poolId]
	if !ok {
		return c.JSON(http.StatusOK, startGameAck{
			Message: "Game start failure",
			Success: false,
		})
	}

	// for c := range pool.Clients {
	// 	c.Conn.WriteJSON(model.SocketMessage{
	// 		Type:    5,
	// 		Content: "hi there from socket",
	// 	})
	// 	break
	// }

	pool.HasGameStarted = true

	return c.JSON(http.StatusOK, startGameAck{
		Message: "Game start success",
		Success: true,
	})
}
