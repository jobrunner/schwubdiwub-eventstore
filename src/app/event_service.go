package app

import (
	"eventstore/domain"
)

// EventService defines a application service interface for handling events
type EventService struct {
	Repo domain.EventRepository
}

// NewEventService creates a new instance of EventService
func NewEventService(repo domain.EventRepository) *EventService {
	return &EventService{Repo: repo}
}

// AppendEvent appends a new event to the repository
func (s *EventService) AppendEvent(event domain.Event) error {
	return s.Repo.Append(event)
}

// AppendEvents appends new events to the repository
func (s *EventService) AppendEvents(events []domain.Event) error {
	return s.Repo.AppendAll(events)
}

// GetEvents retrieves events with optional pagination
func (s *EventService) GetEvents(start, limit int) ([]domain.Event, error) {
	return s.Repo.GetAll(start, limit)
}
