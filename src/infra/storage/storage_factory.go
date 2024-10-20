package storage

import (
	"eventstore/config"
	"eventstore/core"
	"fmt"
)

func NewEventStoreRepository(cfg *config.Config) (core.StorageRepository, error) {
	switch cfg.StorageType {
	case "memory":
		return NewMemoryStorageRepository(cfg)
	case "file":
		return NewFileStorageRepository(cfg)
	case "azure":
		return NewAzureBlobStorageRepository(cfg)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.StorageType)
	}
}
