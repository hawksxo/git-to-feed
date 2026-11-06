package platform_test

import (
	"context"
	"testing"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/platform"
	"github.com/hawksxo/git-to-feed/internal/platform/github"
	"github.com/hawksxo/git-to-feed/internal/publisher"
	"github.com/hawksxo/git-to-feed/internal/webhook"
)

type mockComparator struct{}

func (m *mockComparator) CompareCommits(ctx context.Context, owner, repo, base, head string) ([]github.Commit, error) {
	return []github.Commit{}, nil
}

func TestEventOrchestrator_FullFlow(t *testing.T) {
	ctx := context.Background()

	// 1. Instantiate Repositories and Services
	approvalRepo := approval.NewInMemoryApprovalRepository()
	approvalUseCase := approval.NewApprovalUseCase(approvalRepo)

	postProcessor := pipeline.NewPostProcessor()

	publisherRepo := publisher.NewInMemoryPublisherRepository()
	mockLinkedInClient := publisher.NewMockLinkedInClient("urn:li:share:123456", false)
	publishUseCase, _ := publisher.NewPublishApprovedPostUseCase(mockLinkedInClient, publisherRepo, "urn:li:person:author")

	// 2. Instantiate Orchestrator
	fewShotRepo := pipeline.NewInMemoryFewShotRepository()
	orchestrator := platform.NewEventOrchestrator(approvalUseCase, postProcessor, publishUseCase, fewShotRepo, nil, nil)

	// 3. Test Ingestion and Generation (ProcessGitHubEvent)
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

	// 4. Test Approval and Publishing (ApproveAndPublish)
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

func TestEventOrchestrator_ReleaseSupersedesPendingDrafts(t *testing.T) {
	ctx := context.Background()
	approvalRepo := approval.NewInMemoryApprovalRepository()
	approvalUseCase := approval.NewApprovalUseCase(approvalRepo)
	postProcessor := pipeline.NewPostProcessor()
	publisherRepo := publisher.NewInMemoryPublisherRepository()
	mockLinkedInClient := publisher.NewMockLinkedInClient("urn:li:share:123", false)
	publishUseCase, _ := publisher.NewPublishApprovedPostUseCase(mockLinkedInClient, publisherRepo, "urn:li:person:author")
	fewShotRepo := pipeline.NewInMemoryFewShotRepository()

	orchestrator := platform.NewEventOrchestrator(approvalUseCase, postProcessor, publishUseCase, fewShotRepo, nil, &mockComparator{})

	// 1. Create a pending draft
	pendingPR, err := approvalUseCase.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "PR draft to be superseded"})
	if err != nil {
		t.Fatalf("unexpected error creating pending PR: %v", err)
	}

	// 2. Process Release event
	payload := webhook.GitHubPayload{
		EventType:  "release",
		Release:    webhook.ReleaseInfo{TagName: "v1.2.0"},
		Repository: webhook.RepositoryInfo{FullName: "hawksxo/git-to-feed"},
	}
	_, err = orchestrator.ProcessGitHubPayload(ctx, payload)
	if err != nil {
		t.Fatalf("unexpected error in ProcessGitHubPayload: %v", err)
	}

	// 3. Verify that the pending draft has been invalidated
	updated, err := approvalRepo.FindByUUID(ctx, pendingPR.UUID)
	if err != nil {
		t.Fatalf("unexpected error finding PR: %v", err)
	}
	if updated.Status != approval.StatusSupersededByRelease {
		t.Errorf("updated.Status = %s, want %s", updated.Status, approval.StatusSupersededByRelease)
	}
	if updated.SupersededBy != "v1.2.0" {
		t.Errorf("updated.SupersededBy = %s, want %s", updated.SupersededBy, "v1.2.0")
	}
}
