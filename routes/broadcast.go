package routes

import (
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

func (pool *Pool) BroadcastMsg(message model.SocketMessage) {
	PrintSocketMessage(message)

	// broadcasts the given message to all clients in the pool
	for _, c := range pool.Clients {
		c.mu.Lock()
		err := c.Conn.WriteJSON(message)
		c.mu.Unlock()

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (pool *Pool) BeginBroadcastClientInfoMessage() {
	// to be run as a go routine
	// starts an infinite loop to broadcast client info after every regular interval
	for {
		time.Sleep(RenderClientsEvery)
		pool.BroadcastMsg(pool.getClientInfoList())

		// stop broadcasting when game ends
		if pool.HasGameEnded || len(pool.Clients) == 0 {
			utils.Cp("yellowBg", "Stopped broadcasting client info")
			break
		}
	}
}
