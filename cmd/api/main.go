package main


import (
	"fmt"
	"net/http"

	"github.com/hawksxo/git-to-feed/internal/platform/http"
)

func main() {
	router := server.NewRouter()

	fmt.Println("Servidor iniciado en http://localhost:8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %v\n", err)
	}
}
