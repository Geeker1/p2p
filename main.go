package main

import (
	"github.com/Geeker1/p2p/server"
	// "github.com/Geeker1/p2p/tracker"
)

func main()  {
	// Start Tracker
	// Start Clients
	// Start Server

	// Wait for algo to be done, then exit application.
	// tracker.StartTracker()
	server.StartServer("/home/devcode/shing.mp4", 1000000)
}
