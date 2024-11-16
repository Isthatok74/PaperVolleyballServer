package data

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// This is the Websocket implementation for the server

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
		res, err := processws(msg)
		if err != nil {
			log.Println(err)
			return
		}

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
func parsews(msgType int, msgBody []byte) ([]byte, error) {
	switch msgType {
	case websocket.TextMessage:
		return msgBody, nil
	//case websocket.BinaryMessage: // todo: implement?
	default:
		return nil, fmt.Errorf("unsupported message type: %d", msgType)
	}
}
func processws(msgBody []byte) (string, error) {

	/*
		// deserialize the message
		var data map[string]interface{}
		err := json.Unmarshal(msgBody, &data)
		if err != nil {
			fmt.Println("Error parsing incoming message: ", err)
			return "", err
		}

		// search for the "type" and "game" key-value pairs to determine what type of data was pased in, which game it corresponds to
		const jsonTagType string = "type"
		typeVal := ""
		gameVal := ""
		for key := range data {
			val := data[key].(string)
			if strings.Contains(strings.ToLower(key), jsonTagType) {
				typeVal = val
			} else if strings.Contains(strings.ToLower(key), states.JsonTagGame) {
				gameVal = val
			}
		}
		if len(typeVal) == 0 {
			return "", fmt.Errorf("error finding type key in json string; unidentifiable message")
		}
		if len(gameVal) == 0 {
			return "", fmt.Errorf("error finding game identifier key in json string; unidentifiable message")
		}

		// read the wrapped data
		if strings.Contains(strings.ToLower(typeVal), states.JsonTagPlayer) {

			// player update, just rebroadcast the same message but to all connected clients

		} else if strings.Contains(strings.ToLower(typeVal), states.JsonTagBall) {

			// ball update, check whether it is a valid hit or something else happened to the ball already

		} else {
			return "", fmt.Errorf("unrecognized json tag in received data; unidentifiable message")
		}
	*/

	// check for any hard-syncing events that need to be broadcasted, e.g.
	// * ending a rally
	// * ending the game
	// if any of these events occur, it is important that all connected clients be notified and synced up with the current state of the game

	// send a verification message back to the client who delivered this message
	return fmt.Sprintf("Processed message: %s", msgBody), nil
}
