package socket

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID, Name, Color string
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
		var clientMsg SocketMessage
		err = json.Unmarshal(msgByte, &clientMsg)

		// broadcast the message to all clients in the pool
		c.Pool.Broadcast <- clientMsg
	}
}
