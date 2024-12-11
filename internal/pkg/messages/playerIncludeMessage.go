package messages

import (
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
)

// a message that contains information about a single player in the game or lobby, used to instruct a client to add them in if not done already
type PlayerIncludeMessage struct {
	Attributes states.PlayerAttributes `json:"Attributes"`
	Action     states.PlayerAction     `json:"Action"`
}
