package routes

const (
	GameStartDurationInSeconds = 1000
	TimeForEachWordInSeconds   = 45 // 75 finally
	ScoreForCorrectGuess       = 25
	RenderClientsEvery         = 8
	NumberOfRounds             = 3
)

// map of {poolId: pool}
var HUB = map[string]*Pool{}

var messageTypeMap = map[int]string{
	1:  "client_connect",    // server b=> clients
	2:  "client_disconnect", // server b=> clients
	3:  "text_msg",          // client b=> clients
	4:  "canvas_data",       // client b=> clients
	5:  "clear_canvas",      // client b=> clients
	6:  "client_info",       // server b=> clients --at regular intervals
	7:  "start_game",        // client  => server  --to start the game
	71: "round_num",         // server b=> clients
	8:  "word_assigned",     // server b=> clients
	9:  "end_game",          // server b=> clients
}
