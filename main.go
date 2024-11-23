package main

// Deploying instructions:
// * Build command (run via Terminal): `go build`
// * You can test on local machine by simply running the .exe to start the server
// * For servers hosted on Windows, connect to the RDP instance (assuming it is a Windows server), copy the .exe into the server, and run it via Command Prompt
// * The url for a request will be <http or ws>://<serverAddress>:<port>/<command>. If running locally, the value of <serverAddress> is `localhost``. If deploying on the cloud, then it is the IP address of the server.

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
	http.Handle("/status", rateLimitHandler(http.HandlerFunc(serverData.HandleStatus), &(serverData.Info)))
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
