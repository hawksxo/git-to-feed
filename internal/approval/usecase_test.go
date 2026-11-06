package approval

import (
	"context"
	"testing"
	"time"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

func TestSubmitForApproval(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	genPost := pipeline.GeneratedPost{
		Content:   "Automated Post",
		Archetype: pipeline.ArchetypeFeature,
	}

	appPost, err := uc.SubmitForApproval(ctx, genPost)
	if err != nil {
		t.Fatalf("SubmitForApproval() unexpected error: %v", err)
	}

	if appPost.UUID == "" {
		t.Errorf("expected non-empty UUID")
	}
	if appPost.Status != StatusPending {
		t.Errorf("expected status %v, got %v", StatusPending, appPost.Status)
	}
}

func TestApproveSuccessAndTransitionError(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Test content"})

	approvedPost, err := uc.Approve(ctx, appPost.UUID)
	if err != nil {
		t.Fatalf("Approve() unexpected error: %v", err)
	}
	if approvedPost.Status != StatusApproved {
		t.Errorf("expected status %v, got %v", StatusApproved, approvedPost.Status)
	}

	_, err = uc.Approve(ctx, appPost.UUID)
	if err != ErrInvalidTransition {
		t.Errorf("expected ErrInvalidTransition on already approved post, got %v", err)
	}

	_, err = uc.Reject(ctx, appPost.UUID, "Too late")
	if err != ErrInvalidTransition {
		t.Errorf("expected ErrInvalidTransition on already approved post, got %v", err)
	}
}

func TestRejectSuccessAndEmptyReason(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Test content"})

	_, err := uc.Reject(ctx, appPost.UUID, "   ")
	if err != ErrEmptyRejectReason {
		t.Errorf("expected ErrEmptyRejectReason, got %v", err)
	}

	rejectedPost, err := uc.Reject(ctx, appPost.UUID, "Low quality content")
	if err != nil {
		t.Fatalf("Reject() unexpected error: %v", err)
	}
	if rejectedPost.Status != StatusRejected {
		t.Errorf("expected status %v, got %v", StatusRejected, rejectedPost.Status)
	}
	if rejectedPost.RejectReason != "Low quality content" {
		t.Errorf("expected reject reason 'Low quality content', got %v", rejectedPost.RejectReason)
	}
}

func TestEditAndApprove(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Original content"})

	_, err := uc.EditAndApprove(ctx, appPost.UUID, "")
	if err != ErrEmptyEditedContent {
		t.Errorf("expected ErrEmptyEditedContent, got %v", err)
	}

	editedPost, err := uc.EditAndApprove(ctx, appPost.UUID, "Refined high-value post")
	if err != nil {
		t.Fatalf("EditAndApprove() unexpected error: %v", err)
	}
	if editedPost.Status != StatusApproved {
		t.Errorf("expected status %v, got %v", StatusApproved, editedPost.Status)
	}
	if editedPost.EditedContent != "Refined high-value post" {
		t.Errorf("expected edited content 'Refined high-value post', got %v", editedPost.EditedContent)
	}

	stored, _ := repo.FindByUUID(ctx, appPost.UUID)
	if stored.EditedContent != "Refined high-value post" {
		t.Errorf("stored post does not reflect edited content")
	}
}

func TestListPending(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	post1, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Post 1"})
	_, _ = uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Post 2"})

	_, _ = uc.Approve(ctx, post1.UUID)

	pending, err := uc.ListPending(ctx)
	if err != nil {
		t.Fatalf("ListPending() unexpected error: %v", err)
	}
	if len(pending) != 1 {
		t.Errorf("expected 1 pending post, got %d", len(pending))
	}
}

func TestApproveAndRejectExpiredDraft(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Old post"})
	appPost.CreatedAt = time.Now().Add(-15 * 24 * time.Hour)
	_ = repo.Update(ctx, appPost)

	_, err := uc.Approve(ctx, appPost.UUID)
	if err != ErrDraftExpired {
		t.Errorf("Approve() error = %s, want %s", err, ErrDraftExpired)
	}

	appPost2, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Old post"})
	appPost2.CreatedAt = time.Now().Add(-15 * 24 * time.Hour)
	_ = repo.Update(ctx, appPost2)

	_, err2 := uc.Reject(ctx, appPost2.UUID, "rejected after expired")
	if err2 != ErrDraftExpired {
		t.Errorf("Reject() error = %s, want %s", err2, ErrDraftExpired)
	}
}

func TestSupersedeDrafts(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	draft, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "New post"})
	err := uc.SupersedeDrafts(ctx, []string{draft.UUID}, "v1.2.0")
	if err != nil {
		t.Fatalf("SupersedeDrafts() unexpected error: %v", err)
	}

	postFind, err := repo.FindByUUID(ctx, draft.UUID)
	if err != nil {
		t.Fatalf("FindByUUID() unexpected error: %v", err)
	}
	if postFind.Status != StatusSupersededByRelease {
		t.Errorf("post.Status = %s, want %s", postFind.Status, StatusSupersededByRelease)
	}
	if postFind.SupersededBy != "v1.2.0" {
		t.Errorf("post.SupersededBy = %s, want %s", postFind.SupersededBy, "v1.2.0")
	}
}
