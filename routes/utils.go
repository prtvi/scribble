package routes

import (
	"encoding/json"
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

const (
	GameStartDurationInSeconds = 30
	TimeForEachWordInSeconds   = 20
	ScoreForCorrectGuess       = 15
	RenderClientsEvery         = 10
)

var messageTypeMap = map[int]string{
	1: "connected client",
	2: "disconnected client",
	3: "text message",
	4: "canvas data",
	5: "clear canvas",
	6: "client info",
	7: "start game req",
	8: "req next word",
	9: "all clients done playing",
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

//

func startGameAck(pool *Pool, messageType int) model.SocketMessage {
	// returns if game has started or not embedded in model.SocketMessage

	if pool.HasGameStarted {
		return model.SocketMessage{
			Type:              messageType,
			Content:           "true",
			CurrSketcherId:    pool.CurrSketcher.ID,
			CurrWord:          pool.CurrWord,
			CurrWordExpiresAt: pool.CurrWordExpiresAt,
		}
	}

	// flag game started variable for the pool as true
	pool.HasGameStarted = true
	utils.Cp("yellow", "Game started!")

	beginClientSketchingFlow(pool)

	return model.SocketMessage{
		Type:              messageType,
		Content:           "true",
		CurrSketcherId:    pool.CurrSketcher.ID,
		CurrWord:          pool.CurrWord,
		CurrWordExpiresAt: pool.CurrWordExpiresAt,
	}
}

func nextClientForSketching(pool *Pool, messageType int) model.SocketMessage {
	// if this request was previously made which means the current word is set, which means the expiry of the word is in future, then just return the curr stat

	if pool.CurrWordExpiresAt.Sub(time.Now()) > 0 {
		return model.SocketMessage{
			Type:              messageType,
			Content:           "true",
			CurrSketcherId:    pool.CurrSketcher.ID,
			CurrWord:          pool.CurrWord,
			CurrWordExpiresAt: pool.CurrWordExpiresAt,
		}
	}

	// else begin the client sketching flow
	isClient := beginClientSketchingFlow(pool)
	if !isClient {
		// if no client left to pick then end the game by sending the scores, type 9
		pool.HasGameEnded = true
		fmt.Println("no client found")
		return getClientInfoList(pool, 9)
	}

	return model.SocketMessage{
		Type:              messageType,
		Content:           "true",
		CurrSketcherId:    pool.CurrSketcher.ID,
		CurrWord:          pool.CurrWord,
		CurrWordExpiresAt: pool.CurrWordExpiresAt,
	}
}

func beginClientSketchingFlow(pool *Pool) bool {

	client := pickClient(pool)
	if client == nil {
		return false
	}

	pool.CurrSketcher = client
	pool.CurrWord = utils.GetRandomWord()
	pool.CurrWordExpiresAt = time.Now().Add(time.Second * TimeForEachWordInSeconds)

	fmt.Println("Current word:", pool.CurrWord)

	// reset client.HasGuessed when called upon for next word
	for _, c := range pool.Clients {
		c.HasGuessed = false
	}

	return true
}
