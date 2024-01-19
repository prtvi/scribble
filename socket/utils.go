package socket

import (
	"runtime"
	model "scribble/model"
	utils "scribble/utils"
	"time"

	"github.com/gorilla/websocket"
)

func (pool *Pool) printSocketMsg(m model.SocketMessage) {
	if !debug || m.Type == 4 {
		return
	}

	from := m.ClientName
	if from == "" {
		from = "server"
	}

	var color string
	switch m.Type {
	case 35, 82, 84, 87, 88, 89:
		color = "cyan"
	case 69, 8, 33, 81, 83:
		color = "yellow"
	case 1, 2, 3:
		color = "blue"
	case 4, 41, 5:
		color = "purple"
	case 7, 34:
		color = "red"
	default:
		color = "green"
	}

	// utils.Cp("reset", "SOCKET_MSG> pool id:", utils.Cs("green", pool.ID),
	// 	"from:", utils.Cs(color, fmt.Sprintf("%-15s ", from)),
	// 	"msg type:", utils.Cs("red", fmt.Sprintf("%2d ", m.Type)),
	// 	utils.Cs(color, messageTypeMap[m.Type]))
	utils.Cp(color, pool.ID, "SOCKET_MSG>", "from:", from, "msg type:", m.Type, messageTypeMap[m.Type])
}

func newPool(players, drawTime, rounds, wordCount, hints int, wordMode string, customWords []string, useCustomWordsOnly bool) *Pool {
	return &Pool{
		ID:                 utils.GenerateUUID(),
		Capacity:           players,
		DrawTime:           time.Duration(time.Second * time.Duration(drawTime)),
		Rounds:             rounds,
		WordCount:          wordCount,
		Hints:              hints,
		WordMode:           wordMode,
		CustomWords:        customWords,
		UseCustomWordsOnly: useCustomWordsOnly,
		WordsForGame:       utils.SelectWordsForPool(utils.WORDS, customWords, useCustomWordsOnly),
		InitCurrWord:       make(chan string),
		Register:           make(chan *Client),
		Unregister:         make(chan *Client),
		Clients:            make([]*Client, 0),
		Broadcast:          make(chan model.SocketMessage),
		CreatedTime:        time.Now(),
		HasGameStarted:     false,
	}
}

func newClient(id, name string, conn *websocket.Conn, pool *Pool, ac model.AvatarConfig) *Client {
	return &Client{
		ID:            id,
		Name:          name,
		AvatarConfig:  ac,
		IsSketching:   false,
		DoneSketching: false,
		HasGuessed:    false,
		Score:         0,
		Conn:          conn,
		Pool:          pool,
	}
}

func Maintainer() {
	utils.Cp("yellow", "started maintainer")

	// clears the pools in which the game has ended
	// can be implemented using channel
	for {
		utils.Sleep(DurationToCheckForRemovingJunkyPools)

		for poolId, pool := range hub {
			// if pool exists and game has ended
			if pool != nil && pool.HasGameEnded {
				utils.Cp("red", "game ended, removing pool from hub, poolId:", poolId)
				pool = nil

				delete(hub, poolId)
			}
		}

		for poolId, pool := range hub {
			// if pool exists and game hasn't started for RemovePoolAfterGameNotStarted duration
			if now := time.Now(); pool != nil && now.Sub(pool.CreatedTime) > RemovePoolAfterGameNotStarted {
				utils.Cp("red", "removing junky pool, poolId:", poolId)
				pool.HasGameEnded = true
				pool = nil

				delete(hub, poolId)
			}
		}

		utils.Cp("greenBg", "Number of goroutines running:", runtime.NumGoroutine())
		utils.Cp("red", "len hub:", len(hub))
	}
}

func InitDebugEnv(isDebugEnv bool) {
	debug = isDebugEnv
	utils.Cp("red", "len hub:", len(hub))

	if isDebugEnv {
		utils.Cp("greenBg", "----------- DEV/DEBUG ENV -----------")
	} else {
		utils.Cp("redBg", "----------- PROD ENV -----------")
	}
}
