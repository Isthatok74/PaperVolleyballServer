package data

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pv-server/data/states"
	"pv-server/util"
	"sync"

	"github.com/gorilla/websocket"
)

type ServerData struct {
	Info    states.ServerState
	Games   map[string]states.GameState
	Clients sync.Map
}

// Constructor function to initialize ServerData
func NewServerData() *ServerData {
	serverData := &ServerData{
		Info:  *states.NewServerState(),          // Initialize Info field with zero value
		Games: make(map[string]states.GameState), // Initialize Games as an empty map
	}
	return serverData
}

// FindGame searches for a game by its ID and returns the GameState if found, or an error if not.
func (s *ServerData) FindGame(id string) (*states.GameState, error) {
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
	currentTime := util.CurrentTimeUTC().Format("15:04:05.000")
	fmt.Fprintf(w, "%s", currentTime)
}

func (s *ServerData) HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server start time: %s \n", s.Info.StartTime)
	fmt.Fprintf(w, "Number of requests processed: %d \n", s.Info.ReqCount)
	fmt.Fprintf(w, "Number of active games: %d \n", len(s.Games))
	fmt.Fprintf(w, "Number of clients connected: %d\n", util.GetSyncMapSize(&(s.Clients)))
}

func (s *ServerData) HandleCreate(w http.ResponseWriter, r *http.Request) {
	gameState := *states.NewGameState()
	s.Games[gameState.ID] = gameState
	fmt.Fprintf(w, "Created game with ID: %s \n", gameState.ID)
}

// AddPlayerToGame is the handler that adds a player to an existing game
func (s *ServerData) HandleAddPlayer(w http.ResponseWriter, r *http.Request) {

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
	wm.HandlePost(s)

	// Process the message (for now, we just print it to the console)
	fmt.Printf("Received message: %+v\n", wm)

	// Respond with a simple success message
	response := map[string]string{
		"status": "Message received successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Websocket implementation
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow any connections to this endpoint regardless of what it is
}

func (s *ServerData) HandleWS(w http.ResponseWriter, r *http.Request) {

	// upgrade the connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Client Successfully Connected: %s", r.RemoteAddr)

	// store the connection to the map
	clientAddr := conn.RemoteAddr().String()
	s.Clients.LoadOrStore(clientAddr, conn)

	// send a verification message to the client
	verifMsg := fmt.Sprintf("Server registry of client %s successful!", clientAddr)
	sendws(conn, websocket.TextMessage, []byte(verifMsg))

	// start reading
	s.readerws(conn)
}
func (s *ServerData) readerws(conn *websocket.Conn) {
	for {

		// receive a message when it arrives
		msgType, msgBody, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if len(msgBody) > 0 {
			s.Info.CountRequests()
		}

		// parse it
		msg, err := parsews(msgType, msgBody)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Message received from %s: %s", conn.RemoteAddr(), msg)

		// process it
		res := processws(msg)

		// send a result message
		sendws(conn, websocket.TextMessage, []byte(res))
	}
}
func sendws(conn *websocket.Conn, messageType int, msgBody []byte) {
	if err := conn.WriteMessage(messageType, msgBody); err != nil {
		log.Println(err)
		return
	}
}
func parsews(msgType int, msgBody []byte) (string, error) {
	switch msgType {
	case websocket.TextMessage:
		return string(msgBody), nil
	//case websocket.BinaryMessage: // todo: implement
	default:
		return "", fmt.Errorf("unsupported message type: %d", msgType)
	}
}
func processws(msgBody string) string {
	return fmt.Sprintf("Processed message: %s", msgBody)
}
