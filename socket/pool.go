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
	ID, Name string
	Conn     *websocket.Conn
	Pool     *Pool
}

// 0 ack (joined/exited) => "CONNECTED" / "DISCONNECTED"
// 1 string
// 2 interface{} / json
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// Pool

// returns a new Pool
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

// start listening to pool connections and messages
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			// on client register, append the client to Pool.Client map
			pool.Clients[client] = true

			fmt.Println("Size of Connection Pool: ", len(pool.Clients), "client connected", client.Name)

			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 0,
					Body: fmt.Sprintf("CONNECTED_%s_%s", client.ID, client.Name),
				})
			}
			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client map
			delete(pool.Clients, client)

			fmt.Println("Size of Connection Pool: ", len(pool.Clients), "client disconnected", client.Name)

			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 0,
					Body: fmt.Sprintf("DISCONNECTED_%s_%s", client.ID, client.Name),
				})
			}
			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message to all clients
			fmt.Println("Sending message to all clients in Pool:", message)

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

// register the socket connection from client
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return ws, err
	}
	return ws, nil
}

// Client

// read messages received from client
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

		// broadcast the message to all clients in the pool
		c.Pool.Broadcast <- message
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

	fmt.Println("New client:", client.ID, client.Name)

	// register and notify other clients
	pool.Register <- client
	client.Read()

	return nil
}
