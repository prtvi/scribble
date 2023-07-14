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
	ID string

	Capacity           int
	DrawTime           time.Duration
	Rounds             int
	WordMode           string
	WordCount          int
	Hints              int
	HintsRevealed      int
	CustomWords        []string
	UseCustomWordsOnly bool

	JoiningLink string
	CurrWord    string
	CurrRound   int

	Register, Unregister chan *Client
	Broadcast            chan model.SocketMessage
	Clients              []*Client
	CurrSketcher         *Client

	CreatedTime, GameStartedAt, CurrWordExpiresAt               time.Time
	HasGameStarted, HasGameEnded, HasClientInfoBroadcastStarted bool
}
