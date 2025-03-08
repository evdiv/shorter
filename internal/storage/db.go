package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"shorter/internal/models"
	"shorter/internal/urlkey"
)

type DBStorage struct {
	db         *sql.DB
	connection string
}

func NewDBStorage(connection string) (*DBStorage, error) {
	db, err := sql.Open("pgx", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", connection)
	}
	return &DBStorage{
		connection: connection,
		db:         db,
	}, nil
}

// Migrate - creates the tables if they don't exist
func (storage *DBStorage) Migrate() error {
	_, err := storage.db.Exec(`CREATE TABLE IF NOT EXISTS Links 
		(ID SERIAL PRIMARY KEY, CorrelationID VARCHAR(128) NULL,
    	ShortURL VARCHAR(128) NOT NULL,
    	OriginalURL VARCHAR(512) NOT NULL UNIQUE,
    	AddedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		return fmt.Errorf("failed to create migration table: %s", err)
	}
	return nil
}

func (storage *DBStorage) IsAvailable() bool {
	err := storage.db.Ping()
	return err == nil
}

func (storage *DBStorage) Set(OriginalURL string) (string, error) {
	ShortURL := urlkey.GenerateSlug(OriginalURL)
	if ShortURL == "" {
		return "", fmt.Errorf("the short url is empty")
	}

	result, err := storage.db.Exec(`INSERT INTO Links (ShortURL, OriginalURL) 
											VALUES ($1, $2) 
											ON CONFLICT (OriginalURL)
											DO NOTHING`, ShortURL, OriginalURL)
	if err != nil {
		return "", NewStorageError("failed to insert", OriginalURL, ShortURL, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get affected rows: %s", err)
	}
	if rowsAffected == 0 {
		return ShortURL, NewStorageError("already exists", OriginalURL, ShortURL, err)
	}
	return ShortURL, nil
}

func (storage *DBStorage) SetBatch(jReqBatch []models.JSONReq) ([]models.JSONRes, error) {
	jResBatch := []models.JSONRes{}

	stmt, err := storage.db.Prepare(`INSERT INTO Links (ShortURL, OriginalURL) VALUES ($1, $2)`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %s", err)
	}
	defer stmt.Close()

	for _, el := range jReqBatch {
		ShortURL := urlkey.GenerateSlug(el.OriginalURL)
		if ShortURL == "" {
			return nil, fmt.Errorf("the short url for Original Url: %s is empty", el.OriginalURL)
		}
		_, err := stmt.Exec(ShortURL, el.OriginalURL)
		if err != nil {
			return nil, NewStorageError("failed to insert", ShortURL, el.OriginalURL, err)
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

func (storage *DBStorage) Get(ShortURL string) (string, error) {
	var OriginalURL string
	row := storage.db.QueryRow("SELECT OriginalURL FROM Links WHERE ShortURL = $1", ShortURL)

	err := row.Scan(&OriginalURL)
	if err != nil {
		return "", NewStorageError("failed to select", OriginalURL, ShortURL, err)
	}
	return OriginalURL, nil
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}
