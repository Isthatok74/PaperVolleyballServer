package states

// represents a game instance on the server, with all its associated data stored
type LobbyState struct {
	BaseState
	RegisteredInstance
	RoomCode   string `json:"RoomCode"`
	Background string `json:"Background"`
	//mu       sync.Mutex // Mutex to protect concurrent access to member variables
}

// initialize a new gameState object
func NewLobbyState() *LobbyState {
	l := &LobbyState{}
	l.GetGUID()
	l.generateRoomCode()
	l.RegisteredInstance.UpdateTime()
	return l
}

// generate room code
func (l *LobbyState) generateRoomCode( /* take an input list of existing codes */ ) {
	// todo: generate random code with random 4 consonents,
	l.RoomCode = "XXFP"
}

// update the background of the lobby
func (l *LobbyState) UpdateBackground(bgName string) {
	l.Background = bgName
}

// update the team of a player in the lobby
func (l *LobbyState) UpdateTeam(pguid string, rightSide bool) {
	// todo: switch the player's team in the data
}

// create a game instance and migrate the current instance information to it
func (l *LobbyState) CreateGameInstance() *GameState {
	g := NewGameState()
	g.RegisteredInstance = *l.RegisteredInstance.Clone()
	g.RegisteredInstance.UpdateTime()
	return g
}
