package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pv-server/data/states"
	st "pv-server/data/states"
	"pv-server/util"
)

type ServerData struct {
	Info  st.ServerState
	Games map[string]st.GameState
}

// Constructor function to initialize ServerData
func NewServerData() *ServerData {
	serverData := &ServerData{
		Info:  *st.NewServerState(),          // Initialize Info field with zero value
		Games: make(map[string]st.GameState), // Initialize Games as an empty map
	}
	return serverData
}

// FindGame searches for a game by its ID and returns the GameState if found, or an error if not.
func (s *ServerData) FindGame(id string) (*st.GameState, error) {
	// Look up the game by ID in the map
	game, exists := s.Games[id]
	if !exists {
		// If the game is not found, return nil and an error
		return nil, fmt.Errorf("game not found with ID %s", id)
	}
	// Return the found game and nil error
	return &game, nil
}

func (s *ServerData) HandlePing(w http.ResponseWriter, r *http.Request) {
	s.Info.CountRequests()
	currentTime := util.CurrentTimeUTC().Format("15:04:05.000")
	fmt.Fprintf(w, "%s", currentTime)
}

func (s *ServerData) HandleStatus(w http.ResponseWriter, r *http.Request) {
	s.Info.CountRequests()
	fmt.Fprintf(w, "Server start time: %s \n", s.Info.StartTime)
	fmt.Fprintf(w, "Number of requests processed: %d \n", s.Info.ReqCount)
	fmt.Fprintf(w, "Number of active games: %d \n", len(s.Games))
}

func (s *ServerData) HandleCreate(w http.ResponseWriter, r *http.Request) {
	s.Info.CountRequests()
	gameState := *st.NewGameState()
	s.Games[gameState.ID] = gameState
	fmt.Fprintf(w, "Created game with ID: %s \n", gameState.ID)
}

// AddPlayerToGame is the handler that adds a player to an existing game
func (s *ServerData) HandleAddPlayer(w http.ResponseWriter, r *http.Request) {

	s.Info.CountRequests()

	// Get the game GUID from the URL path (assuming it's a part of the URL like /games/{guid}/addplayer)
	gameID := r.URL.Query().Get("gameID") // Example URL: /addplayer?gameID={guid}

	// Check if the game exists
	game, exists := s.Games[gameID]
	if !exists {
		http.Error(w, "Game not found", http.StatusNotFound)
		return
	}

	// Decode the JSON body to get player information
	newPlayer := states.NewPlayerState()

	// Add the player to the game
	game.Players[newPlayer.ID] = *newPlayer

	// Respond with the updated game state
	fmt.Fprintf(w, "In game %s added player: %s\n", game.ID, newPlayer.ID)
}

// Handle POST request for the /message endpoint
func (s *ServerData) HandlePostMessage(w http.ResponseWriter, r *http.Request) {

	s.Info.CountRequests()

	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON body into a Message struct
	var wm WrappedMessage
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&wm)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	wm.HandleWrappedMessage(s)

	// Process the message (for now, we just print it to the console)
	fmt.Printf("Received message: %+v\n", wm)

	// Respond with a simple success message
	response := map[string]string{
		"status": "Message received successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
