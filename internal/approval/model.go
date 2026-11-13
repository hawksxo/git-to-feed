package approval

import (
	"time"
)

type ApprovalStatus string

const (
	StatusPending             ApprovalStatus = "PENDING"
	StatusApproved            ApprovalStatus = "APPROVED"
	StatusRejected            ApprovalStatus = "REJECTED"
	StatusExpired             ApprovalStatus = "EXPIRED"
	StatusSupersededByRelease ApprovalStatus = "SUPERSEDED_BY_RELEASE"
)

func (s ApprovalStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusApproved, StatusRejected, StatusExpired, StatusSupersededByRelease:
		return true
	default:
		return false
	}
}

type ChannelConfig struct {
	DefaultChannelID  string
	ReleasesChannelID string
	FeaturesChannelID string
	AuditChannelID    string
}

type ApprovalPost struct {
	UUID           string
	IdempotencyKey string
	Content        string
	Archetype      string
	Status         ApprovalStatus
	SupersededBy   string
	EditedContent  string
	RejectReason   string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
