package socket

import (
	model "scribble/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID            string
	Name          string
	AvatarConfig  model.AvatarConfig
	IsSketching   bool
	DoneSketching bool
	HasGuessed    bool
	Score         int
	Conn          *websocket.Conn
	Pool          *Pool
	mu            sync.Mutex
}

type Pool struct {
	ID                 string
	Capacity           int
	Rounds             int
	WordCount          int
	DrawTime           time.Duration
	Hints              int
	WordMode           string
	CustomWords        []string
	UseCustomWordsOnly bool

	Register   chan *Client
	Unregister chan *Client
	Clients    []*Client
	Broadcast  chan model.SocketMessage

	WordsForGame         []string
	InitCurrWord         chan string
	CurrWord             string
	HintString           string
	NumHintsForCurrWord  int
	NumHintsRevealed     int
	CurrRound            int
	CurrSketcher         *Client
	CurrWordExpiresAt    time.Time
	SleepingForSketching bool

	HasGameStarted                bool
	HasGameEnded                  bool
	HasClientInfoBroadcastStarted bool

	JoiningLink   string
	CreatedTime   time.Time
	GameStartedAt time.Time
}
