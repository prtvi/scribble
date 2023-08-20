package socket

import (
	"encoding/json"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

// sends the given message to every client except the client with excludeId
func (pool *Pool) sendExcludingClientId(excludeId string, message model.SocketMessage) {
	pool.printSocketMsg(message)

	// broadcasts the given message to all clients in the pool
	for _, c := range pool.Clients {
		if c.ID == excludeId {
			continue
		}

		c.send(message)
	}
}

// sends m1: message to client with id -> id1, sends m: message to rest of the clients
func (pool *Pool) sendCorrespondingMessages(id1 string, m1, m model.SocketMessage) {
	pool.printSocketMsg(m1)
	pool.printSocketMsg(m)

	for _, c := range pool.Clients {
		if c.ID == id1 {
			c.send(m1)
		} else {
			c.send(m)
		}
	}
}

// sends the message to the given client id
func (pool *Pool) sendToClientId(clientId string, m model.SocketMessage) {
	pool.printSocketMsg(m)

	for _, c := range pool.Clients {
		if c.ID == clientId {
			c.send(m)
			break
		}
	}
}

// broadcast the given message to all clients in pool
func (pool *Pool) broadcast(message model.SocketMessage) {
	pool.printSocketMsg(message)

	for _, c := range pool.Clients {
		c.send(message)
	}
}

// 10
func (pool *Pool) broadcastConfigs() {
	cfg := model.SharedConfig{
		MessageTypeMap:               messageTypeMap,
		TimeForEachWordInSeconds:     utils.DurationToSeconds(pool.DrawTime),
		TimeForChoosingWordInSeconds: utils.DurationToSeconds(TimeoutForChoosingWord),
		PrintLogs:                    debug,
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

// 33, 35
func (pool *Pool) broadcastWordList(words []string) {
	byteInfo, _ := json.Marshal(words)
	m1 := model.SocketMessage{
		Type:             33,
		TypeStr:          messageTypeMap[33],
		Content:          string(byteInfo),
		CurrSketcherId:   pool.CurrSketcher.ID,
		CurrSketcherName: pool.CurrSketcher.Name,
		TimeoutAfter:     utils.FormatTimeLong(time.Now().Add(TimeoutForChoosingWord)),
	}

	m := model.SocketMessage{
		Type:             35,
		TypeStr:          messageTypeMap[35],
		CurrSketcherName: pool.CurrSketcher.Name,
	}

	pool.sendCorrespondingMessages(pool.CurrSketcher.ID, m1, m)
}

// 8, 87, 88
func (pool *Pool) broadcastCurrentWordDetails() {
	m1 := model.SocketMessage{
		Type:              8,
		TypeStr:           messageTypeMap[8],
		CurrSketcherId:    pool.CurrSketcher.ID,
		CurrWord:          pool.CurrWord,
		CurrWordExpiresAt: utils.FormatTimeLong(pool.CurrWordExpiresAt),
	}

	m := model.SocketMessage{
		Type:              88,
		TypeStr:           messageTypeMap[88],
		CurrWordLen:       len(pool.CurrWord),
		CurrSketcherName:  pool.CurrSketcher.Name,
		CurrWordExpiresAt: utils.FormatTimeLong(pool.CurrWordExpiresAt),
	}

	// send sketcher is now drawing event to everyone except the sketcher
	pool.sendExcludingClientId(pool.CurrSketcher.ID, model.SocketMessage{
		Type:             87,
		TypeStr:          messageTypeMap[87],
		CurrSketcherName: pool.CurrSketcher.Name,
	})

	pool.sendCorrespondingMessages(pool.CurrSketcher.ID, m1, m)
}

// 81, 82
func (pool *Pool) broadcastTurnOver() {
	m1 := model.SocketMessage{
		Type:    81,
		TypeStr: messageTypeMap[81],
	}

	m := model.SocketMessage{
		Type:    82,
		TypeStr: messageTypeMap[82],
	}

	pool.sendCorrespondingMessages(pool.CurrSketcher.ID, m1, m)
}

// 83, 84
func (pool *Pool) broadcastTurnOverBeforeTimeout() {
	m1 := model.SocketMessage{
		Type:    83,
		TypeStr: messageTypeMap[83],
	}

	m := model.SocketMessage{
		Type:    84,
		TypeStr: messageTypeMap[84],
	}

	pool.sendCorrespondingMessages(pool.CurrSketcher.ID, m1, m)
}

// 32
func (pool *Pool) broadcastWordReveal(word string) {
	pool.broadcast(model.SocketMessage{
		Type:    32,
		TypeStr: messageTypeMap[32],
		Content: word,
	})
}
