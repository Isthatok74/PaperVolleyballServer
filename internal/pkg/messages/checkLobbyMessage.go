package messages

// a message that checks whether a lobby with the given room code is open
type CheckLobbyMessage struct {
	Exists   bool   `json:"Exists"`
	RoomCode string `json:"RoomCode"`
}
