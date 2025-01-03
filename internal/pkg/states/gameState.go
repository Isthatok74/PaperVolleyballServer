package states

import (
	"sync"
)

// represents a game instance on the server, with all its associated data stored
type GameState struct {
	RegisteredInstance
	Ball *BallState `json:"Ball"`
	mu   sync.Mutex // Mutex to protect concurrent access to Ball
}

// initialize a new gameState object
func NewGameState() *GameState {
	gameState := &GameState{
		Ball: nil, // no ball exists yet
	}
	gameState.GenerateGUID()
	gameState.RegisteredInstance.UpdateTime()
	return gameState
}

// update the ball data on the map
func (g *GameState) UpdateBall(b *BallState) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Ball = b
	g.RegisteredInstance.UpdateTime()
}

// return a copy of the game ball's data for threadsafe operations
func (g *GameState) GetBallCopy() *BallState {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.Ball == nil {
		return nil
	}
	return g.Ball.Clone()
}
