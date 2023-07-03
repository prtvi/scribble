package socket

import (
	model "scribble/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client.Color: string
// color string hash value without the #

type Client struct {
	ID, Name, Color                    string
	IsOwner, DoneSketching, HasGuessed bool
	Score                              int
	Conn                               *websocket.Conn
	Pool                               *Pool
	mu                                 sync.Mutex
}

type Pool struct {
	ID string

	Capacity           int
	DrawTime           time.Duration
	Rounds             int
	WordMode           string
	WordCount          int
	Hints              int
	CustomWords        []string
	UseCustomWordsOnly bool

	JoiningLink string
	CurrWord    string
	CurrRound   int

	Register, Unregister chan *Client
	Broadcast            chan model.SocketMessage
	Clients              []*Client
	CurrSketcher         *Client

	ColorList                                                    []string
	CreatedTime, GameStartTime, GameStartedAt, CurrWordExpiresAt time.Time
	HasGameStarted, HasGameEnded, HasClientInfoBroadcastStarted  bool
}
