package routes

import (
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

func PrintSocketMessage(m model.SocketMessage) {
	utils.Cp("cyan", "msg type:", utils.Cs("yellow", fmt.Sprintf("%d", m.Type)), utils.Cs("reset", messageTypeMap[m.Type], utils.Cs("cyan", "from:"), m.ClientName))
}

func NewPool(uuid string, capacity int) *Pool {
	// returns a new Pool
	now := time.Now()
	later := now.Add(GameStartDurationInSeconds)

	return &Pool{
		ID:                            uuid,
		JoiningLink:                   "",
		Capacity:                      capacity,
		CurrRound:                     1,
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
		time.Sleep(time.Minute * 10) // TODO - to be tested

		for poolId, pool := range HUB {
			if pool != nil && pool.HasGameEnded {
				utils.Cp("yellowBg", "Removing pool from HUB, poolId:", poolId)
				delete(HUB, poolId)

				fmt.Println("Size of HUB:", len(HUB))
			}

			if now := time.Now(); now.Sub(pool.CreatedTime) > time.Duration(time.Minute*10) {
				utils.Cp("yellowBg", "Removing pool from HUB after game not started for 10 mins, poolId:", poolId)
				delete(HUB, poolId)

				fmt.Println("Size of HUB:", len(HUB))
			}
		}
	}
}

func DebugMode() {
	GameStartDurationInSeconds = time.Duration(time.Second * 500)
	TimeForEachWordInSeconds = time.Duration(time.Second * 30)
	RenderClientsEvery = time.Duration(time.Second * 10)
	ScoreForCorrectGuess = 25
	NumberOfRounds = 3

	poolId := "debug"
	pool := NewPool(poolId, 4)

	HUB[poolId] = pool
	go pool.Start()

	link := "/app?join=" + poolId
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", link)

	utils.Cp("greenBg", "----------- DEBUG MODE -----------")
}
