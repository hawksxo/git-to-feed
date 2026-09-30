package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var ErrEmptyDatabaseURL = errors.New("database URL cannot be empty")

func NewPostgresDB(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, ErrEmptyDatabaseURL
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error abriendo conexion postgres: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error conectando a postgresql: %w", err)
	}

	return db, nil
}