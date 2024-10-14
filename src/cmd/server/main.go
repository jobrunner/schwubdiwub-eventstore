package main

import (
	"eventstore/config"
	"eventstore/ports/rest"
	"eventstore/ports/storage"
	"log"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Create repository using StorageFactory
	repo, err := storage.NewStorageFactory(cfg)
	if err != nil {
		log.Fatalf("Error initializing storage: %v", err)
	}

	// Initialize and configure HTTP server
	server := rest.NewRestServer(cfg.ServerAddress, repo)
	server.ConfigureRoutes()

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
