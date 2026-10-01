package main

import (
	"fmt"
	"net/http"

	"github.com/hawksxo/git-to-feed/internal/platform/config"
	"github.com/hawksxo/git-to-feed/internal/platform/http"
	"github.com/hawksxo/git-to-feed/internal/platform/storage"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.DatabaseURL != "" {
		fmt.Println("🚀 Running SQL auto-migrations on Supabase/PostgreSQL...")
		if err := storage.RunMigrations(cfg.DatabaseURL, ""); err != nil {
			fmt.Printf("⚠️ Migration warning: %v\n", err)
		}
	}

	router := server.NewRouter(cfg)

	fmt.Println("Server started on http://localhost:8080")

	err := http.ListenAndServe(cfg.Port, router)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
