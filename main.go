package main

// Deploying instructions:
// Testing on local machine: `go run main.go`
// Build command: `go build -o MyServer.exe main.go`
// For servers hosted on Windows, connect to the RDP instance (assuming it is a Windows server), copy the .exe into the server, and run it via Command Prompt
// The url for a request will be http://<serverName>:<port>/<command>

import (
	"fmt"
	"net/http"
	"pv-server/data"
)

func main() {

	fmt.Println("Starting server...")

	fmt.Println("Initiating server data...")
	serverData := data.NewServerData()

	fmt.Println("Starting rate limiter...")

	go resetRateLimit()
	http.Handle("/", rateLimitHandler(http.DefaultServeMux))

	fmt.Println("Setting up function handlers...")
	http.Handle("/ping", rateLimitHandler(http.HandlerFunc(serverData.HandlePing)))
	http.Handle("/status", rateLimitHandler(http.HandlerFunc(serverData.HandleStatus)))
	http.Handle("/create", rateLimitHandler(http.HandlerFunc(serverData.HandleCreate)))
	http.Handle("/addplayer", rateLimitHandler(http.HandlerFunc(serverData.HandleAddPlayer)))
	http.Handle("/message", rateLimitHandler(http.HandlerFunc(serverData.HandlePostMessage)))

	fmt.Println("Attempting to start server...")
	startServer()

	fmt.Println("Setting up shutdown listener...")
	serverData.Info.ListenForShutdown()
}

func startServer() {

	// declare which port the server will be listening on
	port := "13274"
	address := ":" + port

	// start the server
	go func() {
		fmt.Println("Starting server on port" + port + "...")
		if err := http.ListenAndServe(address, nil); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()
}
