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
	_, err := storage.db.Exec("CREATE TABLE IF NOT EXISTS Links " +
		"(ID SERIAL PRIMARY KEY," +
		"ShortURL VARCHAR(128) NOT NULL UNIQUE," +
		"OriginalURL VARCHAR(512) NOT NULL," +
		"AddedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return fmt.Errorf("failed to create migration table: %s", err)
	}
	return nil
}

func (storage *DBStorage) IsAvailable() bool {
	err := storage.db.Ping()
	return err == nil
}

func (storage *DBStorage) Set(url string) string {
	key := urlkey.GenerateSlug(url)
	if key == "" {
		return ""
	}
	_, err := storage.db.Exec("INSERT INTO Links (ShortURL, OriginalURL) VALUES ($1, $2)", key, url)
	if err != nil {
		fmt.Errorf("failed to insert a new record to the Database: %s", err)
		return ""
	}
	return key
}

func (storage *DBStorage) Get(key string) string {
	row := storage.db.QueryRow("SELECT OriginalURL FROM Links WHERE ShortURL = $1", key)

	var url string
	err := row.Scan(&url)
	if err != nil {
		fmt.Errorf("failed to select a record from the Database: %s", err)
		return ""
	}
	return url
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}
