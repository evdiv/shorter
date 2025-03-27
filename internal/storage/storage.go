package storage

import (
	"fmt"
	"shorter/internal/config"
	"shorter/internal/models"
)

type Storer interface {
	Set(url string, userID string) (string, error)
	SetBatch(entries []models.JSONReq, userID string) ([]models.JSONRes, error)
	DeleteBatch(keysToDelete []models.KeysToDelete) (bool, error)
	GetUserURLs(userID string) ([]models.JSONUserRes, error)
	Get(key string) (string, error)
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
