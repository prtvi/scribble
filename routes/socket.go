package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	model "scribble/model"
	utils "scribble/utils"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// register the socket connection from client
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return ws, err
	}
	return ws, nil
}

// serves the websocket and registers the client to the pool
func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) error {
	clientId := r.URL.Query().Get("clientId")
	clientName := r.URL.Query().Get("clientName")
	clientColor := r.URL.Query().Get("clientColor")

	// register to socket connection
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	// create a new client to append to Pool.Clients map
	client := &Client{
		ID:          clientId,
		Name:        clientName,
		Color:       clientColor,
		HasSketched: false,
		HasGuessed:  false,
		Score:       0,
		Conn:        conn,
		Pool:        pool,
	}

	// register and notify other clients
	pool.Register <- client
	client.Read()

	return nil
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

func Maintainer() {
	// clears the pools in which the game has ended every 10 mins

	for {
		time.Sleep(time.Minute * 10) // TODO - to be tested

		for key, pool := range HUB {
			if pool != nil && pool.HasGameEnded {
				utils.Cp("yellowBg", "Removing pool from HUB", key)
				delete(HUB, key)
			}
		}
	}
}

func DebugMode() {
	poolId := "debug"
	pool := NewPool(poolId, 4)

	HUB[poolId] = pool
	go pool.Start()

	link := "/app?join=" + poolId
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", link)
}
