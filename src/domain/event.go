package domain

// Event represents an event in the event store
type Event struct {
	Id        string `json:"id"`
	Timestamp string `json:"timestamp"`
	EventType string `json:"event_type"`
	Payload   string `json:"payload"`
}
