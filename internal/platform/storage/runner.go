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
		return fmt.Errorf("error instanciando motor de migraciones: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error ejecutando migraciones SQL: %w", err)
	}

	log.Println("✅ Auto-migraciones SQL ejecutadas con éxito en Supabase/PostgreSQL")
	return nil
}
