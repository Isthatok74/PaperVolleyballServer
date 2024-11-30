package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"pv-server/data/requests"
	"pv-server/data/states"
	"pv-server/data/structures"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// This file contains the engine logic for processing game requests

// process an message containing information about an in-game event, and returns a message to send back
func (s *ServerData) processws(conn *websocket.Conn, msgBody []byte) ([]byte, error) {

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
		return handleping(msgBody)

	} else if strings.Contains(typeVal, JsonTagCreateRequest) {

		// create game request
		return s.handlecreate()
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
		return handleaddplayer(conn.RemoteAddr(), msgBody, game)

	} else if strings.Contains(typeVal, JsonTagPlayer) {

		// player update, just rebroadcast the same message but to all connected clients
		s.broadcastws(msgBody, game)

	} else if strings.Contains(typeVal, JsonTagBall) {

		// ball update, check whether it is a valid hit or something else happened to the ball already
		return s.handleballevent(msgBody, game)

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
func handleping(msgBody []byte) ([]byte, error) {

	// deserialize the message
	var rq requests.PingRequest
	structures.FromWrappedJSON(&rq, msgBody)

	// re-serialize the message
	return structures.ToWrappedJSON(rq, "")
}

// process a game creation request
func (s *ServerData) handlecreate() ([]byte, error) {

	// create a game in the data
	game := *states.NewGameState()
	s.Games.LoadOrStore(game.GUID, &game)

	// start a routine that times the game out if too much time has passed since it last updated
	checkGameTimeout := func(g *states.GameState) {
		for g != nil {
			if g.IsTimeoutExpired() {
				log.Printf("Deleting game %s due to timeout", g.GUID)
				s.Games.CompareAndDelete(g.GUID, g)
				break
			}
			time.Sleep(time.Minute) // sleep for some time to prevent high CPU usage and avoid tight looping
		}
	}
	go checkGameTimeout(&game)

	// create message to send back, with the game ID
	rq := requests.CreateGameRequest{
		GameID: game.GUID,
	}
	msg, err := structures.ToWrappedJSON(rq, game.GUID)
	return msg, err
}

// process a player add request
func handleaddplayer(addr net.Addr, msgBody []byte, game *states.GameState) ([]byte, error) {

	// deserialize the message
	var rq requests.AddPlayerRequest
	structures.FromWrappedJSON(&rq, msgBody)

	// decode the JSON body to get player information
	newPlayerVars := rq.PlayerVars
	newPlayerVars.GetGUID()
	newPlayerVars.SetAddress(addr)
	newPlayer := states.PlayerState{}
	newPlayer.GUID = newPlayerVars.GUID

	// add the player to the game
	game.UpdatePlayerState(&newPlayer)
	game.UpdatePlayerVars(&newPlayerVars)

	// respond with the updated game state
	retrq := requests.AddPlayerRequest{
		ClientPlayerID: rq.ClientPlayerID,
		ServerPlayerID: newPlayerVars.GUID,
	}
	return structures.ToWrappedJSON(retrq, game.GUID)
}

func (s *ServerData) handleballevent(msgBody []byte, game *states.GameState) ([]byte, error) {

	// deserialize the message
	var clientBall states.BallState
	structures.FromWrappedJSON(&clientBall, msgBody)

	// if for whatever reason the client's copy of the ball is out of date (e.g. someone else has registered a hit before them or the ball has already died), do not process the request and return a harmless error to the client
	denyBallUpdate := func() ([]byte, error) {

		// todo: revise output
		msgDrop := "ball update request denied"
		log.Println(msgDrop)
		return []byte(""), errors.New(msgDrop)
	}

	// accept the ball update and broadcast it
	acceptBallUpdate := func(b *states.BallState) ([]byte, error) {
		sendMsg, err := structures.ToWrappedJSON(b, game.GUID)
		if err == nil {
			s.broadcastws(sendMsg, game)
		} else {
			log.Printf("Error broadcasting accepted ball status: %s", err)
		}
		return []byte(""), err
	}

	// begin processing
	if len(clientBall.GUID) == 0 {

		// handle new ball registry
		clientBall.GetGUID()

		// register it to the game
		if game.Ball == nil {
			game.Ball = &clientBall
			game.Ball.GetGUID()
			log.Printf("Logged new game ball on server : %s", game.Ball.GUID)
			return acceptBallUpdate(game.Ball)
		} else {
			return denyBallUpdate()
		}

	} else {

		// check if ball id matches the one that is live on the server
		matchesLiveID := game.Ball != nil && game.Ball.GUID == clientBall.GUID && game.Ball.IsAlive()
		if !matchesLiveID {
			return denyBallUpdate()
		}

		// check whether the client's ball update indicates that the ball is alive
		isAlive := clientBall.IsAlive()
		if isAlive {

			// if the ball is still alive, it means the player touched it; check if the touch count makes sense
			isTouchCountCorrect := clientBall.TouchCount <= 1 || (clientBall.TouchCount-game.Ball.TouchCount == 1)
			if !isTouchCountCorrect {
				return denyBallUpdate()
			}

			// broadcast the updated client ball to other players
			game.UpdateBall(&clientBall)
			return acceptBallUpdate(game.Ball)

		} else {

			// if game ball was alive but client says it's dead, broadcast the dead ball and kill the ball on game side
			game.UpdateBall(nil)
			return acceptBallUpdate(&clientBall)
		}
	}
}
