package routes

import (
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

type Pool struct {
	ID                                                          string
	Capacity, ColorAssignmentIndex                              int
	Register, Unregister                                        chan *Client
	Clients                                                     []*Client
	Broadcast                                                   chan model.SocketMessage
	CreatedTime, GameStartTime, CurrWordExpiresAt               time.Time
	HasGameStarted, HasGameEnded, HasClientInfoBroadcastStarted bool
	CurrSketcher                                                *Client
	CurrWord                                                    string
}

func NewPool(uuid string, capacity int) *Pool {
	// returns a new Pool
	now := time.Now()
	later := now.Add(time.Second * GameStartDurationInSeconds)

	return &Pool{
		ID:                            uuid,
		Capacity:                      capacity,
		Register:                      make(chan *Client),
		Unregister:                    make(chan *Client),
		Clients:                       make([]*Client, 0),
		Broadcast:                     make(chan model.SocketMessage),
		ColorAssignmentIndex:          0,
		CreatedTime:                   now,
		GameStartTime:                 later,
		CurrWordExpiresAt:             time.Time{},
		HasGameStarted:                false,
		HasGameEnded:                  false,
		HasClientInfoBroadcastStarted: false,
		CurrSketcher:                  nil,
		CurrWord:                      "",
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
			pool.BroadcastMsg(model.SocketMessage{
				Type:       1,
				Content:    fmt.Sprintf("CONNECTED_%s", client.Name),
				ClientId:   client.ID,
				ClientName: client.Name,
			})

			// start broadcasting client info list
			if len(pool.Clients) == 1 && !pool.HasClientInfoBroadcastStarted {
				pool.HasClientInfoBroadcastStarted = true
				utils.Cp("yellowBg", "broadcasting client info start!")

				go pool.BroadcastClientInfoMessage()
			}

			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice
			pool.Clients = removeClientFromList(pool.Clients, client)
			// pool.ColorAssignmentIndex -= 1 // TODO

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client disconnected:"), client.Name)

			// all clients (c from loop) to one (disconnected client): all-1
			pool.BroadcastMsg(model.SocketMessage{
				Type:       2,
				Content:    fmt.Sprintf("DISCONNECTED_%s", client.Name),
				ClientId:   client.ID,
				ClientName: client.Name,
			})

			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message to all clients
			// any of the game logic there is will be applied when clients do something, which will happen after the message is received from any of the clients

			utils.Cp("blue", "sm recv, type:", utils.Cs("yellow", fmt.Sprintf("%d:", message.Type)), utils.Cs("reset", messageTypeMap[message.Type], utils.Cs("blue", "from:"), message.ClientName), utils.Cs("blue", "broadcasting ..."))

			switch message.Type {
			case 3:
				updateScore(pool, message)
				pool.BroadcastMsg(message)

			case 7:
				pool.StartGame()

			case 8:
				message = nextClientForSketching(pool, message.Type)

			case 9:
				pool.EndGame()

			default:
				break
			}

			// dont broadcast everything to everyone
			// for _, c := range pool.Clients {
			// 	c.Conn.WriteJSON(message)
			// }
		}
	}
}
