package approval

import (
	"context"
	"testing"
	"time"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

func TestInMemoryApprovalRepositorySaveAndFindByUUID(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	post := &ApprovalPost{
		UUID: "test-uuid-1",
		GeneratedPost: pipeline.GeneratedPost{
			Content:   "Test post content",
			Archetype: pipeline.ArchetypeRelease,
			CreatedAt: time.Now(),
		},
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Save(ctx, post)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	found, err := repo.FindByUUID(ctx, "test-uuid-1")
	if err != nil {
		t.Fatalf("FindByUUID() unexpected error: %v", err)
	}

	if found.UUID != "test-uuid-1" {
		t.Errorf("found.UUID = %q, want %q", found.UUID, "test-uuid-1")
	}

	if found.Status != StatusPending {
		t.Errorf("found.Status = %q, want %q", found.Status, StatusPending)
	}
}

func TestInMemoryApprovalRepositoryNotFound(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	_, err := repo.FindByUUID(ctx, "non-existent-uuid")
	if err != ErrPostNotFound {
		t.Errorf("FindByUUID() error = %v, want %v", err, ErrPostNotFound)
	}
}

func TestInMemoryApprovalRepositoryFindAllPending(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	post1 := &ApprovalPost{UUID: "uuid-1", Status: StatusPending}
	post2 := &ApprovalPost{UUID: "uuid-2", Status: StatusApproved}
	post3 := &ApprovalPost{UUID: "uuid-3", Status: StatusPending}

	_ = repo.Save(ctx, post1)
	_ = repo.Save(ctx, post2)
	_ = repo.Save(ctx, post3)

	pending, err := repo.FindAllPending(ctx)
	if err != nil {
		t.Fatalf("FindAllPending() unexpected error: %v", err)
	}

	if len(pending) != 2 {
		t.Errorf("len(pending) = %d, want %d", len(pending), 2)
	}
}

func TestInMemoryApprovalRepositoryUpdate(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	post := &ApprovalPost{UUID: "uuid-1", Status: StatusPending}
	_ = repo.Save(ctx, post)

	post.Status = StatusApproved
	err := repo.Update(ctx, post)
	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	updated, _ := repo.FindByUUID(ctx, "uuid-1")
	if updated.Status != StatusApproved {
		t.Errorf("updated.Status = %q, want %q", updated.Status, StatusApproved)
	}

	nonExistent := &ApprovalPost{UUID: "non-existent", Status: StatusApproved}
	err = repo.Update(ctx, nonExistent)
	if err != ErrPostNotFound {
		t.Errorf("Update() non-existent error = %v, want %v", err, ErrPostNotFound)
	}
}