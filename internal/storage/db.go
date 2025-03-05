package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
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
		return "", fmt.Errorf("failed to insert link: %s", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get affected rows: %s", err)
	}
	if rowsAffected == 0 {
		return ShortURL, fmt.Errorf("the original url: %s already exists", OriginalURL)
	}
	return ShortURL, nil
}

func (storage *DBStorage) Get(ShortURL string) (string, error) {
	var OriginalURL string
	row := storage.db.QueryRow("SELECT OriginalURL FROM Links WHERE ShortURL = $1", ShortURL)

	err := row.Scan(&OriginalURL)
	if err != nil {
		return "", fmt.Errorf("failed to get original url: %s", err)
	}
	return OriginalURL, nil
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}
