package routes

import "time"

var (
	GameStartDurationInSeconds = time.Duration(time.Second * 120)
	TimeForEachWordInSeconds   = time.Duration(time.Second * 75)
	RenderClientsEvery         = time.Duration(time.Second * 5)
	ScoreForCorrectGuess       = 25
	NumberOfRounds             = 3
)

// map of {poolId: pool}
var HUB = map[string]*Pool{}

var messageTypeMap = map[int]string{
	1:  "client_connect",    // server b=> clients
	2:  "client_disconnect", // server b=> clients
	3:  "text_msg",          // client b=> clients
	31: "correct_guess",     // server b=> clients
	32: "reveal_word",       // server b=> clients
	4:  "canvas_data",       // client b=> clients
	5:  "clear_canvas",      // client b=> clients
	6:  "client_info",       // server b=> clients --at regular intervals
	7:  "start_game",        // client  => server  --to start the game
	71: "round_num",         // server b=> clients
	8:  "word_assigned",     // server b=> clients
	9:  "end_game",          // server b=> clients
}
