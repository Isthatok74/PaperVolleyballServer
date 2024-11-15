package data

// a container for a json message string, which also includes a specifier `Type` to communicate what type of data is stored, and the gameID that it is relevant to
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
