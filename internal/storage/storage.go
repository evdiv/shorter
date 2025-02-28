package storage

import (
	"fmt"
	"shorter/internal/config"
)

type Storer interface {
	Set(url string) string
	Get(key string) string
	IsAvailable() bool
	Close() error
}

func NewStorage(appConfig config.Config) (Storer, error) {
	// If DB connection is provided, initialize DB storage
	if appConfig.DBConnection != "" {
		dbStorage, err := NewDBStorage(appConfig.DBConnection)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize database storage: %w", err)
		}
		err = dbStorage.Migrate()
		if err != nil {
			return nil, fmt.Errorf("failed to migrate database storage: %w", err)
		}
		return dbStorage, nil
	}
	// If FilePath is provided (but no DB), initialize file storage
	if appConfig.StoragePath != "" {
		fileStorage, err := NewFileStorage(appConfig.StoragePath)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize file storage: %w", err)
		}
		return fileStorage, nil
	}
	// Default to in-memory storage
	return NewMemoryStorage(), nil
}
