package routes

import (
	"fmt"
	"net/http"

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
		Conn:        conn,
		Pool:        pool,
	}

	// register and notify other clients
	pool.Register <- client
	client.Read()

	return nil
}
