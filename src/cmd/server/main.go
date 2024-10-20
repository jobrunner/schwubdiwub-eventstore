package main

import (
	"eventstore/app"
	"eventstore/config"
	"eventstore/infra/rest"
	"eventstore/infra/storage"
	"log"
)

func main() {

	cfg := config.LoadConfig()

	repo, err := storage.NewEventStoreRepository(&cfg)
	if err != nil {
		log.Fatalf("Error initializing storage: %v", err)
	}

	service := app.NewEventStoreService(repo)
	server := rest.NewRestServer(cfg.ServerAddress, service)

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
