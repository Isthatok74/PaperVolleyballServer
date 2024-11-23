package states

import (
	"github.com/google/uuid"
)

// a base class for state containers. all derivates should have an ID string
type BaseState struct {
	ID string `json:"ID"`
}

// a method that should be called to generate a guid for this object
func (b *BaseState) GetGUID() {
	b.ID = uuid.New().String()
}
