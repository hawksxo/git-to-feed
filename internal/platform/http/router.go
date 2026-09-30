package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/platform"
	"github.com/hawksxo/git-to-feed/internal/platform/config"
	"github.com/hawksxo/git-to-feed/internal/publisher"
	"github.com/hawksxo/git-to-feed/internal/webhook"
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

	var linkedInClient publisher.LinkedInClient
	if cfg.LinkedInAccessToken != "" {
		realClient, err := publisher.NewHTTPLinkedInClient(cfg.LinkedInAccessToken, nil)
		if err == nil {
			linkedInClient = realClient
		} else {
			linkedInClient = publisher.NewMockLinkedInClient("urn:li:share:mock", false)
		}
	} else {
		linkedInClient = publisher.NewMockLinkedInClient("urn:li:share:mock", false)
	}

	publishUseCase, _ := publisher.NewPublishApprovedPostUseCase(linkedInClient, publisherRepo, cfg.LinkedInAuthorURN)

	orchestrator := platform.NewEventOrchestrator(approvalUseCase, postProcessor, publishUseCase)

	mux.HandleFunc("POST /api/v1/webhooks/github", func(w http.ResponseWriter, r *http.Request) {
		// Validacion de evento GitHub Webhook HMAC
		secret := cfg.GitHubWebhookSecret
		if secret == "" {
			secret = "default_dev_secret_git_to_feed"
		}
		signature := r.Header.Get("X-Hub-Signature-256")
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error al leer cuerpo de petición", http.StatusBadRequest)
			return
		}
		if signature != "" {
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write(bodyBytes)
			calculatedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
			if !hmac.Equal([]byte(signature), []byte(calculatedSignature)) {
				http.Error(w, "Firma inválida", http.StatusUnauthorized)
				return
			}
		}

		eventType := r.Header.Get("X-GitHub-Event")
		if eventType == "ping" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"OK","message":"Pong"}`))
			return
		}

		var payload webhook.GitHubPayload
		payload.EventType = eventType
		if len(bodyBytes) > 0 {
			_ = json.Unmarshal(bodyBytes, &payload)
		}

		approvalPost, err := orchestrator.ProcessGitHubPayload(r.Context(), payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"CREATED","uuid":"%s"}`, approvalPost.UUID)))
	})

	mux.HandleFunc("POST /api/v1/approvals/{uuid}/approve", func(w http.ResponseWriter, r *http.Request) {
		uuid := r.PathValue("uuid")
		pubRecord, err := orchestrator.ApproveAndPublish(r.Context(), uuid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"status":"PUBLISHED","share_urn":"%s"}`, pubRecord.LinkedInShareURN)))
	})

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, `{"status": "UP", "message": "Git-To-Feed API is healthy"}`)
}