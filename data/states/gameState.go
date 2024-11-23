package states

// represents a game instance on the server, with all its associated data stored
type GameState struct {
	BaseState
	Ball    *BallState             `json:"Ball"`
	Players map[string]PlayerState `json:"Players"`
}

// initialize a new gameState object
func NewGameState() *GameState {
	gameState := &GameState{
		Ball:    nil,                          // no ball exists yet
		Players: make(map[string]PlayerState), // create an empty map of players
	}
	gameState.GetGUID()
	return gameState
}

// updpate the player data on the map
func (g *GameState) UpdatePlayer(p *PlayerState) {
	g.Players[p.ID] = *p
}

// update the ball data on the map
func (g *GameState) UpdateBall(b *BallState) {
	g.Ball = b
}
