package routes

import (
	"encoding/json"
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
}

// returns a new Pool
func NewPool(uuid string, capacity int) *Pool {
	now := time.Now()
	// later := now.Add(time.Minute * 2)

	later := now.Add(time.Second * 15) // later diff to be added on config

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

func removeClientFromList(list []*Client, client *Client) []*Client {
	var idxToRemove int
	for i, c := range list {
		if c == client {
			idxToRemove = i
			break
		}
	}

	list[idxToRemove] = list[len(list)-1]
	list = list[:len(list)-1]
	return list
}

// start listening to pool connections and messages
func (pool *Pool) Start() {
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
			// pool.ColorAssignmentIndex -= 1

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
			case 5: // client info list
				message = ResponseMessageType_5(message.PoolId)

			case 6: // start game ack
				message = ResponseMessageType_6(message.PoolId)

			default:
				break
			}

			for _, c := range pool.Clients {
				c.Conn.WriteJSON(message)
			}
		}
	}
}

func ResponseMessageType_5(poolId string) model.SocketMessage {
	// returns client info list embedded in model.SocketMessage

	type clientInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	clientInfoList := make([]clientInfo, 0)
	pool, ok := HUB[poolId]

	// if pool does not exist then send empty list
	if !ok {
		return model.SocketMessage{
			Type:    5,
			Content: "[]",
		}
	}

	// append client info into an array
	for _, client := range pool.Clients {
		clientInfoList = append(clientInfoList, clientInfo{
			ID:    client.ID,
			Name:  client.Name,
			Color: client.Color,
		})
	}

	// marshall array in byte and send as string
	byteInfo, _ := json.Marshal(clientInfoList)
	return model.SocketMessage{
		Type:    5,
		Content: string(byteInfo),
	}
}

func ResponseMessageType_6(poolId string) model.SocketMessage {
	// returns if game has started or not embedded in model.SocketMessage

	pool, ok := HUB[poolId]

	// if pool does not exist then send false
	if !ok {
		return model.SocketMessage{
			Type:    6,
			Content: "false",
		}
	}

	// flag game started variable for the pool as true
	pool.HasGameStarted = true
	return model.SocketMessage{
		Type:    6,
		Content: "true",
	}
}
