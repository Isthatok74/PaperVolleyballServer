package states

import (
	"pv-server/data/structures"
	"strings"
)

// represents a game ball, along with all of its state variables
type BallState struct {
	BaseState
	Pos          structures.Vector2 `json:"Pos"`
	Vel          structures.Vector2 `json:"Vel"`
	GravityScale float32            `json:"GravityScale"`
	TouchedBy    string             `json:"TouchedBy"`
	TouchCount   int                `json:"TouchCount"`
	LiveState    string             `json:"LiveState"`
	ServeState   string             `json:"ServeState"`
}

func (b *BallState) IsAlive() bool {
	return strings.Contains(strings.ToLower(b.LiveState), "alive")
}
