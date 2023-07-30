package socket

import (
	"encoding/json"
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"strings"
	"time"
)

// read messages received from client
func (c *Client) read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgByte, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		// parse message received from client
		var clientMsg model.SocketMessage
		err = json.Unmarshal(msgByte, &clientMsg)

		// broadcast the message to all clients in the pool
		c.Pool.Broadcast <- clientMsg
	}
}

// send message to the associated client
func (c *Client) send(m model.SocketMessage) {
	c.mu.Lock()
	err := c.Conn.WriteJSON(m)
	c.mu.Unlock()

	if err != nil {
		fmt.Println(err)
	}
}

// returns client info list embedded in model.SocketMessage
func (pool *Pool) getClientInfoList() model.SocketMessage {
	clientInfoList := make([]model.ClientInfo, 0)

	// append client info into an array
	for _, client := range pool.Clients {
		clientInfoList = append(clientInfoList, model.ClientInfo{
			ID:           client.ID,
			Name:         client.Name,
			Score:        client.Score,
			AvatarConfig: client.AvatarConfig,
		})
	}

	// marshall array in byte and init as string
	byteInfo, _ := json.Marshal(clientInfoList)
	return model.SocketMessage{
		Type:    6,
		TypeStr: messageTypeMap[6],
		Content: string(byteInfo),
	}
}

// append the client to the clients list
func (pool *Pool) appendClientToList(client *Client) {
	pool.Clients = append(pool.Clients, client)
}

// remove the client from the client list
func (pool *Pool) removeClientFromList(client *Client) {
	// remove the client from the list
	var idx int
	for i, c := range pool.Clients {
		if c == client {
			idx = i
			break
		}
	}

	pool.Clients = append(pool.Clients[:idx], pool.Clients[idx+1:]...)
}

func (pool *Pool) flagAllClientsAsNotGuessed() {
	for _, cl := range pool.Clients {
		cl.HasGuessed = false
	}
}

func (pool *Pool) flagAllClientsAsNotSketched() {
	for _, cl := range pool.Clients {
		cl.DoneSketching = false
	}
}

// flag the client's turn as over and return the current word
func (pool *Pool) turnOver(c *Client) string {
	currWord := pool.CurrWord

	c.DoneSketching = true
	pool.CurrWord = ""
	pool.CurrSketcher = nil

	return currWord
}

// 70, flag and broadcast the starting of the game
func (pool *Pool) startGameAndBroadcast() {
	pool.HasGameStarted = true
	pool.GameStartedAt = time.Now()

	pool.broadcast(model.SocketMessage{
		Type:    70,
		TypeStr: messageTypeMap[70],
		Success: true,
	})
}

// begin the client's turn to draw, assign them the word automatically based on timeout if not chosen
func (pool *Pool) clientWordAssignmentFlow(client *Client) {
	// select the client
	pool.CurrSketcher = client

	// create a list of words for client to choose
	words := utils.GetNrandomWords(utils.WORDS, pool.WordCount)
	pool.broadcastWordList(words)

	// start a timeout for assigning word if not chosen by client
	go func() {
		// sleep until the duration, assign any random word to the client if timer runs out
		utils.Sleep(TimeoutForChoosingWord)

		if pool.CurrWord == "" {
			pool.CurrWord = utils.GetRandomItem(words)
		}
	}()

	// run an infinite loop until pool.CurrWord is initialised by sketcher client (initialised in pool.Start func), or initialised in word choose countdown goroutine
	for pool.CurrWord == "" {
	}

	// add the word expiry
	pool.CurrWordExpiresAt = time.Now().Add(pool.DrawTime)
	// reinit hints revealed
	pool.HintsRevealed = 0
}

// begin clientInfo broadcast
func (pool *Pool) beginBroadcastClientInfo() {
	// to be run as a go routine
	// starts an infinite loop to broadcast client info after every regular interval
	pool.HasClientInfoBroadcastStarted = true
	utils.Cp("yellow", "Broadcasting client info start!")

	for {
		utils.Sleep(RenderClientsEvery)
		pool.broadcastClientInfoList()

		// stop broadcasting when game ends
		if pool.HasGameEnded || len(pool.Clients) == 0 {
			utils.Cp("yellow", "Stopped broadcasting client info")
			pool.HasClientInfoBroadcastStarted = false
			break
		}
	}
}

// begin the game flow as soon as a client requests to start the game
func (pool *Pool) startGameRequestFromClient(clientId string) {
	// when the client requests to start the game instead of the countdown
	// start the game and broadcast the same

	if len(pool.Clients) < 2 {
		pool.sendToClientId(clientId, model.SocketMessage{
			Type:    69,
			TypeStr: messageTypeMap[69],
		})

		return
	}

	pool.startGameAndBroadcast()
	utils.Cp("greenBg", "Game started!")

	// start game flow
	go pool.beginGameFlow()
}

// update score for the client that guesses the word right
func (pool *Pool) updateScore(message model.SocketMessage) model.SocketMessage {
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

	guessedLower := strings.ToLower(message.Content)
	currWordLower := strings.ToLower(pool.CurrWord)

	// if guesserClient == nil then its the sketcher sending the message, then modify the message if the sketcher tries to reveal the word
	if guesserClient == nil &&
		(guessedLower == currWordLower || strings.Contains(guessedLower, currWordLower)) {

		message.Type = 313
		message.TypeStr = messageTypeMap[313]
		return message
	}

	// if the sketcher is the guesser, then the guesserClient will be nil, hence check if guesserClient is nil
	// check if the word matches with the current word and check if the guesserClient hasn't already guessed
	if guesserClient != nil &&
		!guesserClient.HasGuessed &&
		guessedLower == currWordLower {

		// increment score and flag as guessed
		guesserClient.HasGuessed = true
		guesserClient.Score += utils.CalcScore(ScoreForCorrectGuess, pool.CurrRound, pool.CurrWordExpiresAt)

		// broadcast client info list to update score on UI immediately
		pool.broadcastClientInfoList()

		// if correct guess then modify the message
		message.Type = 31
		message.TypeStr = messageTypeMap[31]
		return message
	}

	// check if the text message contains the word, word exists in message
	// send this response only if client has already guessed the current word
	if guesserClient != nil &&
		guesserClient.HasGuessed &&
		strings.Contains(guessedLower, currWordLower) {

		message.Type = 312
		message.TypeStr = messageTypeMap[312]
	}

	return message
}

// checks if all the clients have guessed the word and acknowledges it on the stopTimer channel
func (pool *Pool) checkIfAllGuessed(stopTimer, stopHints chan bool) {
	// to be run as a separate goroutine
	// every second, check if all clients have guessed the word
	// if yes, then acknowledge the same on the channel and break this loop
	for {
		utils.Sleep(time.Second * 1)

		var count int = 0
		for _, c := range pool.Clients {
			if c.HasGuessed {
				count += 1
			}
		}

		// if gussed clients is everyone except the sketcher
		if count != 0 && count == len(pool.Clients)-1 {
			stopTimer <- true // write to channel and break
			stopHints <- true
			break
		}

		// if current sketcher is reset/done sketching then break
		if pool.CurrSketcher == nil || pool.CurrSketcher.DoneSketching {
			break
		}
	}
}

func (pool *Pool) broadcastHintsForWord(stopHints chan bool) {
	pool.HintsForCurrWord = utils.CalculateMaxHintsAllowedForWord(pool.CurrWord, pool.Hints)
	revealDurationParts := pool.HintsForCurrWord + 2
	revealHintsEvery := time.Duration(utils.DurationToSeconds(pool.DrawTime)/revealDurationParts) * time.Second

	var word string = pool.CurrWord
	var charsLeft []string = strings.Split(word, "")
	var charPicked string
	var hintString string = func(word string) string {
		var res string
		for i := 0; i < len(word); i++ {
			res += "_"
		}
		return res
	}(word)

	go func() {
		for pool.HintsRevealed < pool.HintsForCurrWord {
			interrupted := utils.SleepWithInterrupt(revealHintsEvery, stopHints)
			if interrupted {
				break
			}

			charsLeft, charPicked = utils.PickRandomCharacter(charsLeft)
			hintString = utils.GetHintString(word, charPicked, hintString)

			pool.sendExcludingClientId(pool.CurrSketcher.ID, model.SocketMessage{
				Type:    89,
				TypeStr: messageTypeMap[89],
				Content: hintString,
			})

			pool.HintsRevealed += 1
		}
	}()
}

// 9, flag and broadcast game end
func (pool *Pool) endGame() {
	utils.Cp("greenBg", "All players done playing!")

	pool.HasGameEnded = true
	pool.broadcast(model.SocketMessage{
		Type:    9,
		TypeStr: messageTypeMap[9],
		Content: pool.getClientInfoList().Content,
	})
}
