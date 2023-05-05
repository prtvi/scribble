package model

// SocketMessage.Content: int
// 1 - Connected - send to all
// 2 - Disconnected - send to all
// 3 - Text message - send to all
// 4 - Canvas data as string - send to all
// 5 - Get all client info - send to all but after processing
// 6 - start game ack - send to all but after processing

type SocketMessage struct {
	Type       int    `json:"type"`
	Content    string `json:"content"`
	ClientId   string `json:"clientId,omitempty"`
	ClientName string `json:"clientName,omitempty"`
	PoolId     string `json:"poolId,omitempty"`
}
