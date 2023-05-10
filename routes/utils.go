package routes

import (
	"encoding/json"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

const GameStartDurationInSeconds = 120
const TimeForEachWordInSeconds = 90

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

func responseMessageType5(poolId string) model.SocketMessage {
	// returns client info list embedded in model.SocketMessage

	clientInfoList := make([]model.ClientInfo, 0)
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
		clientInfoList = append(clientInfoList, model.ClientInfo{
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

func responseMessageType6(poolId string) model.SocketMessage {
	// returns if game has started or not embedded in model.SocketMessage

	pool, ok := HUB[poolId]

	// if pool does not exist then send false
	if !ok {
		return model.SocketMessage{
			Type:    6,
			Content: "false",
		}
	}

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

	pool.CurrWord = utils.GetRandomWord()
	pool.CurrWordExpiresAt = time.Now().Add(time.Second * TimeForEachWordInSeconds)
	pool.CurrSketcher = pickClient(pool)

	utils.Cp("yellow", "Game started!")

	return model.SocketMessage{
		Type:              6,
		Content:           "true",
		CurrSketcherId:    pool.CurrSketcher.ID,
		CurrWord:          pool.CurrWord,
		CurrWordExpiresAt: pool.CurrWordExpiresAt,
	}
}
