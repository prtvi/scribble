package model

// Type:              int identifier
// Content:           can contain any kind of data, type of data identified using Type field
// Success:           to acknowledge success for any event
// ClientId:          incoming clientId from clients
// ClientName:        incoming clientName from clients
// PoolId:            associated poolId
// CurrRound:         current round
// CurrSketcherId:    current sketcher's ID
// CurrSketcherName:  current sketcher's name
// CurrWord:          current word to be guessed
// CurrWordExpiresAt: current word's expiry time as a string
// TimeoutAfter:      timeout for any event as such, rn used for choosing word

type SocketMessage struct {
	Type              int    `json:"type"`
	Content           string `json:"content,omitempty"`
	Success           bool   `json:"success,omitempty"`
	MidGameJoinee     bool   `json:"midGameJoinee,omitempty"`
	WordMode          string `json:"wordMode,omitempty"`
	ClientId          string `json:"clientId,omitempty"`
	ClientName        string `json:"clientName,omitempty"`
	PoolId            string `json:"poolId,omitempty"`
	CurrRound         int    `json:"currRound,omitempty"`
	CurrSketcherId    string `json:"currSketcherId,omitempty"`
	CurrSketcherName  string `json:"currSketcherName,omitempty"`
	CurrWord          string `json:"currWord,omitempty"`
	CurrWordLen       int    `json:"currWordLen,omitempty"`
	CurrWordExpiresAt string `json:"currWordExpiresAt,omitempty"`
	TimeoutAfter      string `json:"timeoutAfter,omitempty"`
}

type ClientInfo struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Score        int          `json:"score"`
	IsSketching  bool         `json:"isSketching"`
	HasGuessed   bool         `json:"hasGuessed"`
	AvatarConfig AvatarConfig `json:"avatarConfig"`
}
type AvatarConfig struct {
	Color     Coords `json:"color"`
	Eyes      Coords `json:"eyes"`
	Mouth     Coords `json:"mouth"`
	IsOwner   bool   `json:"isOwner"`
	IsCrowned bool   `json:"isCrowned"`
}
type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SharedConfig struct {
	MessageTypeMap                 map[int]string `json:"messageTypeMap"`
	TimeForEachWordInSeconds       int            `json:"timeForEachWordInSeconds"`
	TimeForChoosingWordInSeconds   int            `json:"timeForChoosingWordInSeconds"`
	CloseSocketConnectionInSeconds int            `json:"closeSocketConnInSeconds"`
	PrintLogs                      bool           `json:"printLogs"`
}

type FormOption struct {
	Value, Label string
	Selected     bool
}
type CreateFormParam struct {
	ID, Label, Desc string
	ImgIdx          int
	Options         []FormOption
}

type PoolStat struct {
	ID               string `json:"id"`
	NumActiveClients int    `json:"numActiveClients"`
	HasGameStarted   bool   `json:"hasGameStarted"`
	HasGameEnded     bool   `json:"hasGameEnded"`
	CurrSketcher     string `json:"currSketcher"`
	CreatedTime      string `json:"createdTime"`
	GameStartedAt    string `json:"gameStartedAt"`
}
type Stats struct {
	LenHub        int        `json:"lenHub"`
	NumGoroutines int        `json:"numGoroutines"`
	Pools         []PoolStat `json:"pools"`
}

type ApiResp struct {
	Message string `json:"message"`
}
