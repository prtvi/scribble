package model

import "time"

// SocketMessage.Type: int
// Determines the content type in SocketMessage.Content field
// 1 - Connected             - send to all
// 2 - Disconnected          - send to all
// 3 - Text message          - send to all
// 4 - Canvas data as string - send to all
// 5 - Get all client info   - send to all but after processing
// 6 - start game ack        - send to all but after processing

// SocketMessage.ClientId & SocketMessage.ClientName: string
// The client that triggers conn/disconnection to server, Register/Unregister event at server

// SocketMessage.PoolId: string
// Used by client to request all client info list and start game

type SocketMessage struct {
	Type              int       `json:"type"`
	Content           string    `json:"content"`
	ClientId          string    `json:"clientId,omitempty"`
	ClientName        string    `json:"clientName,omitempty"`
	PoolId            string    `json:"poolId,omitempty"`
	CurrSketcherId    string    `json:"currSketcherId,omitempty"`
	CurrWord          string    `json:"currWord,omitempty"`
	CurrWordExpiresAt time.Time `json:"currWordExpiresAt,omitempty"`
}

type ClientInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
