package structures

import (
	"encoding/json"
)

type Serializable struct {
}

func (b *Serializable) FromJSON(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), b)
}
