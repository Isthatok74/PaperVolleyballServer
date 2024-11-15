package structures

import (
	"encoding/json"
)

// a base class for objects that are serializable to json
type Serializable struct {
}

// convert a json string to an instance of this object
func (b *Serializable) FromJSON(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), b)
}
