package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DbStorage struct {
	db  sql.DB
	dsn string
}

func NewDbStorage(dsn string) (*DbStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}
	return &DbStorage{
		dsn: dsn,
		db:  *db,
	}, nil
}

func (storage *DbStorage) IsAvailable() bool {
	err := storage.db.Ping()
	if err != nil {
		fmt.Println("Database connection failed:", err)
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
