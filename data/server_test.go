package data

import (
	"pv-server/data/states"
	"testing"
)

// spam the server with concurrent requests until it reaches the limit
func TestLimitRequests(t *testing.T) {
	ss := states.NewServerState()

	for i := 1; i <= states.MaxRequests+1; i++ {
		go ss.CountRequests()
	}
	<-ss.ShutdownCh // block until the shutdown signal is received
	t.Logf("Shutdown channel successfully closed with %d requests", ss.ReqCount)
}
