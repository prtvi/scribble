package socket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// map of {poolId: pool}
var HUB = map[string]*Pool{}

// to handle socket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	GameStartDurationInSeconds      = time.Second * 120
	TimeForEachWordInSeconds        = time.Second * 75
	RenderClientsEvery              = time.Second * 5
	WaitAfterRoundStarts            = time.Second * 2
	WaitAfterTurnEnds               = time.Second * 2
	TimeoutForChoosingWord          = time.Second * 10
	DeletePoolAfterGameEndsDuration = time.Minute * 10
	RemovePoolAfterGameNotStarted   = time.Minute * 20
	ScoreForCorrectGuess            = 25
	NumberOfRounds                  = 3

	DEBUG = false
)

var messageTypeMap = map[int]string{
	// server b=> clients - green
	6:   "sc__client_info",
	8:   "sc__word_assigned",
	9:   "sc__end_game",
	10:  "sc__get_this_map",
	31:  "sc__correct_guess",
	312: "sc__word_in_msg",
	32:  "sc__reveal_word",
	33:  "sc__choose_word",
	51:  "sc__clear_canvas",
	70:  "sc__game_started",
	71:  "sc__round_num",
	81:  "sc__turn_over",

	// client => server b=> clients - blue
	1: "csc__client_connect",
	2: "csc__client_disconnect",
	3: "csc__text_msg",
	4: "csc__canvas_data",
	5: "csc__clear_canvas",

	// client => server - purple
	7:  "cs__req_start_game",
	34: "cs__chosen_word",
}
