package structures

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

// wrap a serializable object in a WrappedMessage container and return it in json format
func ToWrappedJSON(b any, gameID string) ([]byte, error) {
	data, err := json.Marshal(b)
	if err != nil {
		log.Println(err.Error())
		return []byte{}, err
	}
	wm := NewWrappedMessage(reflect.TypeOf(b).String(), string(data), gameID)
	msg, err := json.Marshal(wm)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}

// from json string of a wrapped message, construct the appropriate object
func FromWrappedJSON(b any, jsonData []byte) error {

	// deserialize the message
	var wm WrappedMessage
	err := json.Unmarshal(jsonData, &wm)
	if err != nil {
		log.Println("Error parsing incoming message as wrapped message: ", err)
		return err
	}

	// deserialize the data
	err = json.Unmarshal([]byte(wm.Data), &b)
	if err != nil {
		log.Println("Error parsing incoming data in wrapped message: ", err)
		return err
	}

	// no errors encountered
	return nil
}

// Generalized test function for serialization/deserialization
func CompareSerializeDeserialize[T any](t *testing.T, original T, getField func(T) string) {
	// Serialize to JSON
	msg, err := ToWrappedJSON(original, "")
	if err != nil {
		t.Fatalf("Error serializing: %v", err)
	}

	// Deserialize from JSON
	var copy T
	FromWrappedJSON(&copy, msg)

	// Check if the field value matches
	if getField(original) != getField(copy) {
		t.Errorf("Serialize %T = %s; want %s", original, getField(copy), getField(original))
	}
}
