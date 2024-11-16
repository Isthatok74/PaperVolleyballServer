package data

import (
	"encoding/json"
	"log"
	"pv-server/data/structures"
	"reflect"
)

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

// wrap a data container into a message which also contains its type
func MakeWrappedMessage(obj *structures.Serializable, gameID string) (WrappedMessage, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println(err.Error())
		return WrappedMessage{}, err
	}
	return NewWrappedMessage(reflect.TypeOf(obj).String(), string(data), gameID), nil
}

// define how posted json messages received http are handled here
func (w *WrappedMessage) HandlePost(s *ServerData) string {

	// currently there are no known messages that will be processed this way
	return w.Data
}
