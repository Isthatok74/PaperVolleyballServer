package states

import (
	"sync"
	"time"
)

// Base struct for instances containing live players
type RegisteredInstance struct {
	Players    sync.Map `json:"Players"`    // key: string; value: PlayerState
	PlayerInfo sync.Map `json:"PlayerInfo"` // key: string; value: PlayerVars
	LastUpdate time.Time
	mu         sync.Mutex
}

// updpate the player game state on the map
func (r *RegisteredInstance) UpdatePlayerState(p *PlayerState) {
	r.Players.LoadOrStore(p.GUID, *p)
	r.UpdateTime()
}

// update the player variables on the map
func (r *RegisteredInstance) UpdatePlayerVars(p *PlayerVars) {
	r.PlayerInfo.LoadOrStore(p.GUID, *p)
	r.UpdateTime()
}

func (r *RegisteredInstance) UpdateTime() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.LastUpdate = time.Now()
}
