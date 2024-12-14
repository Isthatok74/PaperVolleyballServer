package messages

import (
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
)

// a wrapper to communicate the state of a player action in the game or lobby with specified id / room code
type PlayerActionMessage struct {
	Action         states.PlayerAction `json:"Action"`
	PlayerServerID string              `json:"PlayerServerID"` // the server's ID of the player
	GameID         string              `json:"GameID"`         // if non-empty, the game that this update is taking place in
	RoomCode       string              `json:"RoomCode"`       // if non-empty, the lobby that this update is taking place in
}
