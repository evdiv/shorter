package storage

import "fmt"

type StorageError struct {
	Type        string
	OriginalURL string
	ShortURL    string
	Err         error
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("Type: %s, OriginalURL: %s, ShortURL: %s) Error: %s",
		e.Type, e.OriginalURL, e.ShortURL, e.Err)
}

func NewStorageError(errorType string, originalURL string, shortURL string, err error) *StorageError {
	return &StorageError{
		Type:        errorType,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		Err:         err,
	}
}
