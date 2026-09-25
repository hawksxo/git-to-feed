package server

import (
	"fmt"
	"net/http"

	"github.com/hawksxo/git-to-feed/internal/webhook"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	useCase := webhook.ProcessWebhookUseCase{}
	handler := webhook.Handler{ProcessWebhookUseCase: useCase}

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/webhooks/github", handler.HandleWebhook)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "UP", "message": "Git-To-Feed API is healthy"}`)
}