package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/Isthatok74/PaperVolleyballServer/data/states"
)

// The rate limit handler is a precautionary middleware that limits the number of requests that a client can make to the server over a specified time period.
// This is useful to mitigate damages in the event of DDoS attacks

var (
	rateLimitMap sync.Map      // a map of all clients that have registered any type of request to the server
	rateLimit    = 300         // max requests per the time window
	limitWindow  = time.Minute // the time window, after which the rate quota gets reset
)

// Rate limit structure for each client
type clientRate struct {
	count     int       // the number of queries that the client has made within the time window
	timestamp time.Time // a tracker to help check how much time has elapsed in the time window
}

// prevent any single client from sending too many requests in close succession
func rateLimitHandler(next http.Handler, s *states.ServerState) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// define a header writing function
		writeHeader := func(statusCode int) {
			w.WriteHeader(statusCode)
			s.HTTPHeaderBytesSent(w)
		}

		// also track the total number of requests that have been made to the server
		s.CountRequests()

		// get the client's IP address
		clientIP := r.RemoteAddr

		// retrieve rate limit for this client IP
		value, _ := rateLimitMap.LoadOrStore(clientIP, &clientRate{timestamp: time.Now()})
		clientRate := value.(*clientRate)

		// check if the window has expired and reset if necessary
		if time.Since(clientRate.timestamp) > limitWindow {
			clientRate.count = 0
			clientRate.timestamp = time.Now()
		}

		// increment the request count
		clientRate.count++

		// if the limit is reached, block further requests and return an error
		if clientRate.count > rateLimit {
			writeHeader(http.StatusTooManyRequests)
			return
		}

		// tally the data received in the header
		s.HTTPHeaderBytesReceived(r)

		// estimate the data that will be sent in the header
		s.HTTPHeaderBytesSent(w)

		// continue to process the actual handler function for the query
		next.ServeHTTP(w, r)
	})
}

// periodically resets the quota of requests for each client that has connected by any means
func resetRateLimit() {
	for {
		time.Sleep(limitWindow)
		rateLimitMap.Range(func(key, value interface{}) bool {
			clientRate := value.(*clientRate)
			clientRate.count = 0 // Reset count for all clients
			clientRate.timestamp = time.Now()
			return true
		})
	}
}
