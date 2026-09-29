package publisher_test

import (
	"context"
	"testing"
	"time"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/publisher"
)

func TestInMemoryPublisherRepository(t *testing.T) {
	ctx := context.Background()
	repo := publisher.NewInMemoryPublisherRepository()

	post, _ := publisher.NewPublishedPost("550e8400-e29b-41d4-a716-446655440000")

	t.Run("Debe guardar y recuperar por UUID", func(t *testing.T) {
		err := repo.Save(ctx, post)
		if err != nil {
			t.Fatalf("se esperaba error nil en Save, se obtuvo %v", err)
		}

		found, err := repo.FindByUUID(ctx, post.ID)
		if err != nil {
			t.Fatalf("se esperaba error nil en FindByUUID, se obtuvo %v", err)
		}
		if found.ID != post.ID {
			t.Errorf("se esperaba ID %s, se obtuvo %s", post.ID, found.ID)
		}
	})

	t.Run("Debe buscar por ApprovalUUID", func(t *testing.T) {
		found, err := repo.FindByApprovalUUID(ctx, post.ApprovalPostUUID)
		if err != nil {
			t.Fatalf("se esperaba error nil en FindByApprovalUUID, se obtuvo %v", err)
		}
		if found.ApprovalPostUUID != post.ApprovalPostUUID {
			t.Errorf("se esperaba ApprovalPostUUID %s, se obtuvo %s", post.ApprovalPostUUID, found.ApprovalPostUUID)
		}
	})

	t.Run("Debe retornar ErrNotFound si no existe", func(t *testing.T) {
		_, err := repo.FindByUUID(ctx, "non-existent-uuid")
		if err != publisher.ErrNotFound {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrNotFound, err)
		}
	})
}

func TestPublishApprovedPostUseCase(t *testing.T) {
	ctx := context.Background()
	authorURN := "urn:li:person:author123"

	t.Run("Debe publicar con exito un post en estado APPROVED", func(t *testing.T) {
		mockClient := publisher.NewMockLinkedInClient("urn:li:share:88888", false)
		repo := publisher.NewInMemoryPublisherRepository()
		useCase, _ := publisher.NewPublishApprovedPostUseCase(mockClient, repo, authorURN)

		approvedPost := &approval.ApprovalPost{
			UUID: "550e8400-e29b-41d4-a716-446655440000",
			GeneratedPost: pipeline.GeneratedPost{
				Content: "Lanzamos una nueva arquitectura!",
			},
			Status:    approval.StatusApproved,
			CreatedAt: time.Now().UTC(),
		}

		pubRecord, err := useCase.Execute(ctx, approvedPost)
		if err != nil {
			t.Fatalf("se esperaba error nil, se obtuvo %v", err)
		}
		if pubRecord.Status != publisher.PublishStatusPublished {
			t.Errorf("se esperaba estado PUBLISHED, se obtuvo %s", pubRecord.Status)
		}
		if pubRecord.LinkedInShareURN != "urn:li:share:88888" {
			t.Errorf("se esperaba URN urn:li:share:88888, se obtuvo %s", pubRecord.LinkedInShareURN)
		}
	})

	t.Run("Debe rechazar publicacion si la aprobacion no esta en estado APPROVED", func(t *testing.T) {
		mockClient := publisher.NewMockLinkedInClient("", false)
		repo := publisher.NewInMemoryPublisherRepository()
		useCase, _ := publisher.NewPublishApprovedPostUseCase(mockClient, repo, authorURN)

		pendingPost := &approval.ApprovalPost{
			UUID: "550e8400-e29b-41d4-a716-446655440000",
			GeneratedPost: pipeline.GeneratedPost{
				Content: "Post pendiente",
			},
			Status:    approval.StatusPending,
			CreatedAt: time.Now().UTC(),
		}

		pubRecord, err := useCase.Execute(ctx, pendingPost)
		if err != publisher.ErrPostNotApproved {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrPostNotApproved, err)
		}
		if pubRecord != nil {
			t.Errorf("se esperaba pubRecord nil, se obtuvo %v", pubRecord)
		}
	})

	t.Run("Debe marcar como FAILED si la API de LinkedIn retorna error", func(t *testing.T) {
		mockClient := publisher.NewMockLinkedInClient("", true)
		repo := publisher.NewInMemoryPublisherRepository()
		useCase, _ := publisher.NewPublishApprovedPostUseCase(mockClient, repo, authorURN)

		approvedPost := &approval.ApprovalPost{
			UUID: "550e8400-e29b-41d4-a716-446655440000",
			GeneratedPost: pipeline.GeneratedPost{
				Content: "Post con fallo en API",
			},
			Status:    approval.StatusApproved,
			CreatedAt: time.Now().UTC(),
		}

		pubRecord, err := useCase.Execute(ctx, approvedPost)
		if err == nil {
			t.Fatalf("se esperaba error al fallar LinkedIn API, se obtuvo nil")
		}
		if pubRecord == nil {
			t.Fatalf("se esperaba pubRecord no nil para inspeccionar estado FAILED")
		}
		if pubRecord.Status != publisher.PublishStatusFailed {
			t.Errorf("se esperaba estado FAILED, se obtuvo %s", pubRecord.Status)
		}
	})
}
