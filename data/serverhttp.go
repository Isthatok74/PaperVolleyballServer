package data

import (
	"fmt"
	"net/http"
	"pv-server/util"
)

// Handles the http traffic portion of the server

// handle the status route on http - returns some server metrics
func (s *ServerData) HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server start time: %s \n", s.Info.StartTime)
	fmt.Fprintf(w, "Number of requests processed: %d \n", s.Info.ReqCount)
	fmt.Fprintf(w, "Number of active games: %d \n", util.GetSyncMapSize(&(s.Games)))
	fmt.Fprintf(w, "Number of clients connected: %d\n", util.GetSyncMapSize(&(s.Clients)))
	fmt.Fprintf(w, "Data received: %s \n", util.FormatBytes(s.Info.BytesReceived))
	fmt.Fprintf(w, "Data sent: %s \n", util.FormatBytes(s.Info.BytesSent))
}
