package messages

import (
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
)

// a message that the server sends to force the player to update their position on the scene
type ForcePlayerMessage struct {
	Action         states.PlayerAction `json:"Action"`
	ServerPlayerID string              `json:"ServerPlayerID"`
}
