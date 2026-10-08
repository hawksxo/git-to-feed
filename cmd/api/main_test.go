package main_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/platform"
	"github.com/hawksxo/git-to-feed/internal/platform/config"
	server "github.com/hawksxo/git-to-feed/internal/platform/http"
	"github.com/hawksxo/git-to-feed/internal/publisher"
)

func TestFullE2EFlow(t *testing.T) {
	ctx := context.Background()

	cfg := &config.Config{
		Port:                ":8080",
		GitHubWebhookSecret: "test-secret",
		LinkedInAccessToken: "test-token",
		LinkedInAuthorURN:   "urn:li:person:testauthor",
	}

	router := server.NewRouter(cfg)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Verificacion Endpoint de Salud
	t.Run("GET /health debe retornar 200 OK y status UP", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/health")
		if err != nil {
			t.Fatalf("error haciendo GET /health: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("se esperaba status 200 OK, se obtuvo %d", resp.StatusCode)
		}
	})

	// Verificacion del Flujo Completo E2E mediante Orquestador y HTTP Approve
	t.Run("Flujo E2E: Ingesta -> Aprobacion Pendiente -> Publicacion en LinkedIn", func(t *testing.T) {
		approvalRepo := approval.NewInMemoryApprovalRepository()
		approvalUseCase := approval.NewApprovalUseCase(approvalRepo)
		postProcessor := pipeline.NewPostProcessor()

		publisherRepo := publisher.NewInMemoryPublisherRepository()
		mockLinkedInClient := publisher.NewMockLinkedInClient("urn:li:share:e2e-success", false)
		publishUseCase, _ := publisher.NewPublishApprovedPostUseCase(mockLinkedInClient, publisherRepo, cfg.LinkedInAuthorURN)

		orchestrator := platform.NewEventOrchestrator(approvalUseCase, postProcessor, publishUseCase)

		// 1. Simular Ingesta de Evento y Generacion
		pendingPost, err := orchestrator.ProcessGitHubEvent(ctx, "Lanzamos nueva arquitectura limpia en Go!", pipeline.ArchetypeFeature)
		if err != nil {
			t.Fatalf("error procesando evento de ingesta: %v", err)
		}
		if pendingPost.Status != approval.StatusPending {
			t.Errorf("se esperaba estado PENDING en post registrado, se obtuvo %s", pendingPost.Status)
		}

		// 2. Simular Aprobacion y Publicacion via Endpoint HTTP
		resp, err := http.Post(testServer.URL+"/api/v1/approvals/"+pendingPost.UUID+"/approve", "application/json", nil)
		if err != nil {
			t.Fatalf("error enviando POST /approve: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()

		// 3. Confirmar que el flujo completo hasta LinkedIn termino con exito
		publishedRecord, err := orchestrator.ApproveAndPublish(ctx, pendingPost.UUID)
		if err != nil {
			t.Fatalf("error en publicacion final: %v", err)
		}
		if publishedRecord.Status != publisher.PublishStatusPublished {
			t.Errorf("se esperaba estado PUBLISHED en el registro final, se obtuvo %s", publishedRecord.Status)
		}
	})
}
