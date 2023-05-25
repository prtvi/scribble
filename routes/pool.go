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
	later := now.Add(time.Second * GameStartDurationInSeconds)

	return &Pool{
		ID:                            uuid,
		Capacity:                      capacity,
		Register:                      make(chan *Client),
		Unregister:                    make(chan *Client),
		Clients:                       make([]*Client, 0),
		Broadcast:                     make(chan model.SocketMessage),
		ColorList:                     utils.COLORS[:10],
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

func (pool *Pool) Start() {
	// start listening to pool connections and messages
	for {
		select {
		case client := <-pool.Register:
			// on client register, append the client to Pool.Client slice
			pool.appendClientToList(client)

			pool.BroadcastMsg(model.SocketMessage{
				Type:       1,
				TypeStr:    messageTypeMap[1],
				ClientId:   client.ID,
				ClientName: client.Name,
			})

			// start broadcasting client info list
			if len(pool.Clients) == 1 &&
				!pool.HasClientInfoBroadcastStarted &&
				!pool.HasGameStarted {

				pool.HasClientInfoBroadcastStarted = true
				utils.Cp("yellowBg", "Broadcasting client info start!")

				// begin braodcasting client info at regular intervals
				go pool.BroadcastClientInfoMessage()

				// begin game start countdown
				go pool.StartGameCountdown()
			}

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client connected:"), client.Name)

			break

		case client := <-pool.Unregister:
			// on client disconnect, delete the client from Pool.Client slice
			pool.removeClientFromList(client)

			pool.BroadcastMsg(model.SocketMessage{
				Type:       2,
				TypeStr:    messageTypeMap[2],
				ClientId:   client.ID,
				ClientName: client.Name,
			})

			utils.Cp("yellow", "Size of connection pool:", utils.Cs("reset", fmt.Sprintf("%d", len(pool.Clients))), utils.Cs("yellow", "client disconnected:"), client.Name)

			break

		case message := <-pool.Broadcast:
			// on message received from any of the clients in the pool, broadcast the message to all clients
			// any of the game logic there is will be applied when clients do something, which will happen after the message is received from any of the clients

			utils.Cp("blue", "sm recv, type:", utils.Cs("yellow", fmt.Sprintf("%d:", message.Type)), utils.Cs("reset", messageTypeMap[message.Type], utils.Cs("blue", "from:"), message.ClientName))

			switch message.Type {
			case 3:
				pool.UpdateScore(message)
				pool.BroadcastMsg(message)

			case 4, 5:
				pool.BroadcastMsg(message)

			case 7:
				pool.StartGame()

			default:
				break
			}
		}
	}
}

func (pool *Pool) startGameAndBroadcast() {
	// start the game and broadcast the message
	pool.HasGameStarted = true
	pool.BroadcastMsg(model.SocketMessage{
		Type:    7,
		TypeStr: messageTypeMap[7],
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

func (pool *Pool) GetColorForClient() string {
	return pool.ColorList[0]
}

func (pool *Pool) BroadcastMsg(message model.SocketMessage) {
	// broadcasts the given message to all clients in the pool
	for _, c := range pool.Clients {
		c.Conn.WriteJSON(message)
	}
}

func (pool *Pool) BroadcastClientInfoMessage() {
	// starts a timer to broadcast client info after every regular interval
	for {
		time.Sleep(time.Second * RenderClientsEvery)
		utils.Cp("yellow", "Broadcasting client info")

		msg := pool.getClientInfoList()
		pool.BroadcastMsg(msg)
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
	utils.Cp("greenBg", "Game started! by server using countdown")

	// start game flow
	go pool.BeginGameFlow()
}

func (pool *Pool) StartGame() {
	// when the client requests to start the game instead of the countdown
	// start the game and broadcast the same
	pool.startGameAndBroadcast()
	utils.Cp("greenBg", "Game started! by client using btn")

	// start game flow
	go pool.BeginGameFlow()
}

func (pool *Pool) BeginGameFlow() {
	// schedule timers for current word and current sketcher
	d := time.Duration(time.Second * 2)
	utils.Cp("green", "Starting game in", d.String())
	time.Sleep(d)

	for _, c := range pool.Clients {
		pool.CurrSketcher = c
		pool.CurrWord = utils.GetRandomWord()
		pool.CurrWordExpiresAt = time.Now().Add(time.Second * TimeForEachWordInSeconds)
		c.HasSketched = true

		for _, cl := range pool.Clients {
			cl.HasGuessed = false
		}

		pool.BroadcastMsg(model.SocketMessage{
			Type:              8,
			TypeStr:           messageTypeMap[8],
			CurrSketcherId:    pool.CurrSketcher.ID,
			CurrWord:          pool.CurrWord,
			CurrWordExpiresAt: pool.CurrWordExpiresAt,
		})

		st := pool.CurrWordExpiresAt.Sub(time.Now())
		utils.Cp("yellow", "Sleeping for ...", st.String(), c.Name)
		time.Sleep(st)

		pool.BroadcastMsg(model.SocketMessage{
			Type:    5,
			TypeStr: messageTypeMap[5],
		})
	}

	pool.EndGame()
}

func (pool *Pool) UpdateScore(message model.SocketMessage) {
	// update score for the client that guesses the word right

	// when the game has not begun, the curr sketcher will be nil
	if pool.CurrSketcher == nil {
		return
	}

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
		// increment score and flag as guessed
		guesserClient.Score += ScoreForCorrectGuess * int(utils.GetDiffBetweenTimesInSeconds(time.Now(), pool.CurrWordExpiresAt))
		guesserClient.HasGuessed = true
	}
}

func (pool *Pool) EndGame() {
	utils.Cp("yellowBg", "All players done playing!")
	pool.HasGameEnded = true

	pool.BroadcastMsg(model.SocketMessage{
		Type:    9,
		TypeStr: messageTypeMap[9],
		Content: pool.getClientInfoList().Content,
	})
}
