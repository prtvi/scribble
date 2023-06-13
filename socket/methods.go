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
func (c *Client) Read() {
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

	// remove the color that was picked in getColorForClient func from the list, the first color was picked from the list
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

func (pool *Pool) flagAllClientsAsNotSketched() {
	for _, cl := range pool.Clients {
		cl.DoneSketching = false
	}
}

func (pool *Pool) getColorForClient() string {
	return pool.ColorList[0]
}

func (pool *Pool) wordChooseCountdown(words []string) {
	sleep(TimeoutForChoosingWord)

	if pool.CurrWord == "" {
		pool.CurrWord = utils.GetRandomWord(words)
		return
	}
}

func (pool *Pool) turnOver(c *Client) {
	c.DoneSketching = true
	pool.CurrWord = ""
	pool.CurrSketcher = nil
}

// 70
func (pool *Pool) startGameAndBroadcast() {
	// flag and broadcast the starting of the game
	pool.HasGameStarted = true
	pool.broadcast(model.SocketMessage{
		Type:    70,
		TypeStr: messageTypeMap[70],
		Success: true,
	})
}

func (pool *Pool) clientWordAssignmentFlow(client *Client) {
	// select the client
	pool.CurrSketcher = client

	// create a list of words for client to choose
	words := utils.Get3RandomWords(utils.WORDS)
	pool.broadcast3WordsList(words) // DOTO: dont broadcast to everyone

	// start a timeout for assigning word if not chosen by client
	go pool.wordChooseCountdown(words)

	// run an infinite loop until pool.CurrWord is initialised by sketcher client, initialised in pool.Start func
	for pool.CurrWord == "" {
	}

	// add the word expiry
	pool.CurrWordExpiresAt = time.Now().Add(TimeForEachWordInSeconds)
}

// methods called in Start or BeginGameFlow funcs

func (pool *Pool) BeginBroadcastClientInfo() {
	// to be run as a go routine
	// starts an infinite loop to broadcast client info after every regular interval
	utils.Cp("yellow", "Broadcasting client info start!")

	for {
		sleep(RenderClientsEvery)
		pool.broadcastClientInfoList()

		// stop broadcasting when game ends
		if pool.HasGameEnded || len(pool.Clients) == 0 {
			utils.Cp("yellow", "Stopped broadcasting client info")
			break
		}
	}
}

func (pool *Pool) StartGameCountdown() {
	// as soon as the first player/client joins, start this countdown to start the game, after this timeout, the game begin message will broadcast

	// sleep until its the game starting time
	sleep(pool.GameStartTime.Sub(time.Now()))

	// if the game has already started by the client using the button then exit the countdown
	if pool.HasGameStarted {
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

	// if the sketcher is the guesser, then the guesserClient will be nil, hence check if guesserClient is nil
	// check if the word matches with the current word and check if the guesserClient hasn't already guessed
	if guesserClient != nil &&
		!guesserClient.HasGuessed &&
		guessedLower == currWordLower {

		// increment score and flag as guessed
		guesserClient.HasGuessed = true
		guesserClient.Score += ScoreForCorrectGuess * int(utils.GetDiffBetweenTimesInSeconds(time.Now(), pool.CurrWordExpiresAt))

		// broadcast client info list to update score on UI immediately
		pool.broadcastClientInfoList()

		// if correct guess then modify the message
		message.Type = 31
		message.TypeStr = messageTypeMap[31]
		return message
	}

	// check if the text message contains the word, word exists in message
	if strings.Contains(guessedLower, currWordLower) {
		message.Type = 312
		message.TypeStr = messageTypeMap[312]
	}

	return message
}

// 9
func (pool *Pool) EndGame() {
	// flag and broadcast game end

	utils.Cp("yellowBg", "All players done playing!")

	pool.HasGameEnded = true
	pool.broadcast(model.SocketMessage{
		Type:    9,
		TypeStr: messageTypeMap[9],
		Content: pool.getClientInfoList().Content,
	})
}
