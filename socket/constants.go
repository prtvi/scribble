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
	TimeoutForChoosingWord          = time.Second * 15
	DeletePoolAfterGameEndsDuration = time.Minute * 10
	RemovePoolAfterGameNotStarted   = time.Minute * 20
	ScoreForCorrectGuess            = 50

	debug = false
)

// B=> broadcasting to everyone
// b=> broadcasting to some
// Cs - all clients
// C - single client

// S B=> Cs        - server to all clients
// S b=> Cs        - server to some clients
// S => C          - server to one client
// C => S B=> Cs   - one client to server - broadcasting to all clients
// C => S b=> Cs   - one clients to server - broadcasting to some clients
// C => S          - one client to server

var class1 = "server_b_clients"
var class2 = "server_b_some_clients"
var class3 = "server_client"
var class4 = "client_server_b_clients"
var class5 = "client_server_b_some_clients"
var class6 = "client_server"

var messageTypeMap = map[int]string{
	// S B=> Cs - green
	6:   class1 + "__client_info",
	9:   class1 + "__end_game",
	31:  class1 + "__correct_guess",
	312: class1 + "__word_in_msg",
	313: class1 + "__cant_reveal_word",
	32:  class1 + "__reveal_word",
	51:  class1 + "__clear_canvas",
	70:  class1 + "__game_started",
	71:  class1 + "__round_num",

	// S b=> Cs - cyan
	35: class2 + "__choosing_word",
	82: class2 + "__turn_over",
	84: class2 + "__turn_over_all_guessed",
	87: class2 + "__sketcher_begin_drawing",
	88: class2 + "__sketcher_drawing",
	89: class2 + "__hint",

	// S => C - yellow
	10: class3 + "__shared_config",
	69: class3 + "__game_cant_start",
	8:  class3 + "__word_assigned",
	33: class3 + "__choose_word",
	81: class3 + "__disable_sketching",
	83: class3 + "__disable_sketching_all_guessed",
	86: class3 + "__midgame_timer",

	// C => S B=> Cs - blue
	1: class4 + "__client_connect",
	2: class4 + "__client_disconnect",
	3: class4 + "__text_msg",

	// C => S b=> Cs - purple
	4:  class5 + "__canvas_data",
	41: class5 + "__undo_draw",
	5:  class5 + "__clear_canvas",

	// C => S - red
	7:  class6 + "__req_start_game",
	34: class6 + "__chosen_word",
}
