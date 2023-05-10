package routes

import (
	"encoding/json"
	"fmt"
	model "scribble/model"

	"github.com/gorilla/websocket"
)

// Client.Color: string
// color string hash value without the #

type Client struct {
	ID, Name, Color string
	HasSketched     bool
	Score           int
	Conn            *websocket.Conn
	Pool            *Pool
}

// read messages received from client
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgByte, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// parse message received from client
		var clientMsg model.SocketMessage
		err = json.Unmarshal(msgByte, &clientMsg)

		// broadcast the message to all clients in the pool
		c.Pool.Broadcast <- clientMsg
	}
}
