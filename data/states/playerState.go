package states

import (
	st "github.com/Isthatok74/PaperVolleyballServer/data/structures"
)

// represents an ingame player, with all of its state variables
type PlayerState struct {
	BaseState
	Pos       st.Vector2 `json:"Pos"`       // the player's position
	Vel       st.Vector2 `json:"Vel"`       // the player's velocity
	FaceRight bool       `json:"FaceRight"` // the player's facing direction
	Anim      string     `json:"Anim"`      // the player's primary animation code
	AxisX     float32    `json:"AxisX"`     // the player's x axis control value
}
