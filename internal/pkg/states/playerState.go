package states

import (
	"net"
)

// a container to store clients' session info,
type PlayerState struct {
	PlayerAction               // ingame transient data
	PlayerAttributes           // ingame constant data
	ExpirableInstance          // for handling user timeouts
	addr              net.Addr // make private for to avoid sending the address out
	GameID            string   // the id of the game the user is connected to, if any
	RoomCode          string   // the room code of the lobby that the user is connected to, if any
}

// create a new client container for a user with speicified address
func NewPlayer(address net.Addr) *PlayerState {
	client := &PlayerState{
		addr:     address,
		GameID:   "",
		RoomCode: "",
	}
	client.ExpirableInstance.GenerateGUID()
	client.ExpirableInstance.UpdateTime()
	return client
}

// updpate the player's game state
func (r *PlayerState) UpdatePlayerState(p *PlayerAction) {
	r.PlayerAction = *p
	r.UpdateTime()
}

// update the player variables
func (r *PlayerState) UpdatePlayerAttributes(p *PlayerAttributes) {
	r.PlayerAttributes = *p
	r.UpdateTime()
}

// expose private address variable
func (r *PlayerState) SetAddress(a net.Addr) {
	r.addr = a
}
func (r *PlayerState) GetAddress() net.Addr {
	return r.addr
}
