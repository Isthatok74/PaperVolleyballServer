package requests

import (
	"github.com/Isthatok74/PaperVolleyballServer/data/states"
)

// represents a game ball, along with all of its state variables
type AddPlayerRequest struct {
	ClientPlayerID int               `json:"ClientPlayerID"`
	ServerPlayerID string            `json:"ServerPlayerID"`
	PlayerVars     states.PlayerVars `json:"PlayerVars"`
}
