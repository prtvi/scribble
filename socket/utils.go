package socket

import (
	"fmt"
	"os"
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

func printSocketMsg(m model.SocketMessage) {
	if !debug {
		return
	}

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

	utils.Cp("white",
		"from:", utils.Cs(msgTypeColor, fmt.Sprintf("%-15s ", from)),
		utils.Cs("white", "msg type: "), utils.Cs("red", fmt.Sprintf("%2d ", m.Type)),
		utils.Cs(msgTypeColor, messageTypeMap[m.Type]))
}

func sleep(d time.Duration) {
	time.Sleep(d)
}

func sleepWithInterrupt(d time.Duration, stop chan bool) bool {
	// this func can be used to sleep for d duration, with an interuppt if any to stop this sleep
	// to achieve this interrupt before timeout, pass a channel bool, which will be used to break this timeout
	// this chan needs to be used to pass acknowledgement for stopping this timeout
	// returns boolean whether the timeout was interrupted or not, if interrupted then returns true

	select {
	case <-stop:
		return true
	case <-time.After(d):
		return false
	}
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
		GameStartedAt:                 time.Time{},
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
		printHubStatus()
		sleep(DeletePoolAfterGameEndsDuration)

		for poolId, pool := range hub {
			// if pool exists and game has ended
			if pool != nil && pool.HasGameEnded {
				pool.printStats("Removing pool from hub, poolId:", poolId)
				delete(hub, poolId)
			}

			// if pool exists and game hasn't started for RemovePoolAfterGameNotStarted duration
			if now := time.Now(); now.Sub(pool.CreatedTime) > RemovePoolAfterGameNotStarted {
				pool.printStats("Removing junky pool, poolId:", poolId)
				delete(hub, poolId)
			}
		}

		printHubStatus()
	}
}

func printHubStatus() {
	if !debug || len(hub) == 0 {
		return
	}

	// HubSize
	utils.Cp("white", "HubSize:", utils.Cs("green", fmt.Sprintf("%d", len(hub))))

	for _, pool := range hub {
		pool.printStats()
	}
}

func DebugMode() {
	env := os.Getenv("ENV")
	if env == "" || env == "PROD" {
		return
	}

	debug = true
	utils.Cp("greenBg", "----------- DEV/DEBUG ENV -----------")

	GameStartDurationInSeconds = time.Second * 500
	TimeForEachWordInSeconds = time.Second * 30
	RenderClientsEvery = time.Second * 25

	poolId := "debug"
	pool := newPool(poolId, 4)
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", "/app?join="+poolId)

	hub[poolId] = pool

	go pool.start()
}
