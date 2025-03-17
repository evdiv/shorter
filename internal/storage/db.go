package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"shorter/internal/config"
	"shorter/internal/models"
	"shorter/internal/urlkey"
	"time"
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

	// Set connection pool settings
	db.SetMaxOpenConns(25)                 // Allow up to 25 open connections
	db.SetMaxIdleConns(5)                  // Keep up to 5 idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Recreate connections every 5 minutes

	return &DBStorage{
		connection: connection,
		db:         db,
	}, nil
}

// Migrate - creates the tables if they don't exist
func (storage *DBStorage) Migrate() error {

	query := `CREATE TABLE IF NOT EXISTS Links (
        ID SERIAL PRIMARY KEY,
        UserID VARCHAR(128) NULL,
        CorrelationID VARCHAR(128) NULL,
        ShortURL VARCHAR(128) NOT NULL,
        OriginalURL VARCHAR(512) NOT NULL UNIQUE,
        AddedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	_, err := storage.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migration table: %s", err)
	}
	return nil
}

func (storage *DBStorage) IsAvailable() bool {
	err := storage.db.Ping()
	return err == nil
}

func (storage *DBStorage) Set(OriginalURL string, userID string) (string, error) {

	query := `INSERT INTO Links (ShortURL, OriginalURL, UserID)
		VALUES ($1, $2, $3)
		ON CONFLICT (OriginalURL)
		DO NOTHING`

	urlKey := urlkey.GenerateSlug(OriginalURL)
	if urlKey == "" {
		return "", fmt.Errorf("the short url is empty")
	}

	fmt.Printf("urlKey: %s\n", urlKey)
	fmt.Printf("OriginalURL: %s\n", OriginalURL)
	fmt.Printf("UserID: %s\n", userID)

	result, err := storage.db.Exec(query, urlKey, OriginalURL, userID)
	if err != nil {
		return "", NewStorageError("failed to insert", OriginalURL, urlKey, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get affected rows: %s", err)
	}
	if rowsAffected == 0 {
		return urlKey, NewStorageError("already exists", OriginalURL, urlKey, err)
	}
	return urlKey, nil
}

func (storage *DBStorage) SetBatch(jReqBatch []models.JSONReq, userID string) ([]models.JSONRes, error) {

	query := `INSERT INTO Links (ShortURL, OriginalURL, UserID)
		VALUES ($1, $2, $3)
		ON CONFLICT (OriginalURL)
		DO NOTHING`

	jResBatch := []models.JSONRes{}

	stmt, err := storage.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %s", err)
	}
	defer stmt.Close()

	for _, el := range jReqBatch {
		urlKey := urlkey.GenerateSlug(el.OriginalURL)
		if urlKey == "" {
			return nil, fmt.Errorf("the urlKey for Original Url: %s is empty", el.OriginalURL)
		}
		_, err := stmt.Exec(urlKey, el.OriginalURL, userID)
		if err != nil {
			return nil, NewStorageError("failed to insert", urlKey, el.OriginalURL, err)
		}
		row := models.JSONRes{
			CorrID:      el.CorrID,
			ShortURL:    config.AppConfig.ResultHost + "/" + urlKey,
			OriginalURL: el.OriginalURL,
		}
		jResBatch = append(jResBatch, row)
	}

	return jResBatch, nil
}

func (storage *DBStorage) Get(ShortURL string) (string, error) {
	query := `SELECT OriginalURL FROM Links WHERE ShortURL = $1`

	var OriginalURL string
	row := storage.db.QueryRow(query, ShortURL)

	err := row.Scan(&OriginalURL)
	if err != nil {
		return "", NewStorageError("failed to select", OriginalURL, ShortURL, err)
	}
	return OriginalURL, nil
}

func (storage *DBStorage) GetUserURLs(userID string) ([]models.JSONUserRes, error) {
	jResBatch := make([]models.JSONUserRes, 0)

	log.Println("GetUserURLs from a database.")

	query := `SELECT ShortURL, OriginalURL FROM Links WHERE UserID = $1`
	rows, err := storage.db.Query(query, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve links for user: %s", userID)
	}
	defer rows.Close()

	for rows.Next() {
		var row models.JSONUserRes

		if err := rows.Scan(&row.ShortURL, &row.OriginalURL); err != nil {
			return nil, fmt.Errorf("failed to scan row: %s", err)
		}
		row.ShortURL = config.AppConfig.ResultHost + "/" + row.ShortURL

		log.Println("Before appending a row")

		jResBatch = append(jResBatch, row)
	}

	// **Check for errors during iteration**
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return jResBatch, nil
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}
