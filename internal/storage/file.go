package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"shorter/internal/models"
	"shorter/internal/urlkey"
	"strconv"
	"strings"
)

type Row struct {
	ID          string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
	UserID      string `json:"userid"`
	DeletedFlag bool   `json:"deleted"`
}

type FileStorage struct {
	filePath string
	file     *os.File
	encoder  *json.Encoder
	counter  int // Tracks the number of stored records
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	err := makeDirInPath(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Count existing lines (records)
	count, err := countLines(filePath)
	if err != nil {
		count = 0 // Default to 0 if an error occurs
	}

	// Close the file if an error occurs
	defer func() {
		if err != nil {
			file.Close()
		}
	}()

	return &FileStorage{
		filePath: filePath,
		file:     file,
		encoder:  json.NewEncoder(file),
		counter:  count,
	}, nil
}

func (f *FileStorage) Set(OriginalURL string, userID string) (string, error) {

	urlKey := urlkey.GenerateSlug(OriginalURL)
	if urlKey == "" {
		return "", fmt.Errorf("shortURL is empty")
	}
	//Check for duplications
	if storedURL, _ := f.Get(urlKey); storedURL != "" {
		err := fmt.Errorf("the URL: %s is already stored in the file", storedURL)
		return urlKey, NewStorageError("already exists", storedURL, urlKey, err)
	}

	rowID := strconv.Itoa(f.counter + 1)

	row := Row{
		ID:          rowID,
		UserID:      userID,
		ShortURL:    urlKey,
		OriginalURL: OriginalURL,
	}

	// Write JSON entry
	err := f.encoder.Encode(row)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %s", err)
	}
	f.counter++
	return urlKey, nil
}

func (f *FileStorage) SetBatch(jReqBatch []models.JSONReq, userID string) ([]models.JSONRes, error) {
	jResBatch := make([]models.JSONRes, len(jReqBatch))

	for _, el := range jReqBatch {
		ShortURL, err := f.Set(el.OriginalURL, userID)
		if err != nil {
			return nil, err
		}

		row := models.JSONRes{
			CorrID:      el.CorrID,
			ShortURL:    ShortURL,
			OriginalURL: el.OriginalURL,
		}
		jResBatch = append(jResBatch, row)
	}
	return jResBatch, nil
}

func (f *FileStorage) Get(ShortURL string) (string, error) {
	ShortURL = strings.ToLower(ShortURL)
	if ShortURL == "" {
		return "", fmt.Errorf("shortURL is empty")
	}

	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %s", err)
	}
	// Search for the short URL
	for _, line := range splitLines(string(data)) {
		var row Row
		err := json.Unmarshal([]byte(line), &row)
		if err == nil && row.ShortURL == ShortURL {
			return row.OriginalURL, nil
		}
	}
	return "", fmt.Errorf("failed to find OriginalURL by ShortURL: %s", ShortURL)
}

func (f *FileStorage) DeleteBatch(ShortURLs []string, userID string) (bool, error) {
	if len(ShortURLs) == 0 {
		return false, errors.New("no URLs provided for deletion")
	}

	// Read the file
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read file: %w", err)
	}

	// Convert the file data into lines
	lines := splitLines(string(data))
	updatedLines := make([]string, 0, len(lines))
	var updatedCount int

	// Create a set for quick lookup of URLs to delete
	urlSet := make(map[string]struct{}, len(ShortURLs))
	for _, url := range ShortURLs {
		urlSet[url] = struct{}{}
	}

	// Process each line
	for _, line := range lines {
		var row Row
		if err := json.Unmarshal([]byte(line), &row); err != nil {
			// Skip malformed JSON lines
			updatedLines = append(updatedLines, line)
			continue
		}

		// Check if the URL belongs to the user and should be deleted
		if row.UserID == userID {
			if _, exists := urlSet[row.ShortURL]; exists {
				row.DeletedFlag = true
				updatedCount++
			}
		}

		// Convert back to JSON and store the updated record
		updatedJSON, err := json.Marshal(row)
		if err != nil {
			return false, fmt.Errorf("failed to serialize updated record: %w", err)
		}
		updatedLines = append(updatedLines, string(updatedJSON))
	}

	// Write the updated data back to the file
	err = os.WriteFile(f.filePath, []byte(strings.Join(updatedLines, "\n")), 0644)
	if err != nil {
		return false, fmt.Errorf("failed to write updated file: %w", err)
	}

	// Return true if at least one record was updated
	return updatedCount > 0, nil
}

func (f *FileStorage) GetUserURLs(userID string) ([]models.JSONUserRes, error) {
	jResBatch := make([]models.JSONUserRes, 0)

	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	// Search all records that were created by user
	for _, line := range splitLines(string(data)) {
		var row Row
		err := json.Unmarshal([]byte(line), &row)
		if err == nil && row.UserID == userID {
			row := models.JSONUserRes{
				UserID:      row.UserID,
				ShortURL:    row.ShortURL,
				OriginalURL: row.OriginalURL,
			}
			jResBatch = append(jResBatch, row)
		}
	}
	if len(jResBatch) == 0 {
		return nil, fmt.Errorf("failed to find records for userID: %s", userID)
	}
	return jResBatch, nil
}

// Close the file when FileStorage is no longer needed
func (f *FileStorage) Close() error {
	if f.file != nil {
		return f.file.Close()
	}
	return nil
}

func (f *FileStorage) IsAvailable() bool {
	return f.file != nil
}

// makeDirInPath - creates directories to store the file
func makeDirInPath(filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	return nil
}

// countLines - counts lines in the file
func countLines(filePath string) (int, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	return len(splitLines(string(data))), nil
}

// splitLines - splits the lines by a new line symbol
func splitLines(data string) []string {
	lines := []string{}
	start := 0
	for i, c := range data {
		if c == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	//Get the last line in the file if it doesn't have \n
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
