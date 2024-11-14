package states

type GameState struct {
	BaseState
	Ball    *BallState             `json:"Ball"`
	Players map[string]PlayerState `json:"Players"`
}

// Initialize a new gameState object
func NewGameState() *GameState {
	gameState := &GameState{
		Ball:    nil,                          // no ball exists yet
		Players: make(map[string]PlayerState), // create an empty map of players
	}
	gameState.GetGUID()
	return gameState
}

func (g *GameState) UpdatePlayer(p *PlayerState) {
	g.Players[p.ID] = *p
}

func (g *GameState) UpdateBall(b *BallState) {
	g.Ball = b
}
