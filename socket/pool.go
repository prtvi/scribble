package socket

import (
	"fmt"
	utils "scribble/utils"
)

// SocketMessage.Content field
// Type 1 - Connected
// Type 2 - Disconnected
// Type 3 - Text message
// Type 4 - Canvas data as string

type SocketMessage struct {
	Type       int    `json:"type"`
	Content    string `json:"content"`
	ClientId   string `json:"clientId,omitempty"`
	ClientName string `json:"clientName,omitempty"`
}

type Pool struct {
	ID                   string
	Capacity             int
	Register             chan *Client
	Unregister           chan *Client
	Clients              map[*Client]bool
	Broadcast            chan SocketMessage
	ColorAssignmentIndex int
}

// returns a new Pool
func NewPool(uuid string, capacity int) *Pool {
	return &Pool{
		ID:                   uuid,
		Capacity:             capacity,
		Register:             make(chan *Client),
		Unregister:           make(chan *Client),
		Clients:              make(map[*Client]bool),
		Broadcast:            make(chan SocketMessage),
		ColorAssignmentIndex: 0,
	}
}

// start listening to pool connections and messages
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			// on client register, append the client to Pool.Client map
			pool.Clients[client] = true
			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client connected:"), client.Name)

			// all clients (c from loop) to one (registered client): all-1
			for c := range pool.Clients {
				c.Conn.WriteJSON(SocketMessage{
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
			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client disconnected:"), client.Name)

			// all clients (c from loop) to one (disconnected client): all-1
			for c := range pool.Clients {
				c.Conn.WriteJSON(SocketMessage{
					Type:       2,
					Content:    fmt.Sprintf("DISCONNECTED_%s", client.Name),
					ClientId:   client.ID,
					ClientName: client.Name,
				})
			}
			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message to all clients
			utils.Cp("blue", "SocketMessage received, type:", utils.Cs("reset", fmt.Sprintf("%d,", message.Type)), utils.Cs("blue", "broadcasting ..."))

			// any of the game logic there is will be applied when clients do something, which will happen after the message is received from any of the clients

			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
