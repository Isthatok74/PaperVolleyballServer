package states

import (
	"fmt"
	"pv-server/util"
	"sync"
	"time"
)

// ServerState holds the state of the server
type ServerState struct {
	BaseState
	StartTime  string     // the timestamp on which the server was started
	ReqCount   int        // the number of times requests have been processed
	mu         sync.Mutex // To safely increment the ping count in concurrent requests
	ShutdownCh chan struct{}
}

func NewServerState() *ServerState {
	serverState := &ServerState{
		ShutdownCh: make(chan struct{}),
		StartTime:  util.CurrentTimeUTC().Format(time.RFC3339),
		ReqCount:   0,
	}
	serverState.GetGUID()
	return serverState
}

// Dynamically counts the number of requests to the server
func (s *ServerState) CountRequests() {
	s.mu.Lock()
	s.ReqCount++
	s.mu.Unlock() // Unlock the mutex

	// prevent a server from having to process too many requests
	if s.ReqCount > MaxRequests {
		close(s.ShutdownCh) // Signal to shut down
	}
}

const (
	MaxRequests int = 100000
)

func (s *ServerState) ListenForShutdown() {

	// Wait for the shutdown signal
	<-s.ShutdownCh // Block until the shutdown signal is received
	fmt.Println("Shutting down the server...")

	{
		// insert cleanup logic here as necessary
	}

	fmt.Println("Server has shut down.")
}
