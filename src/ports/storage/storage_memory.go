package storage

import (
	"eventstore/domain"
	"sync"
)

// MemoryRepository stores events in memory
type MemoryRepository struct {
	mu     sync.Mutex
	events []domain.Event
}

// NewMemoryRepository creates a new in-memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

// Append adds a new event to the in-memory store
func (r *MemoryRepository) Append(event domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, event)
	return nil
}

// GetAll returns all events, optionally with pagination
func (r *MemoryRepository) GetAll(start, limit int) ([]domain.Event, error) {
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
