package approval

import (
	"context"
	"testing"
)

func TestInMemoryApprovalRepositorySaveAndFindByUUID(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	post := &ApprovalPost{
		UUID:   "uuid-123",
		Status: StatusPending,
	}

	err := repo.Save(ctx, post)
	if err != nil {
		t.Fatalf("Save() unexpected error: %v", err)
	}

	found, err := repo.FindByUUID(ctx, "uuid-123")
	if err != nil {
		t.Fatalf("FindByUUID() unexpected error: %v", err)
	}

	if found.UUID != post.UUID {
		t.Errorf("FindByUUID() = %v, want %v", found.UUID, post.UUID)
	}
}

func TestInMemoryApprovalRepositoryNotFound(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	_, err := repo.FindByUUID(ctx, "non-existent")
	if err != ErrPostNotFound {
		t.Errorf("FindByUUID() error = %v, want %v", err, ErrPostNotFound)
	}
}

func TestInMemoryApprovalRepositoryFindAllPending(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	p1 := &ApprovalPost{UUID: "1", Status: StatusPending}
	p2 := &ApprovalPost{UUID: "2", Status: StatusApproved}
	p3 := &ApprovalPost{UUID: "3", Status: StatusPending}

	_ = repo.Save(ctx, p1)
	_ = repo.Save(ctx, p2)
	_ = repo.Save(ctx, p3)

	pending, err := repo.FindAllPending(ctx)
	if err != nil {
		t.Fatalf("FindAllPending() unexpected error: %v", err)
	}

	if len(pending) != 2 {
		t.Errorf("FindAllPending() returned %d posts, want 2", len(pending))
	}
}

func TestInMemoryApprovalRepositoryUpdate(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	post := &ApprovalPost{UUID: "1", Status: StatusPending}
	_ = repo.Save(ctx, post)

	post.Status = StatusApproved
	err := repo.Update(ctx, post)
	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	found, _ := repo.FindByUUID(ctx, "1")
	if found.Status != StatusApproved {
		t.Errorf("Update() status = %v, want %v", found.Status, StatusApproved)
	}

	err = repo.Update(ctx, &ApprovalPost{UUID: "non-existent"})
	if err != ErrPostNotFound {
		t.Errorf("Update() non-existent error = %v, want %v", err, ErrPostNotFound)
	}
}

func TestInMemoryApprovalRepositorySupersedePending(t *testing.T) {
	repo := NewInMemoryApprovalRepository()
	ctx := context.Background()

	post := &ApprovalPost{UUID: "uuid-1", Status: StatusPending}
	post2 := &ApprovalPost{UUID: "uuid-2", Status: StatusApproved}
	_ = repo.Save(ctx, post)
	_ = repo.Save(ctx, post2)

	err := repo.SupersedePending(ctx, []string{post.UUID, post2.UUID}, "v1.2.0")
	if err != nil {
		t.Fatalf("SupersedePending() unexpected error: %v", err)
	}

	updated, _ := repo.FindByUUID(ctx, "uuid-1")
	if updated.Status != StatusSupersededByRelease {
		t.Errorf("Updated.Status = %q, want = %q", updated.Status, StatusSupersededByRelease)
	}

	if updated.SupersededBy != "v1.2.0" {
		t.Errorf("Updated.SupersededBy = %q, want = %q", updated.SupersededBy, "v1.2.0")
	}

	updated2, _ := repo.FindByUUID(ctx, "uuid-2")
	if updated2.Status != StatusApproved {
		t.Errorf("Updated.Status = %q, want = %q", updated2.Status, StatusApproved)
	}
}
