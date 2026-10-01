package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/platform"
	"github.com/hawksxo/git-to-feed/internal/platform/config"
	"github.com/hawksxo/git-to-feed/internal/platform/discord"
	"github.com/hawksxo/git-to-feed/internal/platform/storage"
	"github.com/hawksxo/git-to-feed/internal/publisher"
	"github.com/hawksxo/git-to-feed/internal/webhook"
)

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /health/api-router", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "OPERATIONAL",
			"component": "api_router",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("GET /health/webhook-engine", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "OPERATIONAL",
			"component": "webhook_engine",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("GET /health/pipeline-gemini", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := "OPERATIONAL"
		if cfg.GeminiAPIKey == "" {
			status = "DEGRADED"
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    status,
			"component": "pipeline_gemini",
			"model":     cfg.GeminiModelName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("GET /health/discord-bot", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := "OPERATIONAL"
		if cfg.DiscordBotToken == "" || cfg.DiscordChannelID == "" {
			status = "DEGRADED"
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    status,
			"component": "discord_bot",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("GET /health/linkedin-engine", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := "OPERATIONAL"
		if cfg.LinkedInAccessToken == "" || cfg.LinkedInAuthorURN == "" {
			status = "DEGRADED"
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    status,
			"component": "linkedin_engine",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})
	mux.HandleFunc("GET /health/database", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		status := "OPERATIONAL"
		message := "Supabase PostgreSQL connected"
		if cfg.DatabaseURL == "" {
			status = "DEGRADED"
			message = "Running on InMemory repository fallback"
		} else {
			db, err := storage.NewPostgresDB(cfg.DatabaseURL)
			if err != nil || db.Ping() != nil {
				status = "OUTAGE"
				message = "Database connection failed"
				w.WriteHeader(http.StatusServiceUnavailable)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"status":    status,
					"component": "database_supabase",
					"message":   message,
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				})
				return
			}
			db.Close()
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    status,
			"component": "database_supabase",
			"message":   message,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Approval & Publisher Persistence (PostgreSQL / Supabase if DATABASE_URL is set, else InMemory)
	var approvalRepo approval.ApprovalRepository
	var publisherRepo publisher.PublisherRepository

	if cfg.DatabaseURL != "" {
		db, err := storage.NewPostgresDB(cfg.DatabaseURL)
		if err == nil {
			approvalRepo = approval.NewPostgresApprovalRepository(db)
			publisherRepo = publisher.NewPostgresPublisherRepository(db)
		} else {
			approvalRepo = approval.NewInMemoryApprovalRepository()
			publisherRepo = publisher.NewInMemoryPublisherRepository()
		}
	} else {
		approvalRepo = approval.NewInMemoryApprovalRepository()
		publisherRepo = publisher.NewInMemoryPublisherRepository()
	}

	approvalUseCase := approval.NewApprovalUseCase(approvalRepo)
	approvalHandler := approval.NewHandler(approvalUseCase)

	mux.HandleFunc("GET /api/v1/approvals/pending", approvalHandler.HandleListPending)
	mux.HandleFunc("POST /api/v1/approvals/{uuid}/reject", approvalHandler.HandleReject)
	mux.HandleFunc("PUT /api/v1/approvals/{uuid}/edit", approvalHandler.HandleEditAndApprove)

	// Few-Shot Admin Endpoints
	var fewShotRepo pipeline.FewShotRepository
	if cfg.DatabaseURL != "" {
		if db, err := storage.NewPostgresDB(cfg.DatabaseURL); err == nil {
			fewShotRepo = pipeline.NewPostgresFewShotRepository(db)
		} else {
			fewShotRepo = pipeline.NewInMemoryFewShotRepository()
		}
	} else {
		fewShotRepo = pipeline.NewInMemoryFewShotRepository()
	}

	adminHandler := pipeline.NewAdminHandler(fewShotRepo)
	mux.HandleFunc("GET /api/v1/examples", adminHandler.HandleListExamples)
	mux.HandleFunc("POST /api/v1/examples", adminHandler.HandleCreateExample)

	// Publisher module routes
	postProcessor := pipeline.NewPostProcessor()

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

	// Optional initialization of Discord Bot for interactive approvals
	if cfg.DiscordBotToken != "" && cfg.DiscordChannelID != "" {
		bot, err := discord.NewBot(cfg.DiscordBotToken, cfg.DiscordChannelID, orchestrator)
		if err == nil {
			orchestrator.SetNotifier(bot)
		} else {
			fmt.Printf("⚠️ Warning: Failed to connect Discord Bot: %v\n", err)
		}
	}

	mux.HandleFunc("POST /api/v1/webhooks/github", func(w http.ResponseWriter, r *http.Request) {
		// GitHub Webhook HMAC signature validation
		secret := cfg.GitHubWebhookSecret
		if secret == "" {
			secret = "default_dev_secret_git_to_feed"
		}
		signature := r.Header.Get("X-Hub-Signature-256")
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}
		if signature != "" {
			mac := hmac.New(sha256.New, []byte(secret))
			mac.Write(bodyBytes)
			calculatedSignature := "sha256=" + hex.EncodeToString(mac.Sum(nil))
			if !hmac.Equal([]byte(signature), []byte(calculatedSignature)) {
				http.Error(w, "Invalid signature", http.StatusUnauthorized)
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

		// STEP 1: Strict release action filtering (only process published or released events)
		if eventType == "release" && payload.Action != "published" && payload.Action != "released" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"status":"OK","message":"Release action '%s' ignored"}`, payload.Action)))
			return
		}

		// STEP 2: Strict pull request action filtering (only process merged pull requests)
		if eventType == "pull_request" && (payload.Action != "closed" || !payload.PullRequest.Merged) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"status":"OK","message":"PullRequest action '%s' (merged=%t) ignored"}`, payload.Action, payload.PullRequest.Merged)))
			return
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
	response := map[string]any{
		"status":    "OPERATIONAL",
		"service":   "git-to-feed",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"components": map[string]string{
			"api_router":       "OPERATIONAL",
			"webhook_engine":   "OPERATIONAL",
			"pipeline_gemini":  "OPERATIONAL",
			"discord_bot":      "OPERATIONAL",
			"linkedin_engine": "OPERATIONAL",
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}