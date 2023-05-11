package routes

import (
	"encoding/json"
	model "scribble/model"
	utils "scribble/utils"
	"strings"
	"time"
)

const GameStartDurationInSeconds = 120
const TimeForEachWordInSeconds = 30
const ScoreForCorrectGuess = 15

func removeClientFromList(list []*Client, client *Client) []*Client {
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
	var client *Client = pool.Clients[0]

	for _, c := range pool.Clients {
		if !c.HasSketched {
			client = c
			break
		}
	}

	client.HasSketched = true
	return client
}

func updateScore(pool *Pool, message model.SocketMessage) {
	// update score for the client that guesses the word right

	var guesserClient *Client
	for _, c := range pool.Clients {
		// increment score only if the guesser is not the sketcher
		if c.ID == message.ClientId &&
			pool.CurrSketcher.ID != message.ClientId {
			guesserClient = c
			break
		}
	}

	// if the sketcher is the guesser, then the guesserClient will be nil, hence check if guesserClient is nil
	// check if the word matches with the current word
	if guesserClient != nil && strings.ToLower(message.Content) == strings.ToLower(pool.CurrWord) {
		guesserClient.Score += ScoreForCorrectGuess * int(utils.GetDiffBetweenTimesInSeconds(time.Now(), pool.CurrWordExpiresAt))
	}
}

func responseMessageType5(pool *Pool) model.SocketMessage {
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
		Type:    5,
		Content: string(byteInfo),
	}
}

func responseMessageType6(pool *Pool) model.SocketMessage {
	// returns if game has started or not embedded in model.SocketMessage

	if pool.HasGameStarted {
		return model.SocketMessage{
			Type:              6,
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
		Type:              6,
		Content:           "true",
		CurrSketcherId:    pool.CurrSketcher.ID,
		CurrWord:          pool.CurrWord,
		CurrWordExpiresAt: pool.CurrWordExpiresAt,
	}
}

func beginClientSketchingFlow(pool *Pool) {
	pool.CurrWord = utils.GetRandomWord()
	pool.CurrWordExpiresAt = time.Now().Add(time.Second * TimeForEachWordInSeconds)
	pool.CurrSketcher = pickClient(pool)
}
