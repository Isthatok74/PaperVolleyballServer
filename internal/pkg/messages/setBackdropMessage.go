package messages

// a message that communicates switching the backdrop in a lobby
type SetBackdropMessage struct {
	RoomCode     string `json:"RoomCode"`
	ResourceName string `json:"ResourceName"`
}
