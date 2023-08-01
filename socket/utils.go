package socket

import (
	"fmt"
	"os"
	model "scribble/model"
	utils "scribble/utils"
	"text/tabwriter"
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
	case 35, 82, 84, 87, 88:
		msgTypeColor = "cyan"
	case 69, 8, 33, 81, 83:
		msgTypeColor = "yellow"
	case 1, 2, 3:
		msgTypeColor = "blue"
	case 4, 41, 5:
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

func newPool(players, drawTime, rounds, wordCount, hints int, wordMode string, customWords []string, useCustomWordsOnly bool) *Pool {
	return &Pool{
		ID:        utils.GenerateUUID(),
		Capacity:  players,
		DrawTime:  time.Duration(time.Second * time.Duration(drawTime)),
		Rounds:    rounds,
		WordCount: wordCount,
		Hints:     hints,

		// not implemented yet
		WordMode:           wordMode,
		CustomWords:        customWords,
		UseCustomWordsOnly: useCustomWordsOnly,

		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Clients:        make([]*Client, 0),
		Broadcast:      make(chan model.SocketMessage),
		CreatedTime:    time.Now(),
		HasGameStarted: false,
	}
}

func Maintainer() {
	// clears the pools in which the game has ended every 10 mins

	for {
		// TODO - to be tested
		// printHubStatus()
		utils.Sleep(DeletePoolAfterGameEndsDuration)

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

		// printHubStatus()
	}
}

func printHubStatus() {
	if !debug || len(hub) == 0 {
		return
	}

	// HubSize
	utils.Cp("white", "HubSize:", utils.Cs("green", fmt.Sprintf("%d", len(hub))))

	// for _, pool := range hub {
	// pool.printStats()
	// }
}

func DebugMode() {
	env := os.Getenv("ENV")
	if env == "" || env == "PROD" {
		return
	}

	debug = true
	utils.Cp("greenBg", "----------- DEV/DEBUG ENV -----------")

	RenderClientsEvery = time.Second * 30

	pool := newPool(5, 30, 6, 5, 2, "normal", []string{}, false)
	pool.ID = "debug"
	pool.JoiningLink = fmt.Sprintf("localhost:1323%s", "/app?join="+pool.ID)

	hub[pool.ID] = pool

	go pool.start()
}

func (pool *Pool) printStats(event ...string) {
	if !debug {
		return
	}

	// poolId, capacity, size, createdTime, gameStartTime
	// hasGameStarted, hasGameEnded, hasClientInfoBroadcastStarted
	// currRound, currWord, currSketcherName, currWordExpiresAt

	utils.Cp("green", "--------------------------------------------------------------")
	if len(event) != 0 {
		utils.Cp("red", event...)
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 2, 2, ' ', 0)
	fmt.Fprintln(w, "\nPoolId\tCapacity\tSize\tCreatedTime\tGameStartTime\tGameStartedAt")
	fmt.Fprintln(w, fmt.Sprintf("%s\t%d\t%d\t%s\t%s\t%s\n", pool.ID, pool.Capacity, len(pool.Clients), utils.GetTimeString(pool.CreatedTime), utils.GetTimeString(pool.CreatedTime), utils.GetTimeString(pool.GameStartedAt)))

	fmt.Fprintln(w, "HasGameStarted\tHasClientInfoBroadcastStarted\tHasGameEnded")
	fmt.Fprintln(w, fmt.Sprintf("%v\t%v\t%v\n", pool.HasGameStarted, pool.HasClientInfoBroadcastStarted, pool.HasGameEnded))

	if pool.CurrSketcher != nil && pool.HasGameStarted {
		fmt.Fprintln(w, "CurrRound\tCurrWord\tCurrSketcherName")
		fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%v\n", pool.CurrRound, pool.CurrWord, pool.CurrSketcher.Name))
	}

	w.Flush()
	utils.Cp("green", "--------------------------------------------------------------")
}
