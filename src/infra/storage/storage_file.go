package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"eventstore/config"
	"eventstore/core"
	"log"
	"os"

	bloom "github.com/bits-and-blooms/bloom/v3"
)

// FileRepository stores events in a local file
type fileStorageRepository struct {
	filePath string
	cfg      *config.Config
}

// NewFileRepository creates a new file-based repository
func NewFileStorageRepository(cfg *config.Config) (*fileStorageRepository, error) {
	return &fileStorageRepository{
		cfg: cfg,
	}, nil
}

// Append adds a new event to the file
func (r *fileStorageRepository) Append(ctx context.Context, event core.Event) error {
	return r.AppendAll(ctx, []core.Event{event})
}

// Append adds a new event to the file
func (r *fileStorageRepository) AppendAll(ctx context.Context, events []core.Event) error {
	filePath := r.cfg.LocalFileConfig.FilePath
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var b []byte

	for _, event := range events {
		record, err := json.Marshal(event)
		if err != nil {
			return err
		}
		b = append(b, record...)
		b = append(b, '\n')
	}

	if b == nil {
		return nil
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

// GetAll returns all events from the file with pagination
func (r *fileStorageRepository) GetAll(ctx context.Context, start, limit int) ([]core.Event, error) {
	continueOnError := true
	filter := bloom.NewWithEstimates(r.cfg.EstimatedEventCount, 0.1)

	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []core.Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event core.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			if continueOnError {
				log.Println("Failed to decode event: ", err)
				continue
			}

			return nil, err
		}
		if !filter.Test([]byte(event.MessageId)) {
			events = append(events, event)
			filter.Add([]byte(event.MessageId))
		} else {
			log.Println("Ignoring potential duplicate event found: ", event.MessageId)
			if limit > 0 {
				limit++
			}
		}
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
