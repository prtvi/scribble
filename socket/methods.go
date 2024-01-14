package socket

import (
	"encoding/json"
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"sort"
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
		if err != nil {
			fmt.Println(err)
		}

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
func (pool *Pool) getClientInfoList(finalCall bool) model.SocketMessage {
	clientInfoList := make([]model.ClientInfo, 0)

	// append client info into an array
	for _, client := range pool.Clients {
		clientInfoList = append(clientInfoList, model.ClientInfo{
			ID:           client.ID,
			Name:         client.Name,
			Score:        client.Score,
			IsSketching:  client.IsSketching,
			AvatarConfig: client.AvatarConfig,
		})
	}

	// sort descending wrt score
	sort.Slice(clientInfoList, func(i, j int) bool {
		return clientInfoList[i].Score > clientInfoList[j].Score
	})

	if finalCall {
		// crown the player with highest score
		if len(clientInfoList) > 0 && clientInfoList[0].Score > 0 {
			clientInfoList[0].AvatarConfig.IsCrowned = true
		}
	}

	// marshall array in byte and init as string
	byteInfo, _ := json.Marshal(clientInfoList)
	return model.SocketMessage{
		Type:    6,
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
	utils.Cp("yellow", pool.ID, "flagging all clients as not guessed")
	for _, c := range pool.Clients {
		c.HasGuessed = false
	}
}

func (pool *Pool) flagAllClientsAsNotSketched() {
	utils.Cp("yellow", pool.ID, "flagging all clients as not sketched")
	for _, c := range pool.Clients {
		c.DoneSketching = false
	}
}

// flag the client's turn as over and return the current word
func (pool *Pool) turnOver(c *Client) string {
	utils.Cp("blue", pool.ID, "turn over")
	currWord := pool.CurrWord

	c.IsSketching = false
	c.DoneSketching = true
	pool.CurrWord = ""
	pool.CurrSketcher = nil

	return currWord
}

// 70, flag and broadcast the starting of the game
func (pool *Pool) startGameAndBroadcast() {
	utils.Cp("green", pool.ID, "start game broadcasting")
	pool.HasGameStarted = true
	pool.GameStartedAt = time.Now()

	pool.broadcast(model.SocketMessage{
		Type:    70,
		Success: true,
		Content: "Game started!",
	})
}

// begin the client's turn to draw, assign them the word automatically based on timeout if not chosen
func (pool *Pool) clientWordAssignmentFlow(client *Client) {
	utils.Cp("cyan", pool.ID, "client word assignment flow has begun")
	// select the client
	pool.CurrSketcher = client
	client.IsSketching = true

	// create a list of words for client to choose
	words := utils.GetNrandomWords(utils.WORDS, pool.WordCount)
	pool.broadcastWordList(words)

	// start a timeout for assigning word if not chosen by client
	go func() {
		// sleep until the duration, assign any random word to the client if timer runs out
		utils.Sleep(TimeoutForChoosingWord)

		if pool.CurrWord == "" {
			utils.Cp("cyan", pool.ID, "word assigned after timeout")
			pool.InitCurrWord <- utils.GetRandomItem(words)
		}
	}()

	// wait until pool.CurrWord is initialised by sketcher client (initialised in pool.Start func, case: 34), or initialised in word choose countdown goroutine above
	pool.CurrWord = <-pool.InitCurrWord

	// add the word expiry
	pool.CurrWordExpiresAt = time.Now().Add(pool.DrawTime)
	// reinit hints revealed
	pool.NumHintsRevealed = 0
	// clear hint string
	pool.HintString = ""
}

// begin clientInfo broadcast
func (pool *Pool) beginBroadcastClientInfo() {
	// to be run as a go routine
	// starts an infinite loop to broadcast client info after every regular interval
	pool.HasClientInfoBroadcastStarted = true
	utils.Cp("yellow", pool.ID, "broadcasting client info start")

	for {
		utils.Sleep(RenderClientsEvery)
		pool.broadcastClientInfoList()

		// stop broadcasting when game ends
		if pool.HasGameEnded || len(pool.Clients) == 0 {
			utils.Cp("yellow", pool.ID, "stopped broadcasting client info")
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
		pool.sendToClientId(clientId, model.SocketMessage{Type: 69})
		return
	}

	pool.startGameAndBroadcast()
	utils.Cp("greenBg", pool.ID, "game started by client")

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
		return message
	}

	// if the sketcher is the guesser, then the guesserClient will be nil, hence check if guesserClient is nil
	// check if the text message contains the word, word exists in message
	// send this response only if client has already guessed the current word
	if guesserClient != nil &&
		guesserClient.HasGuessed &&
		strings.Contains(guessedLower, currWordLower) {

		message.Type = 312
	}

	return message
}

// checks if all the clients have guessed the word and acknowledges it on the stopTimer channel
func (pool *Pool) checkIfAllGuessed(stopSketching, stopHints chan bool) {
	utils.Cp("yellow", pool.ID, "entered checking if all have guessed flow")
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
			utils.Cp("yellow", pool.ID, "all clients have guessed, breaking from loop")
			// write to channel and break
			stopSketching <- true
			if pool.WordMode == "normal" {
				stopHints <- true
			}

			break
		}

		// if current sketcher is reset/done sketching then break
		if pool.CurrSketcher == nil || pool.CurrSketcher.DoneSketching {
			break
		}
	}
	utils.Cp("yellow", pool.ID, "exited check if all guessed loop")
}

func (pool *Pool) broadcastHintsForWord(stopHints chan bool) {
	utils.Cp("purple", pool.ID, "entered broadcasting hints flow")

	pool.NumHintsForCurrWord = utils.CalculateMaxHintsAllowedForWord(pool.CurrWord, pool.Hints)
	revealDurationParts := pool.NumHintsForCurrWord + 2
	revealHintsEvery := time.Duration(utils.DurationToSeconds(pool.DrawTime)/revealDurationParts) * time.Second

	var word string = pool.CurrWord
	var charsLeft []string = strings.Split(word, "")
	var charPicked string
	pool.HintString = func(word string) string {
		var res string
		for i := 0; i < len(word); i++ {
			res += "_"
		}
		return res
	}(word)

	utils.Cp("blue", pool.ID, "nHints:", pool.NumHintsForCurrWord)

	go func() {
		for pool.NumHintsRevealed < pool.NumHintsForCurrWord {
			interrupted := utils.SleepWithInterrupt(revealHintsEvery, stopHints)
			if interrupted {
				utils.Cp("purple", pool.ID, "interrupted hints reveal timer")
				break
			}

			charsLeft, charPicked = utils.PickRandomCharacter(charsLeft)
			pool.HintString = utils.GetHintString(word, charPicked, pool.HintString)

			utils.Cp("purple", pool.ID, "hintstring:", pool.HintString)
			pool.sendExcludingClientId(pool.CurrSketcher.ID, model.SocketMessage{
				Type:    89,
				Content: pool.HintString,
			})

			pool.NumHintsRevealed += 1
			utils.Cp("purple", pool.ID, "numHintsRevealed:", pool.NumHintsRevealed)
		}

		utils.Cp("purple", pool.ID, "exited hints reveal loop")
	}()
}

// 9, flag and broadcast game end
func (pool *Pool) endGame() {
	utils.Cp("greenBg", pool.ID, "all players done playing")

	pool.HasGameEnded = true
	pool.broadcast(model.SocketMessage{
		Type:    9,
		Content: pool.getClientInfoList(true).Content,
	})
}

func (pool *Pool) getClientForSketching() *Client {
	for _, c := range pool.Clients {
		if !c.DoneSketching {
			utils.Cp("yellow", pool.ID, "returning client:", c.Name)
			return c
		}
	}

	utils.Cp("red", pool.ID, "returning <nil> client")
	return nil
}

func (pool *Pool) allSketched() bool {
	flag := false
	for _, c := range pool.Clients {
		if c.DoneSketching {
			flag = true
		} else {
			flag = false
		}
	}

	utils.Cp("red", pool.ID, "allSketched:", flag)
	return flag
}
