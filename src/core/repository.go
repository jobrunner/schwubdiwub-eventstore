package core

import "context"

type StorageRepository interface {
	Append(ctx context.Context, event Event) error
	AppendAll(ctx context.Context, event []Event) error
	GetAll(ctx context.Context, start, limit int) ([]Event, error)
}
