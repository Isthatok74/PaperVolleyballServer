package states

import (
	"fmt"
	"net/http"
	"pv-server/util"
	"sync"
	"time"
)

// Purpose: Tracks vital metrics for the server.

// ServerState holds the state of the server
type ServerState struct {
	BaseState
	StartTime     string        // the timestamp on which the server was started
	ReqCount      int           // the number of times requests have been processed
	BytesReceived uint64        // the amount of data received since the server started
	BytesSent     uint64        // the amount of data sent since the server started
	mu            sync.Mutex    // To safely increment the ping count in concurrent requests
	ShutdownCh    chan struct{} // basically a listener which shuts down the server once it gets tripped (via `close(ShutdownCh)`)
}

// initialize a new ServerState object and return its pointer
func NewServerState() *ServerState {
	serverState := &ServerState{
		ShutdownCh:    make(chan struct{}),
		StartTime:     util.CurrentTimeUTC().Format(time.RFC3339),
		ReqCount:      0,
		BytesReceived: 0,
		BytesSent:     0,
	}
	serverState.GetGUID()
	return serverState
}

// dynamically counts the number of requests to the server
func (s *ServerState) CountRequests() {

	// increment the request tally (with thread-safe implementation of mutex)
	s.mu.Lock()
	s.ReqCount++
	s.mu.Unlock()

	// a hard breaker prevent a server from having to process too many requests
	if s.ReqCount > MaxRequests {
		close(s.ShutdownCh) // Signal to shut down
	}
}

// dynamically counts the bytes received at the server
func (s *ServerState) CountBytesReceived(v uint64) {
	s.mu.Lock()
	s.BytesReceived += v
	s.mu.Unlock()
}

// dynamically counts the bytes sent by the server
func (s *ServerState) CountBytesSent(v uint64) {
	s.mu.Lock()
	s.BytesSent += v
	s.mu.Unlock()
}

// define a constant number of requests to the server, past which the server will automatically shut down to mitigate further damages
const (
	MaxRequests int = 100000
	// todo: this quota should reset after a specified time period, so that the server can continue running indefinitely unless any such issues arise
)

// listens for a shutdown call, on which the server will immdiately attempt to shut down
// * this can be used if any undesirable situations are detected, such as bandwidth being consumed at abnormally high rates
func (s *ServerState) ListenForShutdown() {

	// Wait for the shutdown signal
	<-s.ShutdownCh // block until the shutdown signal is received
	fmt.Println("Shutting down the server...")

	{
		// insert cleanup logic here as necessary
	}

	fmt.Println("Server has shut down.")
}

// returns the size of the http request header
func (s *ServerState) HTTPHeaderBytesReceived(r *http.Request) {
	s.CountBytesReceived(util.HTTPHeaderSize(r.Header))
}

// returns the size
func (s *ServerState) HTTPHeaderBytesSent(w http.ResponseWriter) {

	// estimate HTTP status line size (e.g., "HTTP/1.1 200 OK\r\n")
	statusAssumed := http.StatusOK
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusAssumed, http.StatusText(statusAssumed))
	statusSize := uint64(len(statusLine))

	// estimate Content-Type line size
	contentTypeAssumed := "text/plain"
	contentTypeLine := fmt.Sprintf("Content-Type: " + contentTypeAssumed + "\r\n")
	contentTypeSize := uint64(len(contentTypeLine))

	// estimate Content-Length line size
	contentLengthAssumed := 9999
	contentLengthLine := fmt.Sprintf("Content-Length: %d\r\n", contentLengthAssumed)
	contentLengthSize := uint64(len(contentLengthLine))

	// get the custom headers size
	customHeader := w.Header()
	customHeaderSize := util.HTTPHeaderSize(customHeader)

	// sum up
	s.CountBytesSent(statusSize + contentTypeSize + contentLengthSize + customHeaderSize)
}
