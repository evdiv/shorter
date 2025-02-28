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
	return ""
}

func (storage *DBStorage) Get(key string) string {
	return ""
}

func (storage *DBStorage) Close() error {
	return storage.db.Close()
}
