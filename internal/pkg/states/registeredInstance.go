package states

import (
	"sync"
	"time"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/defs"
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/util"
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

// update the time of the last change in this instance
func (r *RegisteredInstance) UpdateTime() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.LastUpdate = time.Now()
}

// create a clone of the stored instance
func (r *RegisteredInstance) Clone() *RegisteredInstance {
	return &RegisteredInstance{
		Players:    *util.CopySyncMap(&r.Players),
		PlayerInfo: *util.CopySyncMap(&r.PlayerInfo),
		LastUpdate: r.LastUpdate,
	}
}

// returns whether too much time has elapsed since the last game update
func (g *RegisteredInstance) IsTimeoutExpired() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	durSinceLastUpdate := time.Since(g.LastUpdate)
	return durSinceLastUpdate.Minutes() > defs.TimeoutGameMinutesWS
}
