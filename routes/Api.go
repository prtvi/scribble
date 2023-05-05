package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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

// GET /api/start-game?poolId=123dedf
func StartGame(c echo.Context) error {
	type startGameAck struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}

	poolId := c.QueryParam("poolId")

	pool, ok := Hub[poolId]
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
