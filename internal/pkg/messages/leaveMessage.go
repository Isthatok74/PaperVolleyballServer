package messages

// a request sent by the client to leave the game, if it's nice enough to send it
type LeaveGameMessage struct {
	GameID         string `json:"GameID"`
	PlayerServerID string `json:"PlayerServerID"`
}

// a request sent by the client to leave the lobby, if it's nice enough to send one
type LeaveLobbyMessage struct {
	RoomCode       string `json:"RoomCode"`
	PlayerServerID string `json:"PlayerServerID"`
}
