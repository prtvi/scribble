package routes

import (
	model "scribble/model"
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
}

type Pool struct {
	ID                                                          string
	Capacity, ColorAssignmentIndex                              int
	Register, Unregister                                        chan *Client
	Clients                                                     []*Client
	Broadcast                                                   chan model.SocketMessage
	CreatedTime, GameStartTime, CurrWordExpiresAt               time.Time
	HasGameStarted, HasGameEnded, HasClientInfoBroadcastStarted bool
	CurrSketcher                                                *Client
	CurrWord                                                    string
}
