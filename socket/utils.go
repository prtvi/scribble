package socket

import (
	"fmt"
	"os"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

func printSocketMsg(m model.SocketMessage) {
	from := m.ClientName
	if from == "" {
		from = "server"
	}

	var msgTypeColor string

	switch m.Type {
	case 35, 82:
		msgTypeColor = "cyan"
	case 8, 33, 81, 88:
		msgTypeColor = "yellow"
	case 1, 2, 3:
		msgTypeColor = "blue"
	case 4, 5:
		msgTypeColor = "purple"
	case 7, 34:
		msgTypeColor = "red"
	default:
		msgTypeColor = "green"
	}

	utils.Cp("black",
		"from:", utils.Cs(msgTypeColor, fmt.Sprintf("%-15s ", from)),
		utils.Cs("black", "msg type: "), utils.Cs("red", fmt.Sprintf("%2d ", m.Type)),
		utils.Cs(msgTypeColor, messageTypeMap[m.Type]))
}

func sleep(d time.Duration) {
	time.Sleep(d)
}

func newPool(uuid string, capacity int) *Pool {
	// returns a new Pool
	now := time.Now()
	later := now.Add(GameStartDurationInSeconds)

	return &Pool{
		ID:                            uuid,
		JoiningLink:                   "",
		Capacity:                      capacity,
		CurrRound:                     0,
		Register:                      make(chan *Client),
		Unregister:                    make(chan *Client),
		Clients:                       make([]*Client, 0),
		Broadcast:                     make(chan model.SocketMessage),
		ColorList:                     utils.ShuffleList(utils.COLORS[:10]),
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

func Maintainer() {
	// clears the pools in which the game has ended every 10 mins

	for {
		// TODO - to be tested
		sleep(DeletePoolAfterGameEndsDuration)

		for poolId, pool := range HUB {
			// if pool exists and game has ended
			if pool != nil && pool.HasGameEnded {
				utils.Cp("yellowBg", "Removing pool from HUB, poolId:", poolId)
				delete(HUB, poolId)

				fmt.Println("Size of HUB:", len(HUB))
			}

			// if pool exists and game hasn't started for RemovePoolAfterGameNotStarted duration
			if now := time.Now(); now.Sub(pool.CreatedTime) > RemovePoolAfterGameNotStarted {
				utils.Cp("yellowBg", "Removing pool from HUB after game not started for RemovePoolAfterGameNotStarted duration, poolId:", poolId)
				delete(HUB, poolId)

				fmt.Println("Size of HUB:", len(HUB))
			}
		}
	}
}

func DebugMode() {
	env := os.Getenv("ENV")
	if env == "" || env == "PROD" {
		return
	}

	DEBUG = true
	utils.Cp("greenBg", "----------- DEV/DEBUG ENV -----------")

	GameStartDurationInSeconds = time.Second * 500
	TimeForEachWordInSeconds = time.Second * 30
	RenderClientsEvery = time.Second * 10

	poolId := "debug"
	pool := newPool(poolId, 4)
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", "/app?join="+poolId)

	HUB[poolId] = pool

	go pool.start()
}
