package socket

import (
	model "scribble/model"
	utils "scribble/utils"
	"time"
)

// start listening to pool connections and messages
func (pool *Pool) start() {
	for {
		select {
		case client := <-pool.Register:
			// on client register, append the client to Pool.Clients slice, broadcast messageTypeMap, joining of the client and client info list
			pool.appendClientToList(client)
			pool.broadcastConfigs()
			pool.broadcastClientRegister(client.ID, client.Name)
			pool.broadcastClientInfoList()

			if pool.HasGameStarted {
				// to trigger event 70 for mid game joinee,
				client.send(model.SocketMessage{
					Type:          70,
					TypeStr:       messageTypeMap[70],
					Success:       true,
					Content:       "Game has started!",
					MidGameJoinee: true,
				})

				// trigger round info,
				client.send(model.SocketMessage{
					Type:          71,
					TypeStr:       messageTypeMap[71],
					CurrRound:     pool.CurrRound,
					MidGameJoinee: true,
				})

				if pool.SleepingForSketching {
					// trigger word stat
					client.send(model.SocketMessage{
						Type:    89,
						TypeStr: messageTypeMap[89],
						Content: pool.HintString,
					})

					// trigger timer info
					client.send(model.SocketMessage{
						Type:              86,
						TypeStr:           messageTypeMap[86],
						CurrWordExpiresAt: utils.FormatTimeLong(pool.CurrWordExpiresAt),
					})
				}
			}

			// start broadcasting client info list on first client join
			if len(pool.Clients) == 1 &&
				!pool.HasClientInfoBroadcastStarted &&
				!pool.HasGameStarted {

				// run sep goroutine to begin broadcasting client info at regular intervals
				go pool.beginBroadcastClientInfo()
			}

			// pool.printStats("client connected, clientId:", client.ID)
			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice and broadcast the unregister
			pool.removeClientFromList(client)
			pool.broadcastClientUnregister(client.ID, client.Name)
			pool.broadcastClientInfoList()

			// pool.printStats("client disconnected, clientId:", client.ID)
			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message
			// any of the game logic there is will be applied when clients do something, which will happen after the message is received from any of the clients

			switch message.Type {
			case 3:
				message := pool.updateScore(message)
				pool.broadcast(message)

			case 4, 41, 5:
				if pool.CurrSketcher == nil {
					break
				}

				pool.sendExcludingClientId(pool.CurrSketcher.ID, message) // avoid sending canvas data and clear canvas event to the curr sketcher

			case 7:
				pool.printSocketMsg(message)
				pool.startGameRequestFromClient(message.ClientId)

			case 34:
				pool.printSocketMsg(message)
				pool.CurrWord = message.Content // client choosing word

			default:
				break
			}

			break
		}
	}
}

// begin game flow by scheduling schedule timers
func (pool *Pool) beginGameFlow() {
	// wait for the "game started" overlay
	utils.Sleep(InterGameWaitDuration)

	// loop over the number of rounds
	for i := 0; i < pool.Rounds; i++ {
		pool.CurrRound = i + 1

		// broadcast round number and wait
		pool.broadcastRoundNumber()
		utils.Sleep(InterGameWaitDuration)

		// loop over all clients and assign words to each client and sleep until next client's turn
		for _, c := range pool.Clients {
			// flag all clients as not guessed
			pool.flagAllClientsAsNotGuessed()

			// begin client drawing flow and sleep until the word expires
			// broadcast current word, current sketcher and other details to all clients
			pool.clientWordAssignmentFlow(c)
			pool.broadcastCurrentWordDetails()
			pool.broadcastClientInfoList()

			stopRevealingHints := make(chan bool)
			pool.broadcastHintsForWord(stopRevealingHints)

			// start a timer with interrupt, create a channel to use it to interrupt the timer if required
			// and run a go routine and pass this channel to pass data on this chan on all clients guess
			stopSketchingTime := make(chan bool)
			go pool.checkIfAllGuessed(stopSketchingTime, stopRevealingHints)

			pool.SleepingForSketching = true
			interrupted := utils.SleepWithInterrupt(pool.CurrWordExpiresAt.Sub(time.Now()), stopSketchingTime)
			pool.SleepingForSketching = false

			// broadcast turn_over, reveal the word and clear canvas
			if interrupted {
				pool.broadcastTurnOverBeforeTimeout()
			} else {
				pool.broadcastTurnOver()
			}

			// flag sketching done, clear the current word and sketcher
			currWord := pool.turnOver(c)

			utils.Sleep(InterGameWaitDuration)
			pool.broadcastWordReveal(currWord)

			utils.Sleep(InterGameWaitDuration)
			pool.broadcastClearCanvasEvent()
		}

		pool.flagAllClientsAsNotSketched()
	}

	// once all clients are done playing, end the game and broadcast the same
	pool.endGame()
}
