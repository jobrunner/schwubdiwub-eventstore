package main

import (
	"context"
	"eventstore/app"
	"eventstore/config"
	"eventstore/infra/storage"
	"fmt"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	repo, err := storage.NewEventStoreRepository(&cfg)
	if err != nil {
		log.Fatalf("Error initializing storage: %v", err)
	}

	ctx := context.TODO()
	service := app.NewEventStoreService(repo)
	events, _ := service.GetEvents(ctx, 0, 0)
	for _, event := range events {
		fmt.Println(event)
	}
}
