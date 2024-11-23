package structures

// a container for a json message string, which also includes a specifier `Type` to communicate what type of data is stored, and the gameID that it is relevant to
type WrappedMessage struct {
	Type   string `json:"Type"`
	Data   string `json:"Data"`
	GameID string `json:"GameID"`
}

// constructor
func NewWrappedMessage(t string, data string, gameid string) WrappedMessage {
	return WrappedMessage{
		Type:   t,
		Data:   string(data),
		GameID: gameid,
	}
}
