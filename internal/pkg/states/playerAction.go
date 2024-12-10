package states

import (
	st "github.com/Isthatok74/PaperVolleyballServer/internal/pkg/structures"
)

// represents an ingame player, with all of its state variables
type PlayerAction struct {
	Pos       st.Vector2 `json:"Pos"`       // the player's position
	Vel       st.Vector2 `json:"Vel"`       // the player's velocity
	FaceRight bool       `json:"FaceRight"` // the player's facing direction
	Anim      string     `json:"Anim"`      // the player's primary animation code
	AxisX     float32    `json:"AxisX"`     // the player's x axis control value
	GameID    string     `json:"GameID"`    // if non-empty, the game that this update is taking place in
	RoomCode  string     `json:"RoomCode"`  // if non-empty, the lobby that this update is taking place in
}
