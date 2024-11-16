package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pv-server/data/states"
	"pv-server/util"
)

// Handles the http traffic portion of the server

// handle the ping route on http - simply returns the current server time
func (s *ServerData) HandlePing(w http.ResponseWriter, r *http.Request) {
	currentTime := util.CurrentTimeUTC().Format("15:04:05.000")
	fmt.Fprintf(w, "%s", currentTime)
}

// handle the status route on http - returns some server metrics
func (s *ServerData) HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server start time: %s \n", s.Info.StartTime)
	fmt.Fprintf(w, "Number of requests processed: %d \n", s.Info.ReqCount)
	fmt.Fprintf(w, "Number of active games: %d \n", util.GetSyncMapSize(&(s.Games)))
	fmt.Fprintf(w, "Number of clients connected: %d\n", util.GetSyncMapSize(&(s.Clients)))
}

// handle the creation of a game instance on the server
func (s *ServerData) HandleCreate(w http.ResponseWriter, r *http.Request) {
	gameState := *states.NewGameState()
	s.Games.LoadOrStore(gameState.ID, gameState)
	fmt.Fprintf(w, "Created game with ID: %s \n", gameState.ID)
}

// handle registering player to an existing game
func (s *ServerData) HandleAddPlayer(w http.ResponseWriter, r *http.Request) {

	// get the game GUID from the URL path (assuming it's a part of the URL like /games/{guid}/addplayer)
	gameID := r.URL.Query().Get("gameID") // Example URL: /addplayer?gameID={guid}

	// check if the game exists
	game, err := s.FindGame(gameID)
	if err != nil {
		http.Error(w, "Failed to add player to game", http.StatusInternalServerError)
	}

	// decode the JSON body to get player information
	newPlayer := states.NewPlayerState()

	// add the player to the game
	game.UpdatePlayer(newPlayer)

	// respond with the updated game state
	fmt.Fprintf(w, "In game %s added player: %s\n", game.ID, newPlayer.ID)
}

// Handle POST request for the /post endpoint
func (s *ServerData) HandlePost(w http.ResponseWriter, r *http.Request) {

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

	// Process the message (for now, we just print it to the console)
	fmt.Printf("Message received: %s\n", wm.Data)
	result := wm.HandlePost(s)

	// Respond with a simple success message
	response := map[string]string{
		"status": result,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
