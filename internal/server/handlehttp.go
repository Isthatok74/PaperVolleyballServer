package server

import (
	"fmt"
	"net/http"

	"github.com/Isthatok74/PaperVolleyballServer/internal/pkg/util"
)

// Handles the http traffic portion of the server

// handle the status route on http - returns some server metrics
func (s *ServerData) HandleStatus(w http.ResponseWriter, r *http.Request) {
	s.WriteHTTP(w, fmt.Sprintf("Server start time: %s \n", s.Info.StartTime))
	s.WriteHTTP(w, fmt.Sprintf("Number of requests processed: %d \n", s.Info.ReqCount))
	s.WriteHTTP(w, fmt.Sprintf("Number of active lobbies: %d \n", util.GetSyncMapSize(&(s.Lobbies))))
	s.WriteHTTP(w, fmt.Sprintf("Number of active games: %d \n", util.GetSyncMapSize(&(s.Games))))
	s.WriteHTTP(w, fmt.Sprintf("Number of clients connected: %d\n", util.GetSyncMapSize(&(s.Clients))))
	s.WriteHTTP(w, fmt.Sprintf("Estimated data received: %s \n", util.FormatBytes(s.Info.BytesReceived)))
	s.WriteHTTP(w, fmt.Sprintf("Estimated data sent: %s \n", util.FormatBytes(s.Info.BytesSent)))
}

// return an empty page
func (s *ServerData) HandleDefault(w http.ResponseWriter, r *http.Request) {
	s.WriteHTTP(w, "")
}

// wrap the writer with keeping track of bandwidth sent
func (s *ServerData) WriteHTTP(w http.ResponseWriter, msg string) {
	fmt.Fprint(w, msg)
	s.Info.CountBytesSent(uint64(len(msg)))
}
