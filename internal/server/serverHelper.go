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

	// return a pointer to the found object and nil error
	return game, nil
}

// searches for a lobby by its room code and returns the LobbyState if found, or nil along with an error if not.
func (s *ServerData) FindLobby(roomCode string) (*states.LobbyState, error) {

	// look up the game by ID in the map
	value, exists := s.Lobbies.Load(roomCode)
	if !exists {

		// if the game is not found, return nil and an error
		return nil, fmt.Errorf("lobby not found with room code: %s", roomCode)
	}

	// attempt to cast it to what it should be, in order to return the correct object
	lobby, ok := value.(*states.LobbyState)
	if !ok {

		// if somehow the game isn't the correct type, return nil and an error
		return nil, fmt.Errorf("value is not of type *LobbyState")
	}

	// return a pointer to the found object and nil error
	return lobby, nil
}

// searches for a player by its ID and returns the PlayerState if found, or nil along with an error if not.
func (s *ServerData) FindPlayer(id string) (*states.PlayerState, error) {

	// look up the ID in the map
	value, exists := s.Players.Load(id)
	if !exists {

		// if the player is not found, return nil and an error
		return nil, fmt.Errorf("player not found with ID %s", id)
	}

	// attempt to cast it to what it should be, in order to return the correct object
	player, ok := value.(*states.PlayerState)
	if !ok {

		// if somehow the game isn't the correct type, return nil and an error
		return nil, fmt.Errorf("value is not of type *PlayerState")
	}

	// return a pointer to the found object and nil error
	return player, nil
}
