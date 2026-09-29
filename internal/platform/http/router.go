package server

import (
	"fmt"
	"net/http"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/webhook"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Webhook module routes
	webhookUseCase := webhook.ProcessWebhookUseCase{}
	webhookHandler := webhook.Handler{ProcessWebhookUseCase: webhookUseCase}

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /api/v1/webhooks/github", webhookHandler.HandleWebhook)

	// Approval module routes (Human-in-the-Loop)
	approvalRepo := approval.NewInMemoryApprovalRepository()
	approvalUseCase := approval.NewApprovalUseCase(approvalRepo)
	approvalHandler := approval.NewHandler(approvalUseCase)

	mux.HandleFunc("GET /api/v1/approvals/pending", approvalHandler.HandleListPending)
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/approve", approvalHandler.HandleApprove)
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/reject", approvalHandler.HandleReject)
	mux.HandleFunc("PUT /api/v1/approvals/{uuid}/edit", approvalHandler.HandleEditAndApprove)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "UP", "message": "Git-To-Feed API is healthy"}`)
}