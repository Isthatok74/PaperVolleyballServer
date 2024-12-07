package states

import (
	"pv-server/defs"
	"sync"
	"time"
)

// represents a game instance on the server, with all its associated data stored
type GameState struct {
	BaseState
	Ball       *BallState `json:"Ball"`
	Players    sync.Map   `json:"Players"`    // key: string; value: PlayerState
	PlayerInfo sync.Map   `json:"PlayerInfo"` // key: string; value: PlayerVars
	lastUpdate time.Time
	mu         sync.Mutex // Mutex to protect concurrent access to Ball and Players
}

// initialize a new gameState object
func NewGameState() *GameState {
	gameState := &GameState{
		Ball:       nil, // no ball exists yet
		lastUpdate: time.Now(),
	}
	gameState.GetGUID()
	return gameState
}

// updpate the player game state on the map
func (g *GameState) UpdatePlayerState(p *PlayerState) {
	g.Players.LoadOrStore(p.GUID, *p)
	g.lastUpdate = time.Now()
}

// update the player variables on the map
func (g *GameState) UpdatePlayerVars(p *PlayerVars) {
	g.PlayerInfo.LoadOrStore(p.GUID, *p)
	g.lastUpdate = time.Now()
}

// update the ball data on the map
func (g *GameState) UpdateBall(b *BallState) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Ball = b
	g.lastUpdate = time.Now()
}

// returns whether too much time has elapsed since the last game update
func (g *GameState) IsTimeoutExpired() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	durSinceLastUpdate := time.Since(g.lastUpdate)
	return durSinceLastUpdate.Minutes() > defs.TimeoutGameMinutesWS
}
