package core

import (
	"errors"
)

var (
	ErrEventMissingMessageId = errors.New("event is missing MessageId")
	ErrEventMissingEventType = errors.New("event is missing EventType")
	ErrEventMissingTimestamp = errors.New("event is missing Timestamp")
	ErrEventMissingPayload   = errors.New("event is missing Payload")
)

// Event represents an event in the event store
type Event struct {
	MessageId string `json:"message_id"`
	Timestamp string `json:"timestamp"`
	EventType string `json:"event_type"`
	Payload   string `json:"payload"`
}
