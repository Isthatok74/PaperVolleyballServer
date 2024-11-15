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

// Purpose: A container for all the data tracked by the server in real time

type ServerData struct {
	Info    states.ServerState // vitals
	Games   sync.Map           // a map of all ongoing games hosted on this server
	Clients sync.Map           // a map of all connected players hosted on this server
}

// constructor function to initialize ServerData
func NewServerData() *ServerData {
	serverData := &ServerData{
		Info: *states.NewServerState(), // Initialize Info field with zero value
	}
	return serverData
}

// searches for a game by its ID and returns the GameState if found, or nil along with an error if not.
func (s *ServerData) FindGame(id string) (*states.GameState, error) {

	// look up the game by ID in the map
	value, exists := s.Games.Load(id)
	if !exists {

		// if the game is not found, return nil and an error
		return nil, fmt.Errorf("game not found with ID %s", id)
	}

	// attempt to cast it to what it should be, in order to return the correct object
	game, ok := value.(states.GameState)
	if !ok {

		// if somehow the game isn't the correct type, return nil and an error
		return nil, fmt.Errorf("value is not of type *GameState")
	}

	// return a pointer to the found game and nil error
	return &game, nil
}

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

		// add to the number of requests that have been processed
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
	//case websocket.BinaryMessage: // todo: implement?
	default:
		return "", fmt.Errorf("unsupported message type: %d", msgType)
	}
}
func processws(msgBody string) string {

	// deserialize the message

	// if it's a player update, just rebroadcast the same message but to all connected clients

	// if it's a ball update, check whether it is a valid hit or something else happened to the ball already

	// check for any hard-syncing events that need to be broadcasted, e.g.
	// * ending a rally
	// * ending the game
	// if any of these events occur, it is important that all connected clients be notified and synced up with the current state of the game

	// send a verification message back to the client who delivered this message
	return fmt.Sprintf("Processed message: %s", msgBody)
}
