package server

import (
	"sync"
)

// Purpose: A container for all the data tracked by the server in real time

type ServerData struct {
	Info        ServerState // vitals
	Games       sync.Map    // a map of all ongoing games hosted on this server (key: game.GUID, value: *states.gameState)
	Lobbies     sync.Map    // a map of all ongoing lobbies hosted on this server (key: lobby.RoomCode, value: *states.lobbyState)
	Connections sync.Map    // a map of all live connections established on this server (key: conn.RemoteAddr(), value: *websocket.Conn)
	Players     sync.Map    // a map of all connected clients hosted on this server (key: player.GUID, value: *states.playerState)
}

// constructor function to initialize ServerData
func NewServerData() *ServerData {
	serverData := &ServerData{
		Info: *NewServerState(), // Initialize Info field with zero value
	}
	return serverData
}

// check if a lobby with the specified roomcode exists
func (s *ServerData) LobbyExists(roomCode string) bool {
	_, exists := s.Lobbies.Load(roomCode)
	return exists
}
