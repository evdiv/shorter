package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbStorage struct {
	db         *sql.DB
	connection string
}

func NewDbStorage(connection string) (*DbStorage, error) {
	db, err := sql.Open("pgx", connection)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", connection)
	}
	return &DbStorage{
		connection: connection,
		db:         db,
	}, nil
}

func (storage *DbStorage) IsAvailable() bool {
	err := storage.db.Ping()
	if err != nil {
		return false
	}
	return true
}

func (storage *DbStorage) Set(url string) string {
	return ""
}

func (storage *DbStorage) Get(key string) string {
	return ""
}

func (storage *DbStorage) Close() error {
	return storage.db.Close()
}
