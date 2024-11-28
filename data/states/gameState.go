package states

import (
	"pv-server/defs"
	"sync"
	"time"
)

// represents a game instance on the server, with all its associated data stored
type GameState struct {
	BaseState
	Ball       *BallState             `json:"Ball"`
	Players    map[string]PlayerState `json:"Players"`
	PlayerInfo map[string]PlayerVars  `json:"PlayerInfo"`
	lastUpdate time.Time
	mu         sync.Mutex // Mutex to protect concurrent access to Ball and Players
}

// initialize a new gameState object
func NewGameState() *GameState {
	gameState := &GameState{
		Ball:       nil,                          // no ball exists yet
		Players:    make(map[string]PlayerState), // create an empty map of players
		lastUpdate: time.Now(),
	}
	gameState.GetGUID()
	return gameState
}

// updpate the player game state on the map
func (g *GameState) UpdatePlayerState(p *PlayerState) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Players[p.GUID] = *p
	g.lastUpdate = time.Now()
}

// update the player variables on the map
func (g *GameState) UpdatePlayerVars(guid string, p *PlayerVars) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.PlayerInfo[guid] = *p
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
