package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/messages"
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/structures"

	"github.com/gorilla/websocket"
)

// This file contains the engine logic for processing game requests
// * It defines all handlers of messages received from the client, and serves as a function directory for the differents types of messages that can be received
// * Helper functions are contained in a separate file

// process an message containing information about an in-game event, and returns a message to send back
func (s *ServerData) processws(conn *websocket.Conn, msgBody []byte) ([]byte, error) {

	// deserialize the message
	var data map[string]interface{}
	err := json.Unmarshal(msgBody, &data)
	if err != nil {
		fmt.Println("Error parsing incoming message: ", err)
		return nil, err
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
		return nil, fmt.Errorf("error finding type key in json string; unidentifiable message")
	}

	// read the wrapped data and direct to the processing function
	if strings.Contains(typeVal, JsonTagPingMsg) {

		// handle ping request
		return handleping(msgBody)

	} else if strings.Contains(typeVal, JsonTagCreateGameMsg) {

		// create game request
		return s.handlecreategame()

	} else if strings.Contains(typeVal, JsonTagCreateLobbyMsg) {

		// create lobby request
		return s.handlecreatelobby()

	} else if strings.Contains(typeVal, JsonTagAdmissionMsg) {

		// register a client to the server
		return s.handleadmitplayer(conn.RemoteAddr(), msgBody)

	} else if strings.Contains(typeVal, JsonTagAddPlayerMsg) {

		// add player to game request
		return s.handleaddplayergame(conn, msgBody)

	} else if strings.Contains(typeVal, JsonTagAddPlayerLobby) {

		// add player to lobby request
		return s.handleaddplayerlobby(conn, msgBody)

	} else if strings.Contains(typeVal, JsonTagRemPlayerLobby) {

		// remove player from lobby request
		return s.handleleavelobby(msgBody)

	} else if strings.Contains(typeVal, JsonTagRemPlayerGame) {

		// remove player from game request
		return s.handleleavegame(msgBody)

	} else if strings.Contains(typeVal, JsonTagSetBackdrop) {

		// set the backdrop resource name
		return s.handlesetbackdrop(msgBody)

	} else if strings.Contains(typeVal, JsonTagCheckLobbyMsg) {

		// check if a room code exists
		return s.handlechecklobby(msgBody)

	} else if strings.Contains(typeVal, JsonTagSwitchMsg) {

		// switch player to other side request
		return s.handleswitch(msgBody)

	} else if strings.Contains(typeVal, JsonTagPlayerEvent) {

		// player update, just rebroadcast the same message but to all connected clients of the corresponding game
		return s.handleplayeraction(msgBody)

	} else if strings.Contains(typeVal, JsonTagBallEvent) {

		// ball update, check whether it is a valid hit or something else happened to the ball already
		return s.handleballevent(msgBody)

	} else {
		return nil, fmt.Errorf("unrecognized json tag in received data; unidentifiable message")
	}
}

// process a ping request
func handleping(msgBody []byte) ([]byte, error) {

	// deserialize the message
	var rq messages.PingMessage
	structures.FromWrappedJSON(&rq, msgBody)

	// re-serialize the message
	return structures.ToWrappedJSON(rq)
}

// process a game creation request
func (s *ServerData) handlecreategame() ([]byte, error) {

	// create a game in the data
	game := *states.NewGameState()
	s.Games.LoadOrStore(game.GUID, &game)

	// start a routine that times the game out if too much time has passed since it last updated
	checkTimeout := func(g *states.GameState) {
		for g != nil {
			if g.RegisteredInstance.IsTimeoutExpired() {
				log.Printf("Deleting game %s due to timeout", g.GUID)
				s.Games.CompareAndDelete(g.GUID, g)
				break
			}
			time.Sleep(time.Minute) // sleep for some time to prevent high CPU usage and avoid tight looping
		}
	}
	go checkTimeout(&game)

	// create message to send back, with the game ID
	rq := messages.CreateGameMessage{
		GameID: game.GUID,
	}
	msg, err := structures.ToWrappedJSON(rq)
	return msg, err
}

// process a lobby creation request
func (s *ServerData) handlecreatelobby() ([]byte, error) {

	// prepare message
	rq := messages.CreateLobbyMessage{}

	// create a lobby in the data
	lobby := states.NewLobbyState(&s.Lobbies)
	if lobby == nil {
		errMsg := "There are too many instances of player-hosted lobbies at the moment. Please try again later."
		rq.ErrMsg = errMsg
		log.Println(errMsg)
	} else {
		rq.RoomCode = lobby.RoomCode
		s.Lobbies.Store(lobby.RoomCode, lobby)
		log.Printf("Succesfully registered a lobby with room code {%s}", lobby.RoomCode)

		// start a routine that times the lobby out if too much time has passed since it last updated
		checkTimeout := func(l *states.LobbyState) {
			for l != nil {
				if l.RegisteredInstance.IsTimeoutExpired() {
					log.Printf("Deleting lobby %s due to timeout", l.RoomCode)
					s.Lobbies.CompareAndDelete(l.RoomCode, l)
					break
				}
				time.Sleep(time.Minute) // sleep for some time to prevent high CPU usage and avoid tight looping
			}
		}
		go checkTimeout(lobby)
	}

	// return message with the lobby ID or containing the error message
	msg, err := structures.ToWrappedJSON(rq)
	return msg, err
}

// initialize a client's data on the server and return their id to the client for communication
func (s *ServerData) handleadmitplayer(addr net.Addr, msgBody []byte) ([]byte, error) {

	// decode the message body to get player's attributes
	var rq messages.AdmissionMessage
	structures.FromWrappedJSON(&rq, msgBody)
	inputAttributes := rq.Attributes

	// create a new player on the server's player map
	newPlayer := states.NewPlayer(addr)
	newPlayer.PlayerAttributes = inputAttributes
	s.Players.LoadOrStore(newPlayer.GUID, newPlayer)

	// return message with the player's ID or containing the error message
	retrq := messages.AdmissionMessage{
		ClientPlayerID: rq.ClientPlayerID,
		ServerPlayerID: newPlayer.GUID,
	}
	msg, err := structures.ToWrappedJSON(retrq)
	return msg, err
}

// process a player add to game request
func (s *ServerData) handleaddplayergame(conn *websocket.Conn, msgBody []byte) ([]byte, error) {

	// deserialize the message
	var rq messages.AddPlayerGameMessage
	structures.FromWrappedJSON(&rq, msgBody)

	// decode the message body
	serverPlayerID := rq.ServerPlayerID
	gameID := rq.GameID

	// find the player's ID on the player map
	player, pErr := s.FindPlayer(serverPlayerID)
	if pErr != nil {
		return nil, fmt.Errorf("could not find player id in registry: %s", serverPlayerID)
	}
	player.GameID = gameID
	player.UpdateTime()

	// find the game's ID on the game map
	game, gErr := s.FindGame(gameID)
	if gErr != nil {
		return nil, fmt.Errorf("could not find game id in registry: %s", gameID)
	}

	// send back existing players
	s.sendGamePlayerIncludes(conn, &game.RegisteredInstance)

	// store the new player
	game.Players.LoadOrStore(serverPlayerID, true)
	game.UpdateTime()

	// broadcast their inclusion into the game
	s.broadcastPlayerJoined(&game.RegisteredInstance, player)

	// assign them as host if there is none
	s.assignHostIfNone(&game.RegisteredInstance, player)

	// respond by echoing the message
	return structures.ToWrappedJSON(rq)
}

// check if a given room code corresponds to a lobby that exists
func (s *ServerData) handlechecklobby(msgBody []byte) ([]byte, error) {
	var rq messages.CheckLobbyMessage
	structures.FromWrappedJSON(&rq, msgBody)

	// decode the message body
	roomCode := rq.RoomCode
	response := messages.CheckLobbyMessage{
		RoomCode: roomCode,
	}
	if len(roomCode) != states.NumRoomCodeChars {
		response.Exists = false
	} else {

		// search for the lobby and figure out whether it exists
		_, err := s.FindLobby(roomCode)
		response.Exists = err == nil
	}
	return structures.ToWrappedJSON(response)
}

// process a request from a player to switch sides
func (s *ServerData) handleswitch(msgBody []byte) ([]byte, error) {
	var rq messages.SwitchSideMessage
	structures.FromWrappedJSON(&rq, msgBody)
	pguid := rq.ServerPlayerID
	roomCode := rq.RoomCode

	// find lobby
	lobby, err := s.FindLobby(roomCode)
	if err != nil {
		return nil, fmt.Errorf("unable to find lobby to switch player in")
	}

	// find player
	player, err := s.FindPlayer(pguid)
	if err != nil {
		return nil, fmt.Errorf("unable to find player in player map during switch request")
	}

	// check that they match
	if player.RoomCode == lobby.RoomCode {

		// process the switch by pushing a forced update and broadcasting the new position
		player.Pos.X *= -1
		player.FaceRight = !player.FaceRight

	} else {
		return nil, fmt.Errorf("player id %s not found in lobby %s during switch request", pguid, roomCode)
	}

	// broadcast an update with the player's new position
	msg, err := structures.ToWrappedJSON(messages.PlayerActionMessage{
		PlayerServerID: player.GUID,
		Action:         player.PlayerAction,
		RoomCode:       lobby.RoomCode,
	})
	if err != nil {
		return nil, err
	} else {
		s.broadcastws(msg, &lobby.RegisteredInstance)
	}

	// send a forced update back to the client to switch the user
	return structures.ToWrappedJSON(messages.ForcePlayerMessage{
		ServerPlayerID: player.GUID,
		Action:         player.PlayerAction,
	})
}

// process a player add to lobby request
func (s *ServerData) handleaddplayerlobby(conn *websocket.Conn, msgBody []byte) ([]byte, error) {
	var rq messages.AddPlayerLobbyMessage
	structures.FromWrappedJSON(&rq, msgBody)

	// decode the message body
	serverPlayerID := rq.ServerPlayerID
	roomCode := rq.RoomCode

	// find the player's ID on the player map
	player, pErr := s.FindPlayer(serverPlayerID)
	if pErr != nil {
		return nil, fmt.Errorf("could not find player id in registry: %s", serverPlayerID)
	}
	player.RoomCode = roomCode
	player.UpdateTime()

	// find the lobby's ID on the lobby map
	lobby, lErr := s.FindLobby(roomCode)
	if lErr != nil {
		return nil, fmt.Errorf("could not find lobby with room code: %s", roomCode)
	}

	// autoassign them to a team and a position on the court
	isRightTeam := s.computeNewPlayerTeam(lobby)
	player.PlayerAction.Pos.X = computeRandomPosX(isRightTeam)
	player.PlayerAction.FaceRight = player.PlayerAction.Pos.X < 0
	msgForcePosition, err := structures.ToWrappedJSON(messages.ForcePlayerMessage{
		Action:         player.PlayerAction,
		ServerPlayerID: player.GUID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send player spawn loation: %s", err)
	}
	s.sendws(conn, msgForcePosition)

	// send the background image resource name to the client
	s.sendCurrentBackdrop(conn, lobby)

	// store the player to the lobby
	lobby.Players.LoadOrStore(serverPlayerID, true)
	lobby.UpdateTime()

	// send back a list of existing players in the lobby
	s.sendGamePlayerIncludes(conn, &lobby.RegisteredInstance)

	// broadcast their inclusion into the game
	s.broadcastPlayerJoined(&lobby.RegisteredInstance, player)

	// assign them as host if there is none
	s.assignHostIfNone(&lobby.RegisteredInstance, player)

	return structures.ToWrappedJSON(rq)
}

// handle a request to remove a player from the lobby
func (s *ServerData) handleleavelobby(msgBody []byte) ([]byte, error) {
	var rq messages.LeaveLobbyMessage
	structures.FromWrappedJSON(&rq, msgBody)
	roomCode := rq.RoomCode
	pid := rq.PlayerServerID
	s.removePlayerLobby(pid, roomCode)
	return msgBody, nil
}

// handle a request to remove a player from the game
func (s *ServerData) handleleavegame(msgBody []byte) ([]byte, error) {
	var rq messages.LeaveGameMessage
	structures.FromWrappedJSON(&rq, msgBody)
	gameID := rq.GameID
	pid := rq.PlayerServerID
	s.removePlayerGame(pid, gameID)
	return msgBody, nil
}

// handle a request to change the backdrop in a lobby
func (s *ServerData) handlesetbackdrop(msgBody []byte) ([]byte, error) {
	var rq messages.SetBackdropMessage
	structures.FromWrappedJSON(&rq, msgBody)
	lobby, err := s.FindLobby(rq.RoomCode)
	if err != nil {
		return nil, fmt.Errorf("could not find lobby with room code in registry: %s", rq.RoomCode)
	}
	lobby.Backdrop = rq.ResourceName
	s.broadcastws(msgBody, &lobby.RegisteredInstance)
	return nil, nil
}

// process a player action received from the client
func (s *ServerData) handleplayeraction(msgBody []byte) ([]byte, error) {

	// deserialize the message
	var amsg messages.PlayerActionMessage
	structures.FromWrappedJSON(&amsg, msgBody)
	gameID := amsg.GameID
	playerID := amsg.PlayerServerID
	roomCode := amsg.RoomCode

	// update the player's action stored on the server
	if len(playerID) > 0 {

		player, pErr := s.FindPlayer(playerID)
		if pErr != nil {
			return nil, fmt.Errorf("could not find player id in registry: %s", playerID)
		}
		player.UpdatePlayerState(&amsg.Action)
		player.UpdateTime()
	}

	// find the game that it applies to
	if len(gameID) > 0 {

		game, err := s.FindGame(gameID)
		if err != nil {
			return nil, fmt.Errorf("could not find game id in registry: %s", gameID)
		}
		game.UpdateTime()

		// just broadcast the action to all clients in the game
		s.broadcastws(msgBody, &game.RegisteredInstance)
		return msgBody, nil
	}

	// or find the lobby that it applies to
	if len(roomCode) > 0 {

		lobby, err := s.FindLobby(roomCode)
		if err != nil {
			return nil, fmt.Errorf("could not find lobby with room code in registry: %s", roomCode)
		}
		lobby.UpdateTime()

		// just broadcast the action to all clients in the lobby
		s.broadcastws(msgBody, &lobby.RegisteredInstance)
		return msgBody, nil
	}

	// somehow if the message was empty on both gameID and roomcode..? send this error
	return msgBody, fmt.Errorf("could not find matching instance for player action (nil gameID and roomCode received?)")
}

// process a ball event received from the client
func (s *ServerData) handleballevent(msgBody []byte) ([]byte, error) {

	// deserialize the message
	var bsm messages.BallStateMessage
	structures.FromWrappedJSON(&bsm, msgBody)
	clientBall := bsm.Ball
	gameID := bsm.GameID
	structures.FromWrappedJSON(&clientBall, msgBody)

	// find the game that it applies to
	game, err := s.FindGame(gameID)
	if err != nil {
		return nil, fmt.Errorf("could not find game id in registry: %s", gameID)
	}
	game.UpdateTime()

	// if for whatever reason the client's copy of the ball is out of date (e.g. someone else has registered a hit before them or the ball has already died), do not process the request and return a harmless error to the client
	denyBallUpdate := func(reason string) ([]byte, error) {
		err := fmt.Errorf("ball touch request denied, reason: %s", reason)
		log.Printf("Ball touch denied from: %s; reason: %s", clientBall.TouchedBy, reason)
		return nil, err
	}

	// accept the ball update and broadcast it
	acceptBallUpdate := func(b *states.BallState, gameID string) ([]byte, error) {
		ballMsg := messages.BallStateMessage{
			Ball:   *b,
			GameID: gameID,
		}
		sendMsg, err := structures.ToWrappedJSON(ballMsg)
		if err == nil {
			s.broadcastws(sendMsg, &game.RegisteredInstance)
		} else {
			log.Printf("Error broadcasting accepted ball status: %s", err)
		}
		return nil, err
	}

	// grab a local copy of the game ball
	cachedGameBall := game.GetBallCopy()

	// begin processing
	if len(clientBall.GUID) == 0 {

		// handle new ball registry
		clientBall.GenerateGUID()

		// register it to the game
		if cachedGameBall == nil {
			game.UpdateBall(&clientBall)
			log.Printf("Logged new game ball on server : %s", clientBall.GUID)
			return acceptBallUpdate(&clientBall, game.GUID)
		} else {
			return denyBallUpdate("A live game ball already exists")
		}

	} else {

		// check if ball id matches the one that is live on the server
		matchesLiveID := cachedGameBall != nil && cachedGameBall.GUID == clientBall.GUID && cachedGameBall.IsAlive()
		if !matchesLiveID {
			return denyBallUpdate("Ball ID doesn't match")
		}

		// check whether the client's ball update indicates that the ball is alive
		isAlive := clientBall.IsAlive()
		if isAlive {

			// if the ball is still alive, it means the player touched it; check if the touch count makes sense
			isTouchCountCorrect := clientBall.TouchCount <= 1 || (clientBall.TouchCount-cachedGameBall.TouchCount == 1)
			if !isTouchCountCorrect {
				return denyBallUpdate(fmt.Sprintf("Touch count incorrect: %d (client) vs %d (server)", clientBall.TouchCount, cachedGameBall.TouchCount))
			}

			// broadcast the updated client ball to other players
			game.UpdateBall(&clientBall)
			return acceptBallUpdate(&clientBall, game.GUID)

		} else {

			// it's possible that the game ball already died and has been set to nil
			if cachedGameBall == nil {
				return denyBallUpdate("Game ball already died or doesn't exist")
			}

			// if game ball was alive but client says it's dead, broadcast the dead ball and kill the ball on game side
			game.UpdateBall(nil)
			return acceptBallUpdate(&clientBall, game.GUID)
		}
	}
}
