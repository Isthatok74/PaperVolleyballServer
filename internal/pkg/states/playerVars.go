package states

import (
	"net"
)

// represents the constant variables of a player that is connected to the game
type PlayerVars struct {
	BaseState
	DisplayName string   `json:"DisplayName"`
	Strength    float32  `json:"Strength"`
	Speed       float32  `json:"Speed"`
	Jump        float32  `json:"Jump"`
	Size        float32  `json:"Size"`
	Tier        int      `json:"Tier"`
	addr        net.Addr // make private to avoid sending the address out
}

func (v *PlayerVars) SetAddress(a net.Addr) {
	v.addr = a
}
func (v *PlayerVars) GetAddress() net.Addr {
	return v.addr
}
