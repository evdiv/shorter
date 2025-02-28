package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func (storage *DBStorage) IsAvailable() bool {
	err := storage.db.Ping()
	return err == nil
}

func (storage *DBStorage) Set(url string) string {
	return ""
}

func (storage *DBStorage) Get(key string) string {
	return ""
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}
