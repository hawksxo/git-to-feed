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
		fmt.Println("🚀 Ejecutando auto-migraciones SQL en Supabase/PostgreSQL...")
		if err := storage.RunMigrations(cfg.DatabaseURL, ""); err != nil {
			fmt.Printf("⚠️ Advertencia en migraciones: %v\n", err)
		}
	}

	router := server.NewRouter(cfg)

	fmt.Println("Servidor iniciado en http://localhost:8080")

	err := http.ListenAndServe(cfg.Port, router)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
