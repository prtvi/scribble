package socket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Pool struct {
	ID         string
	Capacity   int
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// 0 ack (joined/exited) => "CONNECTED" / "DISCONNECTED"
// 1 string
// 2 interface{} / json
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// Pool
func NewPool(uuid string, capacity int) *Pool {
	return &Pool{
		ID:         uuid,
		Capacity:   capacity,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			// fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			fmt.Println("Size of Connection Pool: ", len(pool.Clients), "client connected", client.ID)

			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 0,
					Body: "CONNECTED__" + client.ID,
				})
			}
			break

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			// fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			fmt.Println("Size of Connection Pool: ", len(pool.Clients), "client disconnected", client.ID)

			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 0,
					Body: "DISCONNECTED__" + client.ID,
				})
			}
			break

		case message := <-pool.Broadcast:
			// fmt.Println("Sending message to all clients in Pool:", message)
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

// Websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return ws, err
	}
	return ws, nil
}

// Client
func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		message := Message{
			Type: messageType,
			Body: string(p),
		}

		fmt.Println("Message Received:", message)
		c.Pool.Broadcast <- message
	}
}

// serves the websocket and registers the client to the pool
func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) error {
	clientId := r.URL.Query().Get("clientId")

	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &Client{
		ID:   clientId,
		Conn: conn,
		Pool: pool,
	}

	fmt.Println("New client:", client.ID)

	pool.Register <- client
	client.Read()

	return nil
}
