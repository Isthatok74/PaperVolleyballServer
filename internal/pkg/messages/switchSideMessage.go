package messages

// a message from client that requests to switch sides on the court
type SwitchSideMessage struct {
	ServerPlayerID string `json:"ServerPlayerID"`
	RoomCode       string `json:"RoomCode"`
}
