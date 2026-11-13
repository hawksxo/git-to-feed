package main_test

import (
	"context"
	"testing"
	"time"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/config"
	"github.com/hawksxo/git-to-feed/internal/generator"
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

	// 1. Verificacion del Pipeline Linear: Webhook Data -> Generator -> Approval -> Publisher
	t.Run("Flujo E2E: Ingesta -> Generacion -> Aprobacion Pendiente -> Publicacion en LinkedIn", func(t *testing.T) {
		// Mock LLM
		mockLLM := generator.NewMockLLMClient("Post generado para la comunidad tech en LinkedIn.")
		genSvc := generator.New(mockLLM)

		// Generar post a partir del payload con RichContext
		richCtx := generator.RichContext{
			GitData: generator.GitData{
				EventType:   "pull_request",
				Title:       "feat: nueva arquitectura limpia v2",
				Description: "Migración completa a package-by-feature",
				Author:      "octocat",
				Repository:  "hawksxo/git-to-feed",
				URL:         "https://github.com/hawksxo/git-to-feed/pull/80",
			},
			ProjectIdentity: "Microservicio event-driven en Go",
			CreatedAt:       time.Now(),
		}

		genPost, err := genSvc.Generate(ctx, richCtx, generator.ArchetypeFeature)
		if err != nil {
			t.Fatalf("error en generacion de post: %v", err)
		}

		if genPost.Content == "" {
			t.Fatal("el contenido generado no deberia estar vacio")
		}

		// Aprobacion
		approvalStore := approval.NewInMemoryStore()
		approvalSvc := approval.NewService(approvalStore)

		pendingPost, err := approvalSvc.SubmitForApproval(ctx, genPost.Content, string(generator.ArchetypeFeature), "pr:80")
		if err != nil {
			t.Fatalf("error al someter post a aprobacion: %v", err)
		}
		if pendingPost.Status != approval.StatusPending {
			t.Fatalf("se esperaba estado pendiente, se obtuvo: %s", pendingPost.Status)
		}

		// Publicador
		pubStore := publisher.NewInMemoryStore()
		mockLinkedIn := publisher.NewMockLinkedInClient("urn:li:share:e2e-success", false)
		pubSvc, err := publisher.NewService(mockLinkedIn, pubStore, cfg.LinkedInAuthorURN)
		if err != nil {
			t.Fatalf("error al inicializar servicio de publicacion: %v", err)
		}

		// Simulamos aprobacion del usuario (desde Discord bot o approver)
		approvedPost, err := approvalSvc.Approve(ctx, pendingPost.UUID)
		if err != nil {
			t.Fatalf("error al aprobar post: %v", err)
		}

		pubRecord, err := pubSvc.Publish(ctx, approvedPost.UUID, approvedPost.Content)
		if err != nil {
			t.Fatalf("error al publicar post: %v", err)
		}

		if pubRecord.Status != publisher.PublishStatusPublished {
			t.Fatalf("se esperaba estado publicado, se obtuvo: %s", pubRecord.Status)
		}
		if pubRecord.LinkedInShareURN != "urn:li:share:e2e-success" {
			t.Fatalf("se esperaba URN de LinkedIn valido, se obtuvo: %s", pubRecord.LinkedInShareURN)
		}
	})
}
