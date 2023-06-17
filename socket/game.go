package socket

import (
	"time"
)

func (pool *Pool) start() {
	// start listening to pool connections and messages
	for {
		select {
		case client := <-pool.Register:
			// on client register, append the client to Pool.Clients slice, broadcast messageTypeMap, joining of the client and client info list
			pool.appendClientToList(client)
			pool.broadcastConfigs()
			pool.broadcastClientRegister(client.ID, client.Name)
			pool.broadcastClientInfoList()

			// start broadcasting client info list on first client join
			if len(pool.Clients) == 1 &&
				!pool.HasClientInfoBroadcastStarted &&
				!pool.HasGameStarted {

				// flag the client info broadcast start, run two sep goroutines to begin broadcasting client info at regular intervals and start game countdown
				pool.HasClientInfoBroadcastStarted = true
				go pool.beginBroadcastClientInfo()
				go pool.startGameCountdown()
			}

			pool.printStats("client connected, clientId:", client.ID)
			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice and broadcast the unregister
			pool.removeClientFromList(client)
			pool.broadcastClientUnregister(client.ID, client.Name)

			pool.printStats("client disconnected, clientId:", client.ID)
			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message
			// any of the game logic there is will be applied when clients do something, which will happen after the message is received from any of the clients

			switch message.Type {
			case 3:
				message := pool.updateScore(message)
				pool.broadcast(message)

			case 4, 5:
				pool.sendExcludingClientId(pool.CurrSketcher.ID, message)

			case 7:
				printSocketMsg(message)
				pool.startGameRequestFromClient()

			case 34:
				printSocketMsg(message)
				pool.CurrWord = message.Content // client choosing word

			default:
				break
			}

			break
		}
	}
}

func (pool *Pool) beginGameFlow() {
	// schedule timers for current word and current sketcher

	// wait for the "game started" overlay
	sleep(time.Second * 2)

	// loop over the number of rounds
	for i := 0; i < NumberOfRounds; i++ {
		pool.CurrRound = i + 1

		// broadcast round number and wait
		pool.broadcastRoundNumber()
		sleep(WaitAfterRoundStarts)

		// loop over all clients and assign words to each client and sleep until next client's turn
		for _, c := range pool.Clients {
			// broadcast clear canvas event, flag all clients as not guessed
			pool.broadcastClearCanvasEvent()
			pool.flagAllClientsAsNotGuessed()

			// begin client drawing flow and sleep until the word expires
			// broadcast current word, current sketcher and other details to all clients
			pool.clientWordAssignmentFlow(c)
			pool.broadcastCurrentWordDetails()

			// start a timer with interrupt, create a channel to use it to interrupt the timer if required
			// and run a go routine and pass this channel to pass data on this chan on all clients guess
			stopSleep := make(chan bool)
			go pool.checkIfAllGuessed(stopSleep)
			interrupted := sleepWithInterrupt(pool.CurrWordExpiresAt.Sub(time.Now()), stopSleep)

			// broadcast turn_over, reveal the word and clear canvas
			if interrupted {
				pool.broadcastTurnOverBeforeTimeout()
			} else {
				pool.broadcastTurnOver()
			}

			sleep(time.Second * 2)
			pool.broadcastWordReveal()
			sleep(time.Second * 2)
			pool.broadcastClearCanvasEvent()

			// flag sketching done, clear the current word and sketcher
			pool.turnOver(c)
		}

		pool.flagAllClientsAsNotSketched()
	}

	// once all clients are done playing, end the game and broadcast the same
	pool.endGame()
}
