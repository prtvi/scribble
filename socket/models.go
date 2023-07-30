package socket

import (
	model "scribble/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID, Name                  string
	AvatarConfig              model.AvatarConfig
	DoneSketching, HasGuessed bool
	Score                     int
	Conn                      *websocket.Conn
	Pool                      *Pool
	mu                        sync.Mutex
}

type Pool struct {
	ID                 string
	Capacity           int
	Rounds             int
	WordCount          int
	DrawTime           time.Duration
	Hints              int
	HintsForCurrWord   int
	HintsRevealed      int
	WordMode           string
	CustomWords        []string
	UseCustomWordsOnly bool

	Register   chan *Client
	Unregister chan *Client
	Clients    []*Client
	Broadcast  chan model.SocketMessage

	CurrWord                      string
	CurrRound                     int
	CurrSketcher                  *Client
	CurrWordExpiresAt             time.Time
	HasGameStarted                bool
	HasGameEnded                  bool
	HasClientInfoBroadcastStarted bool

	JoiningLink   string
	CreatedTime   time.Time
	GameStartedAt time.Time
}
