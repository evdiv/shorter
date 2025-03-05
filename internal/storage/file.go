package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"shorter/internal/urlkey"
	"strconv"
	"strings"
)

type Row struct {
	ID          string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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

func (f *FileStorage) Set(OriginalURL string) (string, error) {
	ShortURL := urlkey.GenerateSlug(OriginalURL)
	if ShortURL == "" {
		return "", fmt.Errorf("shortURL is empty")
	}
	rowID := strconv.Itoa(f.counter + 1)

	row := Row{
		ID:          rowID,
		ShortURL:    ShortURL,
		OriginalURL: OriginalURL,
	}
	// Write JSON entry
	err := f.encoder.Encode(row)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %s", err)
	}
	f.counter++
	return ShortURL, nil
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
