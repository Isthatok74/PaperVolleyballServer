package main

import (
	"fmt"
	"net/http"
	"pv-server/data"
)

// global variable to store the server data throughout lifetime of server
var serverData = data.NewServerData()

// this is the main function of the server, which runs when the program begins
func main() {

	fmt.Println("Starting server...")

	fmt.Println("Starting rate limiter...")
	go resetRateLimit()

	fmt.Println("Setting up function handlers...")
	setupRoutesHTTP()
	setupRoutesWS()

	fmt.Println("Attempting to start server...")
	startServer()

	fmt.Println("Setting up shutdown listener...")
	serverData.Info.ListenForShutdown()
}

// all of the HTTP routes are defined here.
// * HTTP requests are used for requests that can be made from anywhere (e.g. a web browser). They typically involve a client simply sending a request, processing the request on the server, and then sending back a message to the client.
func setupRoutesHTTP() {

	// check the status of the server
	http.Handle("/status", rateLimitHandler(http.HandlerFunc(serverData.HandleStatus), &(serverData.Info)))

	// any other route should still go through the middleware for checks
	http.Handle("/", rateLimitHandler(http.HandlerFunc(serverData.HandleDefault), &(serverData.Info)))
}

// all of the WebSocket routes are defined here
// * WebSockets are used for fast and frequent communication between client and server. A connection line is established over perpetual listeners are set up between both sides. Whenever data is transferred, there is little overhead compared to HTTP (which requires writing a header every time data is transferred).
func setupRoutesWS() {
	http.Handle("/ws", rateLimitHandler(http.HandlerFunc(serverData.HandleWS), &(serverData.Info)))
}

// start the server by setting up a listener on the specified port
func startServer() {

	// declare which port the server will be listening on
	port := "13274"
	address := ":" + port

	// start the server
	go func() {
		fmt.Println("Starting server on port " + port + "...")
		if err := http.ListenAndServe(address, nil); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()
}
