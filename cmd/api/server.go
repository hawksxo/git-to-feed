package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/config"
	"github.com/hawksxo/git-to-feed/internal/generator"
	"github.com/hawksxo/git-to-feed/internal/publisher"
	"github.com/hawksxo/git-to-feed/internal/storage"
	"github.com/hawksxo/git-to-feed/internal/webhook"
)

type discordApproverAdapter struct {
	approvalSvc  *approval.Service
	publisherSvc *publisher.Service
}

func (a *discordApproverAdapter) ApproveAndPublish(ctx context.Context, uuid string) (string, error) {
	approvedPost, err := a.approvalSvc.Approve(ctx, uuid)
	if err != nil {
		return "", fmt.Errorf("error al aprobar borrador: %w", err)
	}

	pubRecord, err := a.publisherSvc.Publish(ctx, approvedPost.UUID, approvedPost.Content)
	if err != nil {
		return "", fmt.Errorf("error al publicar post aprobado: %w", err)
	}

	return pubRecord.LinkedInShareURN, nil
}

func (a *discordApproverAdapter) DiscardDraft(ctx context.Context, uuid, reason string) error {
	_, err := a.approvalSvc.Reject(ctx, uuid, reason)
	if err != nil {
		return fmt.Errorf("error al descartar borrador: %w", err)
	}
	return nil
}

func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Salud del servicio
	mux.HandleFunc("GET /health", healthHandler)

	// Persistencia PostgreSQL / Supabase o InMemory
	var approvalStore approval.Store
	var publisherStore publisher.Store

	if cfg.DatabaseURL != "" {
		db, err := storage.NewPostgresDB(cfg.DatabaseURL)
		if err == nil {
			approvalStore = approval.NewPostgresStore(db)
			publisherStore = publisher.NewPostgresStore(db)
		} else {
			log.Printf("Aviso: No se pudo conectar a la base de datos Postgres: %v. Usando almacenamiento en memoria.", err)
			approvalStore = approval.NewInMemoryStore()
			publisherStore = publisher.NewInMemoryStore()
		}
	} else {
		approvalStore = approval.NewInMemoryStore()
		publisherStore = publisher.NewInMemoryStore()
	}

	approvalSvc := approval.NewService(approvalStore)

	// Cliente LinkedIn
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

	authorURN := cfg.LinkedInAuthorURN
	if authorURN == "" {
		authorURN = "urn:li:person:default-author"
	}

	publisherSvc, err := publisher.NewService(linkedInClient, publisherStore, authorURN)
	if err != nil {
		log.Fatalf("Error al inicializar el servicio de publicacion: %v", err)
	}

	// Cliente Gemini LLM y Generador
	var llmClient generator.LLMClient
	if cfg.GeminiAPIKey != "" {
		gClient, err := generator.NewGeminiLLMClient(cfg.GeminiAPIKey, cfg.GeminiModelName)
		if err == nil {
			llmClient = gClient
		} else {
			log.Printf("Aviso: No se pudo conectar a Gemini: %v. Usando generador simulado.", err)
			llmClient = generator.NewMockLLMClient("Post simulado generado para pruebas.")
		}
	} else {
		llmClient = generator.NewMockLLMClient("Post simulado generado sin API key.")
	}
	genSvc := generator.New(llmClient)

	// Bot de Discord (Opcional si estan las credenciales configuradas)
	var discordBot *approval.Bot
	if cfg.DiscordBotToken != "" && cfg.DiscordChannelID != "" {
		approverAdapter := &discordApproverAdapter{
			approvalSvc:  approvalSvc,
			publisherSvc: publisherSvc,
		}
		bot, err := approval.NewBot(cfg.DiscordBotToken, approval.ChannelConfig{
			DefaultChannelID:  cfg.DiscordChannelID,
			ReleasesChannelID: cfg.DiscordChannelReleasesID,
			FeaturesChannelID: cfg.DiscordChannelFeaturesID,
			AuditChannelID:    cfg.DiscordChannelAuditID,
		}, approverAdapter)
		if err == nil {
			discordBot = bot
		} else {
			log.Printf("Aviso: No se pudo inicializar el bot de Discord: %v", err)
		}
	}

	// Endpoint Webhook GitHub
	mux.HandleFunc("POST /api/v1/webhooks/github", func(w http.ResponseWriter, r *http.Request) {
		payload, err := webhook.ParseAndValidate(r, cfg.GitHubWebhookSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if payload.EventType == "ping" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"OK","message":"pong"}`))
			return
		}

		// Filtrado de eventos
		var archetype generator.Archetype
		var gitData generator.GitData

		switch payload.EventType {
		case "pull_request":
			shouldProcess, reason := webhook.ShouldProcessPullRequest(payload.PullRequest)
			if !shouldProcess {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"status":"IGNORED","message":"%s"}`, reason)))
				return
			}
			archetype = generator.ArchetypeFeature
			gitData = generator.GitData{
				EventType:   payload.EventType,
				Title:       payload.PullRequest.Title,
				Description: payload.PullRequest.Body,
				Author:      payload.Sender.User,
				Repository:  payload.Repository.FullName,
				URL:         payload.PullRequest.HTMLURL,
			}
		case "release":
			if payload.Action != "published" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"IGNORED","message":"accion de release no es published"}`))
				return
			}
			archetype = generator.ArchetypeRelease
			gitData = generator.GitData{
				EventType:   payload.EventType,
				Tag:         payload.Release.TagName,
				Title:       payload.Release.Name,
				Description: payload.Release.Body,
				Author:      payload.Sender.User,
				Repository:  payload.Repository.FullName,
				URL:         payload.Release.HTMLURL,
			}
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"IGNORED","message":"tipo de evento no procesable"}`))
			return
		}

		richCtx := generator.RichContext{
			GitData:         gitData,
			ProjectIdentity: payload.Repository.Description,
			CreatedAt:       time.Now(),
		}

		genPost, err := genSvc.Generate(r.Context(), richCtx, archetype)
		if err != nil {
			http.Error(w, fmt.Sprintf("error al generar contenido: %v", err), http.StatusInternalServerError)
			return
		}

		idempotencyKey := fmt.Sprintf("%s:%s", payload.EventType, gitData.URL)
		appPost, err := approvalSvc.SubmitForApproval(r.Context(), genPost.Content, string(archetype), idempotencyKey)
		if err != nil {
			http.Error(w, fmt.Sprintf("error al someter borrador a aprobacion: %v", err), http.StatusInternalServerError)
			return
		}

		if discordBot != nil {
			if err := discordBot.SendApprovalNotification(appPost); err != nil {
				log.Printf("Aviso: Error enviando notificacion a Discord: %v", err)
			}
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fmt.Appendf(nil, `{"status":"CREATED","uuid":"%s"}`, appPost.UUID))
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
			"api_router":      "OPERATIONAL",
			"webhook_engine":  "OPERATIONAL",
			"pipeline_gemini": "OPERATIONAL",
			"discord_bot":     "OPERATIONAL",
			"linkedin_engine": "OPERATIONAL",
		},
	}
	_ = json.NewEncoder(w).Encode(response)
}
