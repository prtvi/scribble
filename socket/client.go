package socket

import (
	"encoding/json"
	"fmt"
	"net/http"
	utils "scribble/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID, Name string
	Conn     *websocket.Conn
	Pool     *Pool
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
		var clientMsg Message
		err = json.Unmarshal(msgByte, &clientMsg)
		utils.Cp("blue", "Message received:", utils.Cs("white", fmt.Sprintf("%+v", clientMsg)))

		// broadcast the message to all clients in the pool
		c.Pool.Broadcast <- clientMsg
	}
}

// serves the websocket and registers the client to the pool
func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) error {
	clientId := r.URL.Query().Get("clientId")
	clientName := r.URL.Query().Get("clientName")

	// register to socket connection
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	// create a new client to append to Pool.Clients map
	client := &Client{
		ID:   clientId,
		Name: clientName,
		Conn: conn,
		Pool: pool,
	}

	// register and notify other clients
	pool.Register <- client
	client.Read()

	return nil
}
