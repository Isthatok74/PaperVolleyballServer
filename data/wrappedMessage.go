package data

type WrappedMessage struct {
	Type   string `json:"Type"`
	Data   string `json:"Data"`
	GameID string `json:"GameID"`
}

// define how posted json messages received http are handled here
func (w *WrappedMessage) HandlePost(s *ServerData) string {
	// currently there are no known messages that will be processed this way
	return w.Data
}
