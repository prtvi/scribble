package socket

import (
	"net/http"
	model "scribble/model"
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
	GameStartDurationInSeconds      = time.Second * 120
	TimeForEachWordInSeconds        = time.Second * 75
	RenderClientsEvery              = time.Second * 5
	InterGameWaitDuration           = time.Second * 2
	TimeoutForChoosingWord          = time.Second * 15
	DeletePoolAfterGameEndsDuration = time.Minute * 10
	RemovePoolAfterGameNotStarted   = time.Minute * 20
	ScoreForCorrectGuess            = 25
	NumberOfRounds                  = 3

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

	// S => C - yellow
	8:  "sc__word_assigned",
	33: "sc__choose_word",
	81: "sc__disable_sketching",
	83: "sc__disable_sketching_all_guessed",

	// C => S B=> Cs - blue
	1: "csc__client_connect",
	2: "csc__client_disconnect",
	3: "csc__text_msg",

	// C => S b=> Cs - red
	4: "csc__canvas_data",
	5: "csc__clear_canvas",

	// C => S - purple
	7:  "cs__req_start_game",
	34: "cs__chosen_word",
}

var FormParams = []model.CreateFormParam{
	{ID: "players", Label: "Players", ImgIdx: 1,
		Options: []model.FormOption{
			{Value: "2", Label: "2"},
			{Value: "3", Label: "3"},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5", Selected: true},
			{Value: "6", Label: "6"},
			{Value: "7", Label: "7"},
			{Value: "8", Label: "8"},
			{Value: "9", Label: "9"},
			{Value: "10", Label: "10"}}},

	{ID: "drawTime", Label: "Draw time", ImgIdx: 2,
		Options: []model.FormOption{
			{Value: "15", Label: "15"},
			{Value: "20", Label: "20"},
			{Value: "40", Label: "40"},
			{Value: "50", Label: "50"},
			{Value: "60", Label: "60"},
			{Value: "70", Label: "70"},
			{Value: "80", Label: "80", Selected: true},
			{Value: "90", Label: "90"},
			{Value: "100", Label: "100"},
			{Value: "120", Label: "120"},
			{Value: "150", Label: "150"},
			{Value: "180", Label: "180"},
			{Value: "210", Label: "210"},
			{Value: "240", Label: "240"}}},

	{ID: "rounds", Label: "Rounds", ImgIdx: 3,
		Options: []model.FormOption{
			{Value: "2", Label: "2"},
			{Value: "3", Label: "3", Selected: true},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5"},
			{Value: "6", Label: "6"},
			{Value: "7", Label: "7"},
			{Value: "8", Label: "8"},
			{Value: "9", Label: "9"},
			{Value: "10", Label: "10"}}},

	{ID: "wordMode", Label: "Word mode", ImgIdx: 4,
		Options: []model.FormOption{
			{Value: "normal", Label: "Normal", Selected: true},
			{Value: "hidden", Label: "Hidden"},
			{Value: "combination", Label: "Combination"}}},

	{ID: "wordCount", Label: "Word count", ImgIdx: 5,
		Options: []model.FormOption{
			{Value: "1", Label: "1"},
			{Value: "2", Label: "2"},
			{Value: "3", Label: "3", Selected: true},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5"}}},

	{ID: "hints", Label: "Hints", ImgIdx: 6,
		Options: []model.FormOption{
			{Value: "1", Label: "1"},
			{Value: "2", Label: "2", Selected: true},
			{Value: "3", Label: "3"},
			{Value: "4", Label: "4"},
			{Value: "5", Label: "5"}}},
}
