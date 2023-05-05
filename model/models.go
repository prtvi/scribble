package model

// SocketMessage.Content: int
// 1 - Connected
// 2 - Disconnected
// 3 - Text message
// 4 - Canvas data as string
// 5 - Get all client info
// 6 - start game ack

type SocketMessage struct {
	Type       int    `json:"type"`
	Content    string `json:"content"`
	ClientId   string `json:"clientId,omitempty"`
	ClientName string `json:"clientName,omitempty"`
	PoolId     string `json:"poolId,omitempty"`
}
