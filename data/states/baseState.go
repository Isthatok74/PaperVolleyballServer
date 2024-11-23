package states

import (
	"github.com/google/uuid"
)

// a base class for state containers. all derivates should have an ID string
type BaseState struct {
	GUID string `json:"GUID"`
}

// a method that should be called to generate a guid for this object
func (b *BaseState) GetGUID() {
	b.GUID = uuid.New().String()
}
