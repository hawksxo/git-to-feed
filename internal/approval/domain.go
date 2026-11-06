package approval

import (
	"context"
	"time"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

type ApprovalStatus string

const StatusSupersededByRelease ApprovalStatus = "SUPERSEDED_BY_RELEASE"

const (
	StatusPending  ApprovalStatus = "PENDING"
	StatusApproved ApprovalStatus = "APPROVED"
	StatusRejected ApprovalStatus = "REJECTED"
	StatusExpired  ApprovalStatus = "EXPIRED"
)

func (s ApprovalStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusApproved, StatusRejected, StatusExpired, StatusSupersededByRelease:
		return true
	default:
		return false
	}
}

type ApprovalPost struct {
	UUID           string
	IdempotencyKey string
	GeneratedPost  pipeline.GeneratedPost
	Status         ApprovalStatus
	SupersededBy   string
	EditedContent  string
	RejectReason   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ApprovalRepository interface {
	Save(ctx context.Context, post *ApprovalPost) error
	FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error)
	FindAllPending(ctx context.Context) ([]ApprovalPost, error)
	Update(ctx context.Context, post *ApprovalPost) error
	SupersedePending(ctx context.Context, uuids []string, releaseTag string) error
}
