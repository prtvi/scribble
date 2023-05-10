package routes

import (
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

type Pool struct {
	ID                   string
	Capacity             int
	Register             chan *Client
	Unregister           chan *Client
	Clients              []*Client
	Broadcast            chan model.SocketMessage
	ColorAssignmentIndex int
	CreatedTime          time.Time
	GameStartTime        time.Time
	HasGameStarted       bool

	// not initialised when start
	CurrWord          string
	CurrWordExpiresAt time.Time
	CurrSketcher      *Client
}

func NewPool(uuid string, capacity int) *Pool {
	// returns a new Pool
	now := time.Now()
	later := now.Add(time.Second * GameStartDurationInSeconds)

	return &Pool{
		ID:                   uuid,
		Capacity:             capacity,
		Register:             make(chan *Client),
		Unregister:           make(chan *Client),
		Clients:              make([]*Client, 0),
		Broadcast:            make(chan model.SocketMessage),
		ColorAssignmentIndex: 0,
		CreatedTime:          now,
		GameStartTime:        later,
		HasGameStarted:       false,
	}
}

func (pool *Pool) Start() {
	// start listening to pool connections and messages
	for {
		select {
		case client := <-pool.Register:
			// on client register, append the client to Pool.Client slice
			pool.Clients = append(pool.Clients, client)
			pool.ColorAssignmentIndex += 1

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client connected:"), client.Name)

			// all clients (c from loop) to one (registered client): all-1
			for _, c := range pool.Clients {
				c.Conn.WriteJSON(model.SocketMessage{
					Type:       1,
					Content:    fmt.Sprintf("CONNECTED_%s", client.Name),
					ClientId:   client.ID,
					ClientName: client.Name,
				})
			}
			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice
			pool.Clients = removeClientFromList(pool.Clients, client)
			// pool.ColorAssignmentIndex -= 1 // TODO

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client disconnected:"), client.Name)

			// all clients (c from loop) to one (disconnected client): all-1
			for _, c := range pool.Clients {
				c.Conn.WriteJSON(model.SocketMessage{
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

			switch message.Type {
			case 3: // update client score on correct guess of word
				updateScore(pool, message)

			case 5: // client info list
				message = responseMessageType5(pool)

			case 6: // start game
				message = responseMessageType6(pool)

			default:
				break
			}

			for _, c := range pool.Clients {
				c.Conn.WriteJSON(message)
			}
		}
	}
}
