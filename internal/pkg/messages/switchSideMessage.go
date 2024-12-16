package messages

type SwitchSideMessage struct {
	ServerPlayerID string `json:"ServerPlayerID"`
	RoomCode       string `json:"RoomCode"`
}
