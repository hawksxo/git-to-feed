package storage

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(databaseURL string, migrationsPath string) error {
	if databaseURL == "" {
		return errors.New("database URL is empty, skipping auto-migrations")
	}
	if migrationsPath == "" {
		migrationsPath = "file://migrations"
	}

	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("error instantiating migration engine: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error executing SQL migrations: %w", err)
	}

	log.Println("✅ SQL auto-migrations executed successfully on Supabase/PostgreSQL")
	return nil
}
