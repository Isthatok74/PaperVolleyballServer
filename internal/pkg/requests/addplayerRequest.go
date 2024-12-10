package requests

// a request from the client to add a player with specified server id to the game with specified id
type AddPlayerGameRequest struct {
	ErrMsg         string `json:"ErrMsg"`
	ServerPlayerID string `json:"ServerPlayerID"`
	GameID         string `json:"GameID"`
}

// a rqeuest from the client to add a player with specified server id to the lobby with specified room code
type AddPlayerLobbyRequest struct {
	ErrMsg         string `json:"ErrMsg"`
	ServerPlayerID string `json:"ServerPlayerID"`
	RoomCode       string `json:"RoomCode"`
}
