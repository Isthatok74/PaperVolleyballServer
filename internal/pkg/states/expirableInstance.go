package states

import (
	"sync"
	"time"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/defs"
)

// A simple component for tracking when an instance of an object was last updated, in order to handle timeouts
type ExpirableInstance struct {
	BaseState
	LastUpdate time.Time
	mu         sync.Mutex
}

// update the time of the last change in this instance
func (r *ExpirableInstance) UpdateTime() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.LastUpdate = time.Now()
}

// returns whether too much time has elapsed since the last game update
func (g *ExpirableInstance) IsTimeoutExpired() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	durSinceLastUpdate := time.Since(g.LastUpdate)
	return durSinceLastUpdate.Minutes() > defs.TimeoutGameMinutesWS
}
