package server

import (
	"fmt"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
)

// Contains helper functions related to managing server data and requests

// searches for a game by its ID and returns the GameState if found, or nil along with an error if not.
func (s *ServerData) FindGame(id string) (*states.GameState, error) {

	// look up the game by ID in the map
	value, exists := s.Games.Load(id)
	if !exists {

		// if the game is not found, return nil and an error
		return nil, fmt.Errorf("game not found with ID %s", id)
	}

	// attempt to cast it to what it should be, in order to return the correct object
	game, ok := value.(*states.GameState)
	if !ok {

		// if somehow the game isn't the correct type, return nil and an error
		return nil, fmt.Errorf("value is not of type *GameState")
	}

	// return a pointer to the found game and nil error
	return game, nil
}
