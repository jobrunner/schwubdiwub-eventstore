package storage

import (
	"eventstore/config"
	"eventstore/domain"
	"fmt"
)

// NewStorageFactory returns the appropriate storage adapter based on the config
func NewStorageFactory(cfg config.Config) (domain.EventRepository, error) {
	switch cfg.StorageType {
	case "memory":
		return NewMemoryRepository(), nil
	case "file":
		return NewFileRepository(cfg.FilePath), nil
	// case "aws":
	// 	return NewS3Repository(cfg.AWSConfig), nil
	// case "azure":
	// 	return NewAzureBlobRepository(cfg.AzureConfig), nil
	// case ... FireStore, MySQL, Dolt what ever
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.StorageType)
	}
}
