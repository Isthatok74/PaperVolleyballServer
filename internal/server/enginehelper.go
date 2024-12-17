package server

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/defs"
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/messages"
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/states"
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/structures"
	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/util"

	"github.com/gorilla/websocket"
)

// This file contains the helper functions for the engine

// helper function to assign a host to a registered instance if there is none specified
func (s *ServerData) assignHostIfNone(r *states.RegisteredInstance, player *states.PlayerState) {

	// assign host if none assigned
	if len(r.HostID) == 0 {
		r.HostID = player.GUID
		s.broadcastSyncHostMessage(r, r.HostID)
	}
}

// for a given leaving player, check if it is the host and reassign the host as necessary
func (s *ServerData) assignHostIfLeave(r *states.RegisteredInstance, leavingPlayerID string) {
	if leavingPlayerID == r.HostID {
		r.HostID = ""
		r.Players.Range(func(pid, _ interface{}) bool {
			if pid.(string) == leavingPlayerID {
				return true
			}
			r.HostID = pid.(string)
			return false
		})
		s.broadcastSyncHostMessage(r, r.HostID)
	}
}

// helper function to send one joining player's info to all connections in a registered instance
func (s *ServerData) broadcastPlayerJoined(r *states.RegisteredInstance, player *states.PlayerState) {
	includeMsg := messages.PlayerIncludeMessage{
		Attributes:     player.PlayerAttributes,
		Action:         player.PlayerAction,
		ServerPlayerID: player.GUID,
	}
	msg, err := structures.ToWrappedJSON(includeMsg)
	if err != nil {
		log.Printf("Unable to wrap PlayerIncludeMessage in a json: %s", err)
	} else {
		s.broadcastws(msg, r)
	}
}

// broadcast the new host in a lobby
func (s *ServerData) broadcastSyncHostMessage(r *states.RegisteredInstance, hostID string) {
	msg := messages.SyncHostMessage{
		HostID: hostID,
	}
	sendMsg, err := structures.ToWrappedJSON(msg)
	if err != nil {
		log.Printf("Unable to send new host message to lobby: %s", err)
	} else {
		s.broadcastws(sendMsg, r)
	}
}

// assigns a new player to the team with fewer players, or the left if both have same; returns the team that they are on; left = false, right = true
func (s *ServerData) computeNewPlayerTeam(l *states.LobbyState) bool {
	lCount, rCount := s.countTeamPlayers(&l.RegisteredInstance)
	return lCount > rCount
}

// a helper to return either 1 or -1 corresponding to the sides of the court
func computeSideMultiplier(isRightSide bool) float32 {
	var sideSign float32
	if isRightSide {
		sideSign = 1.0
	}
	if !isRightSide {
		sideSign = -1.0
	}
	return sideSign
}

// return a random x position on the court on the given side
func computeRandomPosX(isRightSide bool) float32 {
	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	randX := 1 + rand.Float64()*(defs.MaxCourtSpawnX-1)
	sideSign := computeSideMultiplier(isRightSide)
	return sideSign * float32(math.Abs(float64(randX)))
}

// count the number of players on either team
func (s *ServerData) countTeamPlayers(r *states.RegisteredInstance) (int, int) {
	lCount := 0
	rCount := 0
	r.Players.Range(func(pid, value interface{}) bool {
		ptr, err := s.FindPlayer(pid.(string))
		if err == nil {
			isRight := ptr.PlayerAction.Pos.X > 0
			if isRight {
				rCount++
			} else {
				lCount++
			}
		}
		return true
	})
	return lCount, rCount
}

// send the current backdrop's resource name in the lobby to a player, if null
func (s *ServerData) sendCurrentBackdrop(conn *websocket.Conn, lobby *states.LobbyState) {
	msgBackdrop, err := structures.ToWrappedJSON(messages.SetBackdropMessage{
		ResourceName: lobby.Backdrop,
		RoomCode:     lobby.RoomCode,
	})
	if err != nil {
		log.Printf("failed to send backdrop resource name: %s", err)
	} else {
		s.sendws(conn, msgBackdrop)
	}
}

// helper function to send data of all players in a game to a connection
func (s *ServerData) sendGamePlayerIncludes(conn *websocket.Conn, r *states.RegisteredInstance) {
	r.Players.Range(func(pid, _ interface{}) bool {
		peer, err := s.FindPlayer(pid.(string))
		if err != nil {
			log.Printf("Could not find expected player in game with id %s, player id: %s", r.GUID, pid.(string))
		}
		includeMsg := messages.PlayerIncludeMessage{
			Attributes:     peer.PlayerAttributes,
			Action:         peer.PlayerAction,
			ServerPlayerID: peer.GUID,
		}
		msg, err := structures.ToWrappedJSON(includeMsg)
		if err != nil {
			log.Printf("Unable to wrap PlayerIncludeMessage in a json: %s", err)
		} else {
			s.sendws(conn, msg)
		}
		return true
	})
}

// remove player from the specified game
func (s *ServerData) removePlayerGame(playerID string, gameID string) {

	if gameID != "" {
		game, err := s.FindGame(gameID)
		if err != nil {
			log.Println(err)
		} else {
			msg, err := structures.ToWrappedJSON(messages.LeaveGameMessage{
				PlayerServerID: playerID,
				GameID:         gameID,
			})
			if err != nil {
				log.Printf("Unable to send leave game message to game: %s", err)
			} else {

				// send an update to all players
				s.broadcastws(msg, &game.RegisteredInstance)
			}

			// remove from the instance's player map
			game.RegisteredInstance.Players.Delete(playerID)

			// delete the instance if no players remain
			if util.GetSyncMapSize(&game.RegisteredInstance.Players) == 0 {
				s.Games.Delete(gameID)
			} else {
				s.assignHostIfLeave(&game.RegisteredInstance, playerID)
			}

			// remove from the global player map
			s.Players.Delete(playerID)

		}
	}
}

// remove player from the specified lobby
func (s *ServerData) removePlayerLobby(playerID string, roomCode string) {

	if roomCode != "" {
		lobby, err := s.FindLobby(roomCode)
		if err != nil {
			log.Println(err)
		} else {
			msg, err := structures.ToWrappedJSON(messages.LeaveLobbyMessage{
				PlayerServerID: playerID,
				RoomCode:       roomCode,
			})
			if err != nil {
				log.Printf("Unable to send leave game message to lobby: %s", err)
			} else {

				// send an update to all players
				s.broadcastws(msg, &lobby.RegisteredInstance)
			}

			// remove from the instance's player map
			lobby.RegisteredInstance.Players.Delete(playerID)

			// remove from instance's player map
			if util.GetSyncMapSize(&lobby.RegisteredInstance.Players) == 0 {
				s.Lobbies.Delete(roomCode)
			} else {
				s.assignHostIfLeave(&lobby.RegisteredInstance, playerID)
			}

			// remove from the global player map
			s.Players.Delete(playerID)
		}
	}
}

// handle loss of a connection for any reason
func (s *ServerData) processdisconnect(conn *websocket.Conn) {

	// find players with matching connection
	var lstRemove []string
	s.Players.Range(func(pid, value interface{}) bool {

		// ensure type
		ptr, ok := value.(*states.PlayerState)
		if !ok {
			// Handle the error if the value is not of the expected type
			log.Printf("Value for key %v is not of type *states.PlayerState\n", pid)
			return true // Continue the iteration
		}

		// check for matched address
		if ptr.GetAddress().String() == conn.RemoteAddr().String() {

			// this is the player that disconnected; mark for deletion
			lstRemove = append(lstRemove, ptr.GUID)

			// remove from the player's game if it exists
			gameID := ptr.GameID
			s.removePlayerGame(ptr.GUID, gameID)

			// remove from the player's lobby if it exists
			roomCode := ptr.RoomCode
			s.removePlayerLobby(ptr.GUID, roomCode)
		}
		return true
	})

	// remove the player from the player map (or multiple in case multiple registries were made)
	for _, pid := range lstRemove {
		s.Players.Delete(pid)
	}
}
