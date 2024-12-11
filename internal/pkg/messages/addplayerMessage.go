package messages

// a message that facilitates adding a player with specified server id to the game with specified id
type AddPlayerGameMessage struct {
	ErrMsg         string `json:"ErrMsg"`
	ServerPlayerID string `json:"ServerPlayerID"`
	GameID         string `json:"GameID"`
}

// a message that facilitates adding a player with specified server id to the lobby with specified room code
type AddPlayerLobbyMessage struct {
	ErrMsg         string `json:"ErrMsg"`
	ServerPlayerID string `json:"ServerPlayerID"`
	RoomCode       string `json:"RoomCode"`
}
