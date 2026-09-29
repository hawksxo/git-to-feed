package main

import (
	"fmt"
	"net/http"

	"github.com/hawksxo/git-to-feed/internal/platform/config"
	"github.com/hawksxo/git-to-feed/internal/platform/http"
)

func main() {
	cfg := config.LoadConfig()
	router := server.NewRouter(cfg)

	fmt.Println("Servidor iniciado en http://localhost:8080")

	err := http.ListenAndServe(cfg.Port, router)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
