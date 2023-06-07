package socket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	GameStartDurationInSeconds = time.Duration(time.Second * 120)
	TimeForEachWordInSeconds   = time.Duration(time.Second * 75)
	RenderClientsEvery         = time.Duration(time.Second * 5)
	WaitAfterRoundStarts       = time.Duration(time.Second * 2)
	WaitAfterTurnEnds          = time.Duration(time.Second * 2)
	TimeoutForChoosingWord     = time.Duration(time.Second * 10)
	ScoreForCorrectGuess       = 25
	NumberOfRounds             = 3
)

// map of {poolId: pool}
var HUB = map[string]*Pool{}

var messageTypeMap = map[int]string{
	// server b=> clients
	1:  "sc___client_connect",
	2:  "sc___client_disconnect",
	6:  "sc___client_info",
	8:  "sc___word_assigned",
	9:  "sc___end_game",
	10: "sc___get_this_map",
	31: "sc___correct_guess",
	32: "sc___reveal_word",
	33: "sc___choose_word",
	70: "sc___game_started",
	71: "sc___round_num",
	81: "sc___turn_over",

	// client => server b=> clients
	3: "csc___text_msg",
	4: "csc___canvas_data",
	5: "csc___clear_canvas",

	// client => server
	7:  "cs___req_start_game",
	34: "cs___chosen_word",
}