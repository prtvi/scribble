package routes

import (
	"fmt"
	model "scribble/model"
	utils "scribble/utils"
	"strings"
	"time"
)

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
		// utils.Cp("yellow", "broadcasting message 6 - client info")

		msg := getClientInfoList(pool, 6)
		pool.BroadcastMsg(msg)
	}
}

func (pool *Pool) startGameAndBroadcast() {
	// start the game and broadcast the message
	pool.HasGameStarted = true
	pool.BroadcastMsg(model.SocketMessage{
		Type:    7,
		Content: "Game has started",
		Success: true,
	})
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
	pool.begin()
}

func (pool *Pool) StartGame() {
	// when the client requests to start the game instead of the countdown
	// start the game and broadcast the same
	pool.startGameAndBroadcast()
	utils.Cp("greenBg", "Game started! by client using btn")

	// start game flow
	pool.begin()
}

func (pool *Pool) begin() {
	// schedule timers for current word and current sketcher
	time.Sleep(time.Second * 1)
	fmt.Println("new word flow")

	for i, c := range pool.Clients {

		fmt.Println("iteration #", i+1)

		pool.CurrSketcher = c
		pool.CurrWord = utils.GetRandomWord()
		pool.CurrWordExpiresAt = time.Now().Add(time.Second * TimeForEachWordInSeconds)
		c.HasSketched = true

		for _, cl := range pool.Clients {
			cl.HasGuessed = false
		}

		pool.BroadcastMsg(model.SocketMessage{
			Type:              8,
			Content:           "new word for a client",
			CurrSketcherId:    pool.CurrSketcher.ID,
			CurrWord:          pool.CurrWord,
			CurrWordExpiresAt: pool.CurrWordExpiresAt,
		})

		st := pool.CurrWordExpiresAt.Sub(time.Now())
		fmt.Println("sleeping for ...", st)
		time.Sleep(st)
		fmt.Println("sleep over")
	}
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

//

func (pool *Pool) EndGame() {

}
