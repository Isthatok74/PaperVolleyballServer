package data

import (
	"pv-server/data/states"
	"sync"
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
