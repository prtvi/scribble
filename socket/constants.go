package socket

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// map of {poolId: pool}
var hub = map[string]*Pool{}

// to handle socket connections
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var (
	RenderClientsEvery              = time.Second * 5
	InterGameWaitDuration           = time.Second * 2
	TimeoutForChoosingWord          = time.Second * 20
	DeletePoolAfterGameEndsDuration = time.Minute * 10
	RemovePoolAfterGameNotStarted   = time.Minute * 20
	ScoreForCorrectGuess            = 50

	debug = false
)

// B=> broadcasting to everyone
// b=> broadcasting to some

// S B=> Cs
// S b=> Cs
// S => C
// C => S B=> Cs
// C => S b=> Cs
// C => S

var messageTypeMap = map[int]string{
	// S B=> Cs - green
	6:   "sc__client_info",
	9:   "sc__end_game",
	10:  "sc__get_this_map",
	31:  "sc__correct_guess",
	312: "sc__word_in_msg",
	313: "sc__cant_reveal_word",
	32:  "sc__reveal_word",
	51:  "sc__clear_canvas",
	70:  "sc__game_started",
	71:  "sc__round_num",

	// S b=> Cs - cyan
	35: "sc__choosing_word",
	82: "sc__turn_over",
	84: "sc__turn_over_all_guessed",
	87: "sc__sketcher_begin_drawing",
	88: "sc__sketcher_drawing",
	89: "sc__hint",

	// S => C - yellow
	69: "sc__game_cant_start",
	8:  "sc__word_assigned",
	33: "sc__choose_word",
	81: "sc__disable_sketching",
	83: "sc__disable_sketching_all_guessed",
	86: "sc__midgame_timer",

	// C => S B=> Cs - blue
	1: "csc__client_connect",
	2: "csc__client_disconnect",
	3: "csc__text_msg",

	// C => S b=> Cs - purple
	4:  "csc__canvas_data",
	41: "csc__undo_draw",
	5:  "csc__clear_canvas",

	// C => S - red
	7:  "cs__req_start_game",
	34: "cs__chosen_word",
}
