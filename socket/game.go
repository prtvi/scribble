package socket

import (
	"fmt"
	utils "scribble/utils"
	"time"
)

func (pool *Pool) Start() {
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
				go pool.BeginBroadcastClientInfo()
				go pool.StartGameCountdown()
			}

			utils.Cp("yellow", "Size of pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client connected:"), client.Name)
			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice and broadcast the unregister
			pool.removeClientFromList(client)
			pool.broadcastClientUnregister(client.ID, client.Name)

			utils.Cp("yellow", "Size of pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client disconnected:"), client.Name)
			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message
			// any of the game logic there is will be applied when clients do something, which will happen after the message is received from any of the clients

			switch message.Type {
			case 3:
				message := pool.UpdateScore(message)
				pool.broadcast(message)

			case 4, 5:
				message.CurrSketcherId = pool.CurrSketcher.ID // to disable redrawing on sketcher's canvas
				pool.broadcast(message)

			case 7:
				PrintSocketMessage(message)
				pool.StartGameRequest()

			case 34:
				PrintSocketMessage(message)
				pool.CurrWord = message.Content // client choosing word

			default:
				break
			}

			break
		}
	}
}

func (pool *Pool) BeginGameFlow() {
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
			// TODO: send the whole thing to client who's sketching, send minimal details to rest
			pool.clientWordAssignmentFlow(c)
			pool.broadcastCurrentWordDetails()

			sleep(pool.CurrWordExpiresAt.Sub(time.Now()))

			// broadcast turn_over, reveal the word and clear canvas
			pool.broadcastTurnOver() // TODO: show these events on overlay
			sleep(time.Second * 2)
			pool.broadcastWordReveal()
			sleep(time.Second * 2)
			pool.broadcastClearCanvasEvent()

			// flag sketching done, clear the current word and sketcher
			pool.turnOver(c)
			sleep(WaitAfterTurnEnds)
		}

		pool.flagAllClientsAsNotSketched()
	}

	// once all clients are done playing, end the game and broadcast the same
	pool.EndGame()
}
