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

// 10
func (pool *Pool) broadcastConfigs() {
	cfg := model.SharedConfig{
		MessageTypeMap:  messageTypeMap,
		TimeForEachWord: utils.DurationToSeconds(TimeForEachWordInSeconds),
	}

	byteInfo, _ := json.Marshal(cfg)
	pool.broadcast(model.SocketMessage{
		Type:    10,
		TypeStr: messageTypeMap[10],
		Content: string(byteInfo),
	})
}

// 1
func (pool *Pool) broadcastClientRegister(id, name string) {
	pool.broadcast(model.SocketMessage{
		Type:       1,
		TypeStr:    messageTypeMap[1],
		ClientId:   id,
		ClientName: name,
	})
}

// 2
func (pool *Pool) broadcastClientUnregister(id, name string) {
	pool.broadcast(model.SocketMessage{
		Type:       2,
		TypeStr:    messageTypeMap[2],
		ClientId:   id,
		ClientName: name,
	})
}

// 6
func (pool *Pool) broadcastClientInfoList() {
	pool.broadcast(pool.getClientInfoList())
}

// 71
func (pool *Pool) broadcastRoundNumber() {
	pool.broadcast(model.SocketMessage{
		Type:      71,
		TypeStr:   messageTypeMap[71],
		CurrRound: pool.CurrRound,
	})
}

// 51
func (pool *Pool) broadcastClearCanvasEvent() {
	pool.broadcast(model.SocketMessage{
		Type:    51,
		TypeStr: messageTypeMap[51],
	})
}

// 33
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

// 8
func (pool *Pool) broadcastCurrentWordDetails() {
	pool.broadcast(model.SocketMessage{
		Type:              8,
		TypeStr:           messageTypeMap[8],
		CurrSketcherId:    pool.CurrSketcher.ID,
		CurrWord:          pool.CurrWord,
		CurrWordExpiresAt: utils.FormatTimeLong(pool.CurrWordExpiresAt),
	})
}

// 81
func (pool *Pool) broadcastTurnOver() {
	pool.broadcast(model.SocketMessage{
		Type:           81,
		TypeStr:        messageTypeMap[81],
		CurrSketcherId: pool.CurrSketcher.ID,
	})
}

// 32
func (pool *Pool) broadcastWordReveal() {
	pool.broadcast(model.SocketMessage{
		Type:    32,
		TypeStr: messageTypeMap[32],
		Content: pool.CurrWord,
	})
}
