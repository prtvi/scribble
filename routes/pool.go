package routes

import (
	"encoding/json"
	"fmt"
	model "scribble/model"
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
			byteInfo, _ := json.Marshal(messageTypeMap)
			pool.BroadcastMsg(model.SocketMessage{
				Type:    10,
				TypeStr: messageTypeMap[10],
				Content: string(byteInfo),
			})

			// broadcast the joining of client
			pool.BroadcastMsg(model.SocketMessage{
				Type:       1,
				TypeStr:    messageTypeMap[1],
				ClientId:   client.ID,
				ClientName: client.Name,
			})

			// send client info list once client joins
			pool.BroadcastMsg(pool.getClientInfoList())

			// start broadcasting client info list on first client join
			if len(pool.Clients) == 1 &&
				!pool.HasClientInfoBroadcastStarted &&
				!pool.HasGameStarted {

				pool.HasClientInfoBroadcastStarted = true
				utils.Cp("yellowBg", "Broadcasting client info start!")

				// begin braodcasting client info at regular intervals
				go pool.BeginBroadcastClientInfoMessage()

				// begin start game countdown
				go pool.StartGameCountdown()
			}

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client connected:"), client.Name)

			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice
			pool.removeClientFromList(client)

			// broadcast the leaving of client
			pool.BroadcastMsg(model.SocketMessage{
				Type:       2,
				TypeStr:    messageTypeMap[2],
				ClientId:   client.ID,
				ClientName: client.Name,
			})

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client disconnected:"), client.Name)

			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message
			// any of the game logic there is will be applied when clients do something, which will happen after the message is received from any of the clients

			switch message.Type {
			case 3:
				message := pool.UpdateScore(message)
				pool.BroadcastMsg(message)

			case 4, 5:
				message.CurrSketcherId = pool.CurrSketcher.ID // to disable redrawing on sketcher's canvas
				pool.BroadcastMsg(message)

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
		pool.BroadcastMsg(model.SocketMessage{
			Type:      71,
			TypeStr:   messageTypeMap[71],
			CurrRound: pool.CurrRound,
		})

		time.Sleep(WaitAfterRoundStarts)

		// loop over all clients and assign words to each client and sleep until next client's turn
		for _, c := range pool.Clients {

			// broadcast clear canvas event
			pool.BroadcastMsg(model.SocketMessage{
				Type:    5,
				TypeStr: messageTypeMap[5],
			})

			// flag all clients as not guessed
			pool.flagAllClientsAsNotGuessed()

			// select the client
			pool.CurrSketcher = c
			c.HasSketched = true

			// create a list of words for client to choose
			words := utils.Get3RandomWords(utils.WORDS)

			byteInfo, _ := json.Marshal(words)
			pool.BroadcastMsg(model.SocketMessage{
				Type:             33,
				TypeStr:          messageTypeMap[33],
				Content:          string(byteInfo),
				CurrSketcherId:   pool.CurrSketcher.ID,
				CurrSketcherName: pool.CurrSketcher.Name,
			})

			// start a timeout for assigning word if not chosen by client
			go func() {
				time.Sleep(TimeoutForChoosingWord)

				if pool.CurrWord == "" {
					fmt.Println("auto assigned")
					pool.CurrWord = utils.GetRandomWord(words)
					return
				}

				fmt.Println("exiting timeout wo auto assignment")
			}()

			// run an infinite loop until pool.CurrWord is initialised by sketcher client, initialised in pool.Start func, TODO: create a timeout instead
			for pool.CurrWord == "" {
			}

			// add the word expiry
			pool.CurrWordExpiresAt = time.Now().Add(TimeForEachWordInSeconds)

			// broadcast current word, current sketcher and other details to all clients
			// TODO: send the whole thing to client who's sketching, send minimal details to rest
			pool.BroadcastMsg(model.SocketMessage{
				Type:              8,
				TypeStr:           messageTypeMap[8],
				CurrSketcherId:    pool.CurrSketcher.ID,
				CurrWord:          pool.CurrWord,
				CurrWordExpiresAt: utils.FormatTimeLong(pool.CurrWordExpiresAt),
			})

			// sleep until the word expires
			time.Sleep(pool.CurrWordExpiresAt.Sub(time.Now()))

			// broadcast turn_over
			pool.BroadcastMsg(model.SocketMessage{
				Type:           81,
				TypeStr:        messageTypeMap[81],
				CurrSketcherId: pool.CurrSketcher.ID,
			})

			// reveal the word
			pool.BroadcastMsg(model.SocketMessage{
				Type:    32,
				TypeStr: messageTypeMap[32],
				Content: fmt.Sprintf("%s was the correct word!", pool.CurrWord),
			})

			pool.CurrWord = ""
			pool.CurrSketcher = nil

			time.Sleep(WaitAfterTurnEnds)
		}
	}

	// once all clients are done playing, end the game and broadcast the same
	pool.EndGame()
}
