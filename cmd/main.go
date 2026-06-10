// main.go
package main

import (
	"gotickets/internal/config"
	"gotickets/internal/server"
)

func main() {
	// Load environment variables
    cfg := config.LoadEnv()
	// Connect to the database
	db:= config.ConnectDB(cfg)
	// Start the server
	server.Start(db, cfg)
}
