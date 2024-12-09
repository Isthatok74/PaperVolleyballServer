package requests

// a request sent by the client to register a new game instance
// if successful, the response returned by the server will be the guid of the newly registered game
type CreateGameRequest struct {
	GameID string `json:"GameID"`
}

// a request sent by the client to register a new lobby instance
// if successful, the reponse returned by the server will be the guid of the newly registered lobby, and the room code for display
// otherwise, the server may either not respond or return an error message
type CreateLobbyRequest struct {
	ErrMsg   string `json:"ErrMsg"`
	RoomCode string `json:"RoomCode"`
}
