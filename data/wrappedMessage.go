package data

type WrappedMessage struct {
	Type   string `json:"Type"`
	Data   string `json:"Data"`
	GameID string `json:"GameID"`
}

func (w *WrappedMessage) HandleWrappedMessage(s *ServerData) {

	// todo: construct the appropriate state object based on `Type`

	// todo: based on what type of message it is, direct it to the appropriate function

	// todo: update game state as necessary

	// todo: broadcast a message back to the client
}
