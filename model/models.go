package model

// SocketMessage.Content: int
// 1 - Connected
// 2 - Disconnected
// 3 - Text message
// 4 - Canvas data as string
// 5 - start game ack

type SocketMessage struct {
	Type       int    `json:"type"`
	Content    string `json:"content"`
	ClientId   string `json:"clientId,omitempty"`
	ClientName string `json:"clientName,omitempty"`
}
