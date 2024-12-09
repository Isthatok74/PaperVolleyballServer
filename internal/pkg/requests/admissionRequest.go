package requests

import (
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
)

// for initializing a client's data on the server
type AdmissionRequest struct {
	ClientPlayerID int                     `json:"ClientPlayerID"`
	ServerPlayerID string                  `json:"ServerPlayerID"`
	Attributes     states.PlayerAttributes `json:"Attributes"`
}
