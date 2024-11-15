package main

import (
	"net/http"
	"sync"
	"time"
)

var (
	rateLimitMap sync.Map
	rateLimit    = 1000 // Max requests per minute
	limitWindow  = time.Minute
)

// Rate limit structure for each client
type clientRate struct {
	count     int
	timestamp time.Time
}

// prevent any single client from sending too many requests in close succession
func rateLimitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the client's IP address
		clientIP := r.RemoteAddr

		// Retrieve rate limit for this client IP
		value, _ := rateLimitMap.LoadOrStore(clientIP, &clientRate{timestamp: time.Now()})
		clientRate := value.(*clientRate)

		// Check if the window has expired and reset if necessary
		if time.Since(clientRate.timestamp) > limitWindow {
			clientRate.count = 0
			clientRate.timestamp = time.Now()
		}

		// Increment the request count
		clientRate.count++

		// If the limit is reached, block further requests
		if clientRate.count > rateLimit {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		// Continue with the actual handler
		next.ServeHTTP(w, r)
	})
}

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
