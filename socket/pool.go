package socket

import "fmt"

// connected    type 1
// disconnected type 2
// json data    type 3

type Message struct {
	Type       int    `json:"type"`
	Content    string `json:"content"`
	ClientName string `json:"clientName,omitempty"`
	ClientId   string `json:"clientId,omitempty"`
}

type Pool struct {
	ID         string
	Capacity   int
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

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
			fmt.Println("Size of connection pool: ", len(pool.Clients), "client connected", client.Name)

			// all clients (c from loop) to one (registered client): all-1
			for c := range pool.Clients {
				c.Conn.WriteJSON(Message{
					Type:       1,
					Content:    fmt.Sprintf("CONNECTED_%s", client.Name),
					ClientId:   client.ID,
					ClientName: client.Name,
				})
			}
			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client map
			delete(pool.Clients, client)
			fmt.Println("Size of connection pool: ", len(pool.Clients), "client disconnected", client.Name)

			// all clients (c from loop) to one (disconnected client): all-1
			for c := range pool.Clients {
				c.Conn.WriteJSON(Message{
					Type:       2,
					Content:    fmt.Sprintf("DISCONNECTED_%s", client.Name),
					ClientId:   client.ID,
					ClientName: client.Name,
				})
			}
			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message to all clients
			fmt.Println("Sending message to all clients in pool:", message)

			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
