package messages

import (
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
)

// a wrapper to communicate the state of a game ball in the game with specified id
type BallStateMessage struct {
	Ball   states.BallState `json:"Ball"`
	GameID string           `json:"GameID"` // the game that this update is taking place in
}
