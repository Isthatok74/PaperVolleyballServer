package states

import (
	"pv-server/data/structures"

	"github.com/google/uuid"
)

type BaseState struct {
	structures.Serializable
	ID string
}

func (b *BaseState) GetGUID() {
	b.ID = uuid.New().String()
}
