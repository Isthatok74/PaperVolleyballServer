package states

import (
	"math/rand"
	"sync"
	"time"
)

// represents a game instance on the server, with all its associated data stored
type LobbyState struct {
	RegisteredInstance
	RoomCode   string `json:"RoomCode"`   // the room code that players can enter to join
	Background string `json:"Background"` // the string code for the background asset
}

// initialize a new gameState object
func NewLobbyState(lobbyMap *sync.Map) *LobbyState {
	l := &LobbyState{}
	l.GenerateGUID()
	success := l.generateRoomCode(lobbyMap)
	l.RegisteredInstance.UpdateTime()
	if success {
		return l
	} else {
		return nil
	}
}

// generate room code; returns whether a valid code was found
func (l *LobbyState) generateRoomCode(lobbyMap *sync.Map) bool {

	// generate random codes until successful or too many attempts hitting a duplicate
	// * note: this could be optimized to always find a valid code until all possible codes are filled
	//         it could be accomplished by writing a trie which tracks how many nodes are filled in each child; a "filled" node is one that has exhausted all child combinations;
	//         at each level of the tree, only the non-filled characters will be rolled to determine which branch to follow;
	//         it will fail if all possible combinations are used
	// * in practice, we can just ensure that the number of combinations greatly exceeds the a likely maximum number of lobbies (i.e. by a few orders of magnitude) and code a simple while loop that regenerates until a non-duplicate is found
	l.RoomCode = ""
	numTries := 0
	for numTries < maxNumRoomCodeTries {
		numTries++
		trialCode := randomConsonants(NumRoomCodeChars)
		_, roomAlreadyExists := lobbyMap.Load(trialCode)
		if !roomAlreadyExists {
			l.RoomCode = trialCode
			break
		}
	}
	return len(l.RoomCode) > 0
}

const maxNumRoomCodeTries = 10000
const NumRoomCodeChars = 4

// update the background of the lobby
func (l *LobbyState) UpdateBackground(bgName string) {
	l.Background = bgName
}

// create a game instance and migrate the current instance information to it
func (l *LobbyState) CreateGameInstance() *GameState {
	g := NewGameState()
	g.RegisteredInstance = *l.RegisteredInstance.Clone()
	g.RegisteredInstance.UpdateTime()
	return g
}

// generates 4 random consonents
func randomConsonants(length int) string {
	consonants := "BCDFGHJKLMNPQRSTVWXZ"
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = consonants[rand.Intn(len(consonants))]
	}
	return string(result)
}
