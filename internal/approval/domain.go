package approval

import (
	"context"
	"time"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

type ApprovalStatus string

const (
	StatusPending  ApprovalStatus = "PENDING"
	StatusApproved ApprovalStatus = "APPROVED"
	StatusRejected ApprovalStatus = "REJECTED"
)

func (s ApprovalStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusApproved, StatusRejected:
		return true
	default:
		return false
	}
}

type ApprovalPost struct {
	UUID          string
	GeneratedPost pipeline.GeneratedPost
	Status        ApprovalStatus
	EditedContent string
	RejectReason  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ApprovalRepository interface {
	Save(ctx context.Context, post *ApprovalPost) error
	FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error)
	FindAllPending(ctx context.Context) ([]ApprovalPost, error)
	Update(ctx context.Context, post *ApprovalPost) error
}