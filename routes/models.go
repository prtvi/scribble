package routes

import (
	model "scribble/model"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client.Color: string
// color string hash value without the #

type Client struct {
	ID, Name, Color         string
	HasSketched, HasGuessed bool
	Score                   int
	Conn                    *websocket.Conn
	Pool                    *Pool
	mu                      sync.Mutex
}

type Pool struct {
	ID, JoiningLink, CurrWord                                   string
	Capacity, CurrRound                                         int
	ColorList                                                   []string
	Register, Unregister                                        chan *Client
	Clients                                                     []*Client
	Broadcast                                                   chan model.SocketMessage
	CreatedTime, GameStartTime, CurrWordExpiresAt               time.Time
	HasGameStarted, HasGameEnded, HasClientInfoBroadcastStarted bool
	CurrSketcher                                                *Client
}
