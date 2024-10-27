package app

import (
	"context"
	"eventstore/core"
)

// EventService defines an application service interface for handling events
type EventStoreService struct {
	Repo core.StorageRepository
}

// NewEventService creates a new instance of EventService
func NewEventStoreService(repo core.StorageRepository) *EventStoreService {
	return &EventStoreService{Repo: repo}
}

// AppendEvent appends a new event to the repository
func (s *EventStoreService) AppendEvent(ctx context.Context, event core.Event) error {
	if err := validateEvent(event); err != nil {
		return err
	}
	return s.Repo.Append(ctx, event)
}

// AppendEvents appends new events to the repository
func (s *EventStoreService) AppendEvents(ctx context.Context, events []core.Event) error {
	return s.Repo.AppendAll(ctx, events)
}

// GetEvents retrieves events with optional pagination
func (s *EventStoreService) GetEvents(ctx context.Context, start, limit int) ([]core.Event, error) {
	return s.Repo.GetAll(ctx, start, limit)
}

func validateEvent(event core.Event) error {
	if event.MessageId == "" {
		return core.ErrEventMissingMessageId
	}
	if event.EventType == "" {
		return core.ErrEventMissingEventType
	}
	if event.Timestamp == "" {
		return core.ErrEventMissingTimestamp
	}
	if event.Payload == "" {
		return core.ErrEventMissingPayload
	}
	return nil
}
