package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"pv-server/data/requests"
	"pv-server/data/states"
	"pv-server/data/structures"
	"pv-server/defs"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// This is the Websocket implementation for the server

// the module that upgrades the http connection into a websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow any connections to this endpoint regardless of what it is
}

// connect a client via websocket, and register them to the client map
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
	s.sendws(conn, websocket.TextMessage, []byte(verifMsg))

	// start reading from this client's connection
	go s.readerws(conn)
}

// close the websocket connection
func (s *ServerData) closews(conn *websocket.Conn) {
	conn.Close()
	s.Clients.CompareAndDelete(conn.RemoteAddr().String(), conn)
	log.Printf("[%s] Websocket listener stopped", conn.RemoteAddr())
}

// listener for messages received from websocket connections
func (s *ServerData) readerws(conn *websocket.Conn) {

	// define a panic handling function
	defer func() {

		// handle panic
		if r := recover(); r != nil {
			log.Printf("[%s] Panic during websocket listener: %v", conn.RemoteAddr(), r)
		}

		// close the connection
		s.closews(conn)
	}()

	// handle disconnection error
	logMessageErr := func(err error) {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Printf("[%s] Unexpected close error: %v", conn.RemoteAddr(), err)
		} else if errors.Is(err, io.EOF) {
			log.Printf("[%s] Connection closed by client", conn.RemoteAddr())
		} else if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			log.Printf("[%s] Read timeout: %v", conn.RemoteAddr(), err)
		} else {
			log.Printf("[%s] Error reading message: %v", conn.RemoteAddr(), err)
		}
	}

	// setup a timeout check on this connection
	timeLastMsgReceived := time.Now()

	// continuously listen on the connection
	for {

		// handle timeout timer
		if time.Since(timeLastMsgReceived).Minutes() > defs.TimeoutMinutesWS {
			log.Printf("[%s] Timeout due to no requests received after a long time", conn.RemoteAddr())
			break
		}

		// receive a message when it arrives
		msgType, msgBody, err := conn.ReadMessage()
		if err != nil {
			logMessageErr(err)
			break
		}

		// add to the number of requests that have been processed
		if len(msgBody) > 0 {
			s.Info.CountRequests()
			s.Info.CountBytesReceived(uint64(overheadreceivews(msgBody) + len(msgBody)))
			timeLastMsgReceived = time.Now()
		}

		// parse it
		msg, err := parsews(msgType, msgBody)
		if err != nil {
			log.Printf("Unable to parse a message {%s} from %s: %v", msg, conn.RemoteAddr(), err)
			continue
		}
		log.Printf("[From %s] %s", conn.RemoteAddr(), msg)

		// process it
		res, err := s.processws(msg)
		if err != nil {
			log.Printf("Unable to process a message {%s} from %s: %v", msg, conn.RemoteAddr(), err)
			continue
		}

		// send a result message
		s.sendws(conn, websocket.TextMessage, res)
	}
}

// send a message to the specified websocket connection
func (s *ServerData) sendws(conn *websocket.Conn, messageType int, msgBody []byte) {
	s.Info.CountBytesSent(uint64(overheadsendws(msgBody) + len(msgBody)))
	if err := conn.WriteMessage(messageType, msgBody); err != nil {
		log.Println(err)
		return
	}
}

// returns the number of bytes of overhead bandwidth used to send a message via websockets
func overheadsendws(msgBody []byte) int {
	payloadSize := len(msgBody)
	if payloadSize <= 125 {
		return 2 // Header for small payloads
	} else if payloadSize <= 65535 {
		return 4 // Header for medium payloads
	} else {
		return 10 // Header for large payloads
	}
}

// returns the number of bytes of overhead bandwidth used to receive a message via websockets
func overheadreceivews(msgBody []byte) int {
	return overheadsendws(msgBody) + 4
}

// parse a message received from a websocket connection
func parsews(msgType int, msgBody []byte) ([]byte, error) {
	switch msgType {
	case websocket.TextMessage:
		return msgBody, nil
	//case websocket.BinaryMessage: // todo: implement?
	default:
		return nil, fmt.Errorf("unsupported message type: %d", msgType)
	}
}

// process an message containing information about an in-game event, and returns a message to send back
func (s *ServerData) processws(msgBody []byte) ([]byte, error) {

	// deserialize the message
	var data map[string]interface{}
	err := json.Unmarshal(msgBody, &data)
	if err != nil {
		fmt.Println("Error parsing incoming message: ", err)
		return []byte{}, err
	}

	// search for the "type" key-value pair to determine what type of data was pased in
	const jsonTagType string = "type"
	typeVal := ""
	for key := range data {
		val := data[key].(string)
		if strings.Contains(strings.ToLower(key), jsonTagType) {
			typeVal = strings.ToLower(val)
		}
	}
	if len(typeVal) == 0 {
		return []byte{}, fmt.Errorf("error finding type key in json string; unidentifiable message")
	}

	// read the wrapped data
	if strings.Contains(typeVal, JsonTagPingRequest) {

		// handle ping request
		return handlepingws(msgBody)

	} else if strings.Contains(typeVal, JsonTagCreateRequest) {

		// create game request
		return s.handlecreatews()
	}

	// attempt to find the gameID
	gameVal := ""
	for key := range data {
		val := data[key].(string)
		if strings.Contains(strings.ToLower(key), JsonTagGame) {
			gameVal = val
		}
	}
	if len(gameVal) == 0 {
		return []byte{}, fmt.Errorf("error finding game identifier key in json string; unidentifiable message")
	}

	// verify that the game exists
	game, err := s.FindGame(gameVal)
	if err != nil {
		return []byte{}, fmt.Errorf("could not find game id in registry: %s", gameVal)
	}

	// figure out what kind of game status update the message contains
	if strings.Contains(typeVal, JsonTagAddPlayerRequest) {

		// add player request
		return handleaddplayerws(msgBody, game)

	} else if strings.Contains(typeVal, JsonTagPlayer) {

		// player update, just rebroadcast the same message but to all connected clients
		s.broadcastws(msgBody, game)

	} else if strings.Contains(typeVal, JsonTagBall) {

		// ball update, check whether it is a valid hit or something else happened to the ball already
		log.Println("Processing ball event")

	} else {
		return []byte{}, fmt.Errorf("unrecognized json tag in received data; unidentifiable message")
	}

	// check for any hard-syncing events that need to be broadcasted, e.g.
	// * ending a rally
	// * ending the game
	// if any of these events occur, it is important that all connected clients be notified and synced up with the current state of the game

	// send a verification message back to the client who delivered this message
	return []byte(fmt.Sprintf("Processed message: %s", msgBody)), nil
}

// process a ping request
func handlepingws(msgBody []byte) ([]byte, error) {

	// deserialize the message
	var rq requests.PingRequest
	structures.FromWrappedJSON(&rq, msgBody)

	// re-serialize the message
	return structures.ToWrappedJSON(rq, "")
}

// process a game creation request
func (s *ServerData) handlecreatews() ([]byte, error) {

	// create a game in the data
	game := *states.NewGameState()
	s.Games.LoadOrStore(game.GUID, game)

	// create message to send back, with the game ID
	rq := requests.CreateGameRequest{
		GameID: game.GUID,
	}
	msg, err := structures.ToWrappedJSON(rq, game.GUID)
	return msg, err
}

// process a player add request
func handleaddplayerws(msgBody []byte, game *states.GameState) ([]byte, error) {

	// deserialize the message
	var rq requests.AddPlayerRequest
	structures.FromWrappedJSON(&rq, msgBody)

	// decode the JSON body to get player information
	newPlayer := states.NewPlayerState()

	// add the player to the game
	game.UpdatePlayer(newPlayer)

	// respond with the updated game state
	retrq := requests.AddPlayerRequest{
		ClientPlayerID: rq.ClientPlayerID,
		ServerPlayerID: newPlayer.GUID,
	}
	return structures.ToWrappedJSON(retrq, game.GUID)
}

// send a broadcast message to all clients connected to the specified game
func (s *ServerData) broadcastws(msgBody []byte, game *states.GameState) {

	// for each player connected to the game, send the message to the corresponding client

	// todo: we would need to store the player's ip in the playerState, we would need to resolve the discrepancy between their http ip and their ws ip
	// consider migrating ping, create and add to ws

}
