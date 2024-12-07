package server

import (
	"sync"
)

// Purpose: A container for all the data tracked by the server in real time

type ServerData struct {
	Info    ServerState // vitals
	Games   sync.Map    // a map of all ongoing games hosted on this server (key: game.GUID, value: *states.gameState)
	Clients sync.Map    // a map of all connected players hosted on this server (key: conn.RemoteAddr(), value: *websocket.Conn)
}

// constructor function to initialize ServerData
func NewServerData() *ServerData {
	serverData := &ServerData{
		Info: *NewServerState(), // Initialize Info field with zero value
	}
	return serverData
}
