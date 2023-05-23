package routes

import (
	"encoding/json"
	model "scribble/model"
)

const (
	GameStartDurationInSeconds = 30
	TimeForEachWordInSeconds   = 60
	ScoreForCorrectGuess       = 15
	RenderClientsEvery         = 10
)

var messageTypeMap = map[int]string{
	1: "connected client",    // server b=> clients
	2: "disconnected client", // server b=> clients
	3: "text message",        // client b=> clients
	4: "canvas data",         // client b=> clients
	5: "clear canvas",        // client b=> clients
	6: "client info",         // server b=> clients --at regular intervals
	7: "start game",          // client  => server --to start the game
	8: "word assigned",       // server  => client

	9:  "req next word", //
	10: "all clients done playing",
}

func removeClientFromList(list []*Client, client *Client) []*Client {
	// removes the given client from the given slice and returns the new slice
	var idxToRemove int
	for i, c := range list {
		if c == client {
			idxToRemove = i
			break
		}
	}

	list[idxToRemove] = list[len(list)-1]
	return list[:len(list)-1]
}

func pickClient(pool *Pool) *Client {
	// picks that client that hasn't drawn yet

	var client *Client = nil

	for _, c := range pool.Clients {
		if !c.HasSketched {
			client = c
			break
		}
	}

	if client != nil {
		client.HasSketched = true
		return client
	}

	return nil
}

func getClientInfoList(pool *Pool, messageType int) model.SocketMessage {
	// returns client info list embedded in model.SocketMessage

	clientInfoList := make([]model.ClientInfo, 0)

	// append client info into an array
	for _, client := range pool.Clients {
		clientInfoList = append(clientInfoList, model.ClientInfo{
			ID:    client.ID,
			Name:  client.Name,
			Color: client.Color,
			Score: client.Score,
		})
	}

	// marshall array in byte and send as string
	byteInfo, _ := json.Marshal(clientInfoList)
	return model.SocketMessage{
		Type:    messageType,
		Content: string(byteInfo),
	}
}
