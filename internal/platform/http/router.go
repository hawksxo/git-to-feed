package server

import (
	"fmt"
	"net/http"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/platform"
	"github.com/hawksxo/git-to-feed/internal/platform/config"
	"github.com/hawksxo/git-to-feed/internal/publisher"
)

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	// Approval module routes (Human-in-the-Loop)
	approvalRepo := approval.NewInMemoryApprovalRepository()
	approvalUseCase := approval.NewApprovalUseCase(approvalRepo)
	approvalHandler := approval.NewHandler(approvalUseCase)

	mux.HandleFunc("GET /api/v1/approvals/pending", approvalHandler.HandleListPending)
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/reject", approvalHandler.HandleReject)
	mux.HandleFunc("PUT /api/v1/approvals/{uuid}/edit", approvalHandler.HandleEditAndApprove)

	// Publisher module routes
	postProcessor := pipeline.NewPostProcessor()
	publisherRepo := publisher.NewInMemoryPublisherRepository()
	linkedInClient := publisher.NewMockLinkedInClient("urn:li:share:mock", false)
	publishUseCase, _ := publisher.NewPublishApprovedPostUseCase(linkedInClient, publisherRepo, cfg.LinkedInAuthorURN)

	orchestrator := platform.NewEventOrchestrator(approvalUseCase, postProcessor, publishUseCase)
	mux.HandleFunc("POST /api/v1/webhooks/github", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"status":"RECEIVED"`))
	})
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/approve", func(w http.ResponseWriter, r *http.Request) {
		uuid := r.PathValue("uuid")
		_, err := orchestrator.ApproveAndPublish(r.Context(), uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"PUBLISHED"}`))
	})

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "UP", "message": "Git-To-Feed API is healthy"}`)
}