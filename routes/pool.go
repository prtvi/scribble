package routes

import (
	"encoding/json"
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"strings"
	"time"
)

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

func (pool *Pool) startGameAndBroadcast() {
	// flag and broadcast the starting of the game
	pool.HasGameStarted = true
	pool.BroadcastMsg(model.SocketMessage{
		Type:    70,
		TypeStr: messageTypeMap[70],
		Success: true,
	})
}

func (pool *Pool) getClientInfoList() model.SocketMessage {
	// returns client info list embedded in model.SocketMessage

	clientInfoList := make([]model.ClientInfo, 0)

	// append client info into an array
	for _, client := range pool.Clients {
		clientInfoList = append(clientInfoList, model.ClientInfo{
			ID:    client.ID,
			Name:  client.Name,
			Color: client.Color,
			Score: client.Score,
		})
	}

	// marshall array in byte and send as string
	byteInfo, _ := json.Marshal(clientInfoList)
	return model.SocketMessage{
		Type:    6,
		TypeStr: messageTypeMap[6],
		Content: string(byteInfo),
	}
}

func (pool *Pool) appendClientToList(client *Client) {
	// append the client into the list
	pool.Clients = append(pool.Clients, client)

	// remove the color that was picked in GetColorForClient func from the list, the first color was picked from the list
	pool.ColorList[0] = pool.ColorList[len(pool.ColorList)-1]
	pool.ColorList = pool.ColorList[:len(pool.ColorList)-1]
}

func (pool *Pool) removeClientFromList(client *Client) {
	// take the removed client's color and append it to the color list
	pool.ColorList = append(pool.ColorList, client.Color)

	// remove the client from the list
	var idxToRemove int
	for i, c := range pool.Clients {
		if c == client {
			idxToRemove = i
			break
		}
	}

	pool.Clients[idxToRemove] = pool.Clients[len(pool.Clients)-1]
	pool.Clients = pool.Clients[:len(pool.Clients)-1]
}

func (pool *Pool) flagAllClientsAsNotGuessed() {
	for _, cl := range pool.Clients {
		cl.HasGuessed = false
	}
}

func printSocketMessage(m model.SocketMessage) {
	utils.Cp("cyan", "msg type:", utils.Cs("yellow", fmt.Sprintf("%d", m.Type)), utils.Cs("reset", messageTypeMap[m.Type], utils.Cs("cyan", "from:"), m.ClientName))
}

func (pool *Pool) GetColorForClient() string {
	return pool.ColorList[0]
}

func (pool *Pool) BroadcastMsg(message model.SocketMessage) {
	printSocketMessage(message)

	// broadcasts the given message to all clients in the pool
	for _, c := range pool.Clients {
		c.mu.Lock()
		err := c.Conn.WriteJSON(message)
		c.mu.Unlock()

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (pool *Pool) BeginBroadcastClientInfoMessage() {
	// to be run as a go routine
	// starts an infinite loop to broadcast client info after every regular interval
	for {
		time.Sleep(RenderClientsEvery)
		pool.BroadcastMsg(pool.getClientInfoList())

		// stop broadcasting when game ends
		if pool.HasGameEnded || len(pool.Clients) == 0 {
			utils.Cp("yellowBg", "Stopped broadcasting client info")
			break
		}
	}
}

func (pool *Pool) StartGameCountdown() {
	// as soon as the first player/client joins, start this countdown to start the game, after this timeout, the game begin message will broadcast

	// sleep until its the game starting time
	diff := pool.GameStartTime.Sub(time.Now())
	time.Sleep(diff)

	// if the game has already started by the client using the button then exit the countdown
	if pool.HasGameStarted {
		fmt.Println("game started using button so exiting countdown")
		return
	}

	// else start the game and broadcast the start game message
	pool.startGameAndBroadcast()
	utils.Cp("greenBg", "Game started! by server countdown")

	// start game flow
	go pool.BeginGameFlow()
}

func (pool *Pool) StartGameRequest() {
	// when the client requests to start the game instead of the countdown
	// start the game and broadcast the same
	pool.startGameAndBroadcast()
	utils.Cp("greenBg", "Game started! by client using btn")

	// start game flow
	go pool.BeginGameFlow()
}

func (pool *Pool) UpdateScore(message model.SocketMessage) model.SocketMessage {
	// update score for the client that guesses the word right, return true if correctly guessed

	var correctGuess, wordExistsInMessage bool

	// when the game has not begun, the curr sketcher will be nil
	if pool.CurrSketcher == nil {
		return message
	}

	// get the guesser client
	var guesserClient *Client = nil
	for _, c := range pool.Clients {
		// init guesserClient only if the guesser is not the sketcher
		if c.ID == message.ClientId &&
			pool.CurrSketcher.ID != message.ClientId {
			guesserClient = c
			break
		}
	}

	// if the sketcher is the guesser, then the guesserClient will be nil, hence check if guesserClient is nil
	// check if the word matches with the current word and check if the guesserClient hasn't already guessed
	if guesserClient != nil &&
		strings.ToLower(message.Content) == strings.ToLower(pool.CurrWord) &&
		!guesserClient.HasGuessed {

		correctGuess = true

		// increment score and flag as guessed
		guesserClient.HasGuessed = true
		guesserClient.Score += ScoreForCorrectGuess * int(utils.GetDiffBetweenTimesInSeconds(time.Now(), pool.CurrWordExpiresAt))

		// broadcast client info list to update score on UI immediately
		pool.BroadcastMsg(pool.getClientInfoList())
	}

	// check if the text message contains the word
	if strings.Contains(strings.ToLower(message.Content), strings.ToLower(pool.CurrWord)) {
		wordExistsInMessage = true
	}

	// if correct guess then modify the message
	if correctGuess {
		message.Type = 31
		message.TypeStr = messageTypeMap[31]
		message.Content = fmt.Sprintf("%s guessed the word!", message.ClientName)
	}

	// if word exists in the message
	if wordExistsInMessage {
		message.Type = 31
		message.TypeStr = messageTypeMap[31]
		message.Content = fmt.Sprintf("Naughty huh 😏 @%s", message.ClientName)
	}

	return message
}

func (pool *Pool) EndGame() {
	// flag and broadcast game end

	utils.Cp("yellowBg", "All players done playing!")

	pool.HasGameEnded = true
	pool.BroadcastMsg(model.SocketMessage{
		Type:    9,
		TypeStr: messageTypeMap[9],
		Content: pool.getClientInfoList().Content,
	})
}

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
				printSocketMessage(message)
				pool.StartGameRequest()

			case 34:
				printSocketMessage(message)
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
