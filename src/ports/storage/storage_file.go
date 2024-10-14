package storage

import (
	"bufio"
	"encoding/json"
	"eventstore/domain"
	"os"
)

// FileRepository stores events in a local file
type FileRepository struct {
	filePath string
}

// NewFileRepository creates a new file-based repository
func NewFileRepository(filePath string) *FileRepository {
	return &FileRepository{filePath: filePath}
}

// Append adds a new event to the file
func (r *FileRepository) Append(event domain.Event) error {
	file, err := os.OpenFile(r.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = file.Write(append(data, '\n'))
	return err
}

// GetAll returns all events from the file with pagination
func (r *FileRepository) GetAll(start, limit int) ([]domain.Event, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []domain.Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event domain.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if start > len(events) {
		return nil, nil
	}
	end := start + limit
	if end > len(events) || limit == 0 {
		end = len(events)
	}
	return events[start:end], nil
}
