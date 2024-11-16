package states

import (
	st "pv-server/data/structures"
)

// represents a game ball, along with all of its state variables
type BallState struct {
	BaseState
	Pos        st.Vector2 `json:"Pos"`
	Vel        st.Vector2 `json:"Vel"`
	TouchedBy  string     `json:"TouchedBy"`
	TouchCount int        `json:"TouchCount"`
}
