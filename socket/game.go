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
			// on client register, append the client to Pool.Clients slice
			pool.appendClientToList(client)

			// send the messageTypeMap to clients
			pool.broadcastMessageTypeMap()

			// broadcast the joining of client
			pool.broadcastClientRegister(client.ID, client.Name)

			// send client info list once client joins
			pool.broadcastClientInfoList()

			// start broadcasting client info list on first client join
			if len(pool.Clients) == 1 &&
				!pool.HasClientInfoBroadcastStarted &&
				!pool.HasGameStarted {

				pool.HasClientInfoBroadcastStarted = true
				utils.Cp("yellowBg", "Broadcasting client info start!")

				// begin braodcasting client info at regular intervals
				go pool.BeginBroadcastClientInfo()

				// begin start game countdown
				go pool.StartGameCountdown()
			}

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client connected:"), client.Name)

			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice
			pool.removeClientFromList(client)

			// broadcast the leaving of client
			pool.broadcastClientUnregister(client.ID, client.Name)

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client disconnected:"), client.Name)

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
		}
	}
}

func (pool *Pool) BeginGameFlow() {
	// schedule timers for current word and current sketcher

	// loop over the number of rounds
	for i := 0; i < NumberOfRounds; i++ {
		pool.CurrRound = i + 1

		// broadcast round number
		pool.broadcastRoundNumber()

		time.Sleep(WaitAfterRoundStarts)

		// loop over all clients and assign words to each client and sleep until next client's turn
		for _, c := range pool.Clients {

			// broadcast clear canvas event
			pool.broadcastClearCanvasEvent()

			// flag all clients as not guessed
			pool.flagAllClientsAsNotGuessed()

			// select the client
			pool.CurrSketcher = c
			c.HasSketched = true

			// create a list of words for client to choose
			words := utils.Get3RandomWords(utils.WORDS)
			pool.broadcast3WordsList(words)

			// start a timeout for assigning word if not chosen by client
			go pool.wordChooseCountdown(words)

			// run an infinite loop until pool.CurrWord is initialised by sketcher client, initialised in pool.Start func, TODO: create a timeout instead
			for pool.CurrWord == "" {
			}

			// add the word expiry
			pool.CurrWordExpiresAt = time.Now().Add(TimeForEachWordInSeconds)

			// broadcast current word, current sketcher and other details to all clients
			// TODO: send the whole thing to client who's sketching, send minimal details to rest
			pool.broadcastCurrentWordDetails()

			// sleep until the word expires
			time.Sleep(pool.CurrWordExpiresAt.Sub(time.Now()))

			// broadcast turn_over
			pool.broadcastTurnOver()

			// reveal the word
			pool.broadcastWordReveal()

			pool.CurrWord = ""
			pool.CurrSketcher = nil

			time.Sleep(WaitAfterTurnEnds)
		}
	}

	// once all clients are done playing, end the game and broadcast the same
	pool.EndGame()
}
