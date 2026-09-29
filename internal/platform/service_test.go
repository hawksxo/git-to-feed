package platform_test

import (
	"context"
	"testing"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/platform"
	"github.com/hawksxo/git-to-feed/internal/publisher"
)

func TestEventOrchestrator_FullFlow(t *testing.T) {
	ctx := context.Background()

	// 1. Instanciar Repositorios y Servicios
	approvalRepo := approval.NewInMemoryApprovalRepository()
	approvalUseCase := approval.NewApprovalUseCase(approvalRepo)

	postProcessor := pipeline.NewPostProcessor()

	publisherRepo := publisher.NewInMemoryPublisherRepository()
	mockLinkedInClient := publisher.NewMockLinkedInClient("urn:li:share:123456", false)
	publishUseCase, _ := publisher.NewPublishApprovedPostUseCase(mockLinkedInClient, publisherRepo, "urn:li:person:author")

	// 2. Instanciar Orquestador
	orchestrator := platform.NewEventOrchestrator(approvalUseCase, postProcessor, publishUseCase)

	// 3. Probar Ingesta y Generacion (ProcessGitHubEvent)
	approvalPost, err := orchestrator.ProcessGitHubEvent(ctx, "Lanzamos una game-changer feature!", pipeline.ArchetypeFeature)
	if err != nil {
		t.Fatalf("error insperado en ProcessGitHubEvent: %v", err)
	}
	if approvalPost == nil || approvalPost.UUID == "" {
		t.Fatalf("se esperaba un approvalPost valido")
	}
	if approvalPost.Status != approval.StatusPending {
		t.Errorf("se esperaba estado PENDING, se obtuvo %s", approvalPost.Status)
	}

	// 4. Probar Aprobacion y Publicacion (ApproveAndPublish)
	publishedRecord, err := orchestrator.ApproveAndPublish(ctx, approvalPost.UUID)
	if err != nil {
		t.Fatalf("error inesperado en ApproveAndPublish: %v", err)
	}
	if publishedRecord == nil {
		t.Fatalf("se esperaba publishedRecord no nil")
	}
	if publishedRecord.Status != publisher.PublishStatusPublished {
		t.Errorf("se esperaba estado PUBLISHED, se obtuvo %s", publishedRecord.Status)
	}
	if publishedRecord.LinkedInShareURN != "urn:li:share:123456" {
		t.Errorf("se esperaba URN urn:li:share:123456, se obtuvo %s", publishedRecord.LinkedInShareURN)
	}
}
