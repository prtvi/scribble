package model

import "time"

type SocketMessage struct {
	Type              int       `json:"type"`
	TypeStr           string    `json:"typeStr"`
	Content           string    `json:"content,omitempty"`
	Success           bool      `json:"success,omitempty"`
	ClientId          string    `json:"clientId,omitempty"`
	ClientName        string    `json:"clientName,omitempty"`
	PoolId            string    `json:"poolId,omitempty"`
	CurrRound         int       `json:"currRound"`
	CurrSketcherId    string    `json:"currSketcherId,omitempty"`
	CurrWord          string    `json:"currWord,omitempty"`
	CurrWordExpiresAt time.Time `json:"currWordExpiresAt,omitempty"`
}

type ClientInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Score int    `json:"score"`
}
