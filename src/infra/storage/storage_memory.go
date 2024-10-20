package storage

import (
	"context"
	"eventstore/config"
	"eventstore/core"
	"log"
	"sync"
)

// MemoryRepository stores events in memory
type memoryStorageRepository struct {
	mu     sync.Mutex
	events []core.Event
}

// NewMemoryRepository creates a new memory repository
func NewMemoryStorageRepository(cfg *config.Config) (*memoryStorageRepository, error) {
	log.Printf("waring: using memory storage adapter with more than one instance is not recommended")
	return &memoryStorageRepository{}, nil
}

// Append adds a new event to the memory store
func (r *memoryStorageRepository) Append(ctx context.Context, event core.Event) error {
	return r.AppendAll(ctx, []core.Event{event})
}

// Append adds a new event to the memory store
func (r *memoryStorageRepository) AppendAll(ctx context.Context, events []core.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, events...)
	return nil
}

// GetAll returns all events, optionally with pagination
func (r *memoryStorageRepository) GetAll(ctx context.Context, start, limit int) ([]core.Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if start > len(r.events) {
		return nil, nil
	}
	end := start + limit
	if end > len(r.events) || limit == 0 {
		end = len(r.events)
	}
	return r.events[start:end], nil
}
