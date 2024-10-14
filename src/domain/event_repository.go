package domain

// EventRepository defines the interface for event storage
type EventRepository interface {
	Append(event Event) error
	AppendAll(event []Event) error
	GetAll(start, limit int) ([]Event, error)
}
