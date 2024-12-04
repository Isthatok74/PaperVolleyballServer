package states

import (
	"pv-server/data/structures"
	"strings"
)

// represents a game ball, along with all of its state variables
type BallState struct {
	BaseState
	Pos          structures.Vector2 `json:"Pos"`          // the ball's position
	Vel          structures.Vector2 `json:"Vel"`          // the ball's velocity
	GravityScale float32            `json:"GravityScale"` // the ball's gravity factor
	TouchedBy    string             `json:"TouchedBy"`    // the id of the player who last touched the ball
	TouchCount   int                `json:"TouchCount"`   // the number of touches made on the ball
	LiveState    string             `json:"LiveState"`    // the live/dead status code of the ball
	ServeState   string             `json:"ServeState"`   // the service status code of the ball
}

// returns whether the ball's `LiveState“ indicates that it is alive
func (b *BallState) IsAlive() bool {
	return strings.Contains(strings.ToLower(b.LiveState), "alive")
}
