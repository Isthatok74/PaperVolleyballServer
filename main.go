package main

// Deploying instructions:
// Build command: `go build`
// You can test on local machine by simply running the .exe to start the server
// For servers hosted on Windows, connect to the RDP instance (assuming it is a Windows server), copy the .exe into the server, and run it via Command Prompt
// The url for a request will be <http or ws>://<serverName>:<port>/<command>

import (
	"fmt"
	"net/http"
	"pv-server/data"
)

var (
	serverData data.ServerData // storage for the server data
)

func main() {

	fmt.Println("Starting server...")

	fmt.Println("Initiating server data...")
	serverData := data.NewServerData()

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

func setupRoutesHTTP() {
	http.Handle("/ping", rateLimitHandler(http.HandlerFunc(serverData.HandlePing)))
	http.Handle("/status", rateLimitHandler(http.HandlerFunc(serverData.HandleStatus)))

	// todo: remove obsolete
	http.Handle("/create", rateLimitHandler(http.HandlerFunc(serverData.HandleCreate)))
	http.Handle("/addplayer", rateLimitHandler(http.HandlerFunc(serverData.HandleAddPlayer)))
	http.Handle("/post", rateLimitHandler(http.HandlerFunc(serverData.HandlePost)))
}

func setupRoutesWS() {
	http.Handle("/ws", rateLimitHandler(http.HandlerFunc(serverData.HandleWS)))
}

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
