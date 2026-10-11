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
		Content:   "Generated LinkedIn content",
		Archetype: pipeline.ArchetypeRelease,
	}

	appPost, err := uc.SubmitForApproval(ctx, genPost)
	if err != nil {
		t.Fatalf("SubmitForApproval() unexpected error: %v", err)
	}

	if appPost.UUID == "" {
		t.Errorf("Expected non-empty UUID")
	}

	if appPost.Status != StatusPending {
		t.Errorf("appPost.Status = %q, want %q", appPost.Status, StatusPending)
	}
}

func TestApproveSuccessAndTransitionError(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Test content"})

	approved, err := uc.Approve(ctx, appPost.UUID)
	if err != nil {
		t.Fatalf("Approve() unexpected error: %v", err)
	}

	if approved.Status != StatusApproved {
		t.Errorf("approved.Status = %q, want %q", approved.Status, StatusApproved)
	}

	// Re-approving an already approved post should fail with transition error
	_, err = uc.Approve(ctx, appPost.UUID)
	if err != ErrInvalidTransition {
		t.Errorf("Approve() re-approval error = %v, want %v", err, ErrInvalidTransition)
	}
}

func TestRejectSuccessAndEmptyReason(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Test content"})

	// Empty reason error
	_, err := uc.Reject(ctx, appPost.UUID, "   ")
	if err != ErrEmptyRejectReason {
		t.Errorf("Reject() empty reason error = %v, want %v", err, ErrEmptyRejectReason)
	}

	// Valid reject
	rejected, err := uc.Reject(ctx, appPost.UUID, "Invalid tone")
	if err != nil {
		t.Fatalf("Reject() unexpected error: %v", err)
	}

	if rejected.Status != StatusRejected {
		t.Errorf("rejected.Status = %q, want %q", rejected.Status, StatusRejected)
	}

	if rejected.RejectReason != "Invalid tone" {
		t.Errorf("rejected.RejectReason = %q, want %q", rejected.RejectReason, "Invalid tone")
	}
}

func TestEditAndApprove(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Original content"})

	// Empty content error
	_, err := uc.EditAndApprove(ctx, appPost.UUID, "")
	if err != ErrEmptyEditedContent {
		t.Errorf("EditAndApprove() empty content error = %v, want %v", err, ErrEmptyEditedContent)
	}

	// Valid edit and approve
	edited, err := uc.EditAndApprove(ctx, appPost.UUID, "Edited content for LinkedIn")
	if err != nil {
		t.Fatalf("EditAndApprove() unexpected error: %v", err)
	}

	if edited.Status != StatusApproved {
		t.Errorf("edited.Status = %q, want %q", edited.Status, StatusApproved)
	}

	if edited.EditedContent != "Edited content for LinkedIn" {
		t.Errorf("edited.EditedContent = %q, want %q", edited.EditedContent, "Edited content for LinkedIn")
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
		t.Errorf("len(pending) = %d, want %d", len(pending), 1)
	}
}

func TestApproveAndRejectExpiredDraft(t *testing.T)  {
	repo := NewInMemoryApprovalRepository()
	uc := NewApprovalUseCase(repo)
	ctx := context.Background()

	appPost, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Old post"})
	appPost.CreatedAt = time.Now().Add(-15 * 24 * time.Hour)
	_ = repo.Update(ctx, appPost)
	_, err := uc.Approve(ctx, appPost.UUID)
	if err != ErrDraftExpired {
		t.Errorf("Approve() error = %v, want %v", err, ErrDraftExpired)
	}

	appPost2, _ := uc.SubmitForApproval(ctx, pipeline.GeneratedPost{Content: "Old post"})
	appPost2.CreatedAt = time.Now().Add(-15 * 24 * time.Hour)
	_ = repo.Update(ctx, appPost2)
	_, err2 := uc.Reject(ctx, appPost2.UUID, "discard")
	if err2 != ErrDraftExpired {
		t.Errorf("Reject() error = %s, want %s", err2, ErrDraftExpired)
	}
}
