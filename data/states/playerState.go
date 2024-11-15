package states

import (
	st "pv-server/data/structures"
)

// represents an ingame player, with all of its state variables
type PlayerState struct {
	BaseState
	Pos       st.Vector2 `json:"Pos"`
	Vel       st.Vector2 `json:"Vel"`
	FaceRight bool       `json:"FaceRight"`
	Anim      string     `json:"Anim"`
}

// create a new instance of PlayerState
func NewPlayerState() *PlayerState {
	newPlayer := PlayerState{}
	newPlayer.GetGUID()
	return &newPlayer
}
