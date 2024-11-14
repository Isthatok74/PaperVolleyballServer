package states

import (
	st "pv-server/data/structures"
)

type PlayerState struct {
	BaseState
	Pos       st.Vector2 `json:"Pos"`
	Vel       st.Vector2 `json:"Vel"`
	FaceRight bool       `json:"FaceRight"`
	Anim      string     `json:"Anim"`
}

func NewPlayerState() *PlayerState {
	newPlayer := PlayerState{}
	newPlayer.GetGUID()
	return &newPlayer
}
