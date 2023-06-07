package socket

import (
	"encoding/json"
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
)

func (pool *Pool) broadcast(message model.SocketMessage) {
	PrintSocketMessage(message)

	// broadcasts the given message to all clients in the pool
	for _, c := range pool.Clients {
		c.mu.Lock()
		err := c.Conn.WriteJSON(message)
		c.mu.Unlock()

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (pool *Pool) broadcastMessageTypeMap() {
	byteInfo, _ := json.Marshal(messageTypeMap)
	pool.broadcast(model.SocketMessage{
		Type:    10,
		TypeStr: messageTypeMap[10],
		Content: string(byteInfo),
	})
}

func (pool *Pool) broadcastClientRegister(id, name string) {
	pool.broadcast(model.SocketMessage{
		Type:       1,
		TypeStr:    messageTypeMap[1],
		ClientId:   id,
		ClientName: name,
	})
}

func (pool *Pool) broadcastClientUnregister(id, name string) {
	pool.broadcast(model.SocketMessage{
		Type:       2,
		TypeStr:    messageTypeMap[2],
		ClientId:   id,
		ClientName: name,
	})
}

func (pool *Pool) broadcastClientInfoList() {
	pool.broadcast(pool.getClientInfoList())
}

func (pool *Pool) broadcastRoundNumber() {
	pool.broadcast(model.SocketMessage{
		Type:      71,
		TypeStr:   messageTypeMap[71],
		CurrRound: pool.CurrRound,
	})
}

func (pool *Pool) broadcastClearCanvasEvent() {
	pool.broadcast(model.SocketMessage{
		Type:    5,
		TypeStr: messageTypeMap[5],
	})
}

func (pool *Pool) broadcast3WordsList(words []string) {
	byteInfo, _ := json.Marshal(words)
	pool.broadcast(model.SocketMessage{
		Type:             33,
		TypeStr:          messageTypeMap[33],
		Content:          string(byteInfo),
		CurrSketcherId:   pool.CurrSketcher.ID,
		CurrSketcherName: pool.CurrSketcher.Name,
	})
}

func (pool *Pool) broadcastCurrentWordDetails() {
	pool.broadcast(model.SocketMessage{
		Type:              8,
		TypeStr:           messageTypeMap[8],
		CurrSketcherId:    pool.CurrSketcher.ID,
		CurrWord:          pool.CurrWord,
		CurrWordExpiresAt: utils.FormatTimeLong(pool.CurrWordExpiresAt),
	})
}

func (pool *Pool) broadcastTurnOver() {
	pool.broadcast(model.SocketMessage{
		Type:           81,
		TypeStr:        messageTypeMap[81],
		CurrSketcherId: pool.CurrSketcher.ID,
	})
}

func (pool *Pool) broadcastWordReveal() {
	pool.broadcast(model.SocketMessage{
		Type:    32,
		TypeStr: messageTypeMap[32],
		Content: fmt.Sprintf("%s was the correct word!", pool.CurrWord),
	})
}
