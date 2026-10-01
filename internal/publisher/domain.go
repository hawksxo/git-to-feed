package publisher

import (
	"context"
	"errors"
	"crypto/rand"
	"fmt"
	"time"
)

var (
	ErrEmptyApprovalUUID       = errors.New("approval post UUID cannot be empty")
	ErrInvalidStatusTransition = errors.New("invalid publish status transition")
	ErrEmptyShareURN           = errors.New("linkedin share URN cannot be empty")
	ErrEmptyFailureReason      = errors.New("failure reason cannot be empty")
)

type PublishStatus string

const (
	PublishStatusPending   PublishStatus = "PENDING"
	PublishStatusPublished PublishStatus = "PUBLISHED"
	PublishStatusFailed    PublishStatus = "FAILED"
)

type PublishedPost struct {
	ID               string        `json:"id"`
	ApprovalPostUUID string        `json:"approval_post_uuid"`
	LinkedInShareURN string        `json:"linkedin_share_urn,omitempty"`
	Status           PublishStatus `json:"status"`
	ErrorMessage     string        `json:"error_message,omitempty"`
	PublishedAt      time.Time     `json:"published_at"`
}

func NewPublishedPost(approvalPostUUID string) (*PublishedPost, error) {
	if approvalPostUUID == "" {
		return nil, ErrEmptyApprovalUUID
	}

	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return &PublishedPost{
		ID:               uuid,
		ApprovalPostUUID: approvalPostUUID,
		Status:           PublishStatusPending,
		PublishedAt:      time.Now().UTC(),
	}, nil
}

func (p *PublishedPost) MarkAsPublished(shareURN string) error {
	if p.Status != PublishStatusPending {
		return ErrInvalidStatusTransition
	}
	if shareURN == "" {
		return ErrEmptyShareURN
	}
	p.Status = PublishStatusPublished
	p.LinkedInShareURN = shareURN
	p.PublishedAt = time.Now().UTC()
	return nil
}

func (p *PublishedPost) MarkAsFailed(reason string) error {
	if p.Status != PublishStatusPending {
		return ErrInvalidStatusTransition
	}
	if reason == "" {
		return ErrEmptyFailureReason
	}
	p.Status = PublishStatusFailed
	p.ErrorMessage = reason
	p.PublishedAt = time.Now().UTC()
	return nil
}

type LinkedInClient interface {
	SharePost(ctx context.Context, authorURN string, text string) (shareURN string, err error)
}

type PublisherRepository interface {
	Save(ctx context.Context, post *PublishedPost) error
	FindByUUID(ctx context.Context, uuid string) (*PublishedPost, error)
	FindByApprovalUUID(ctx context.Context, approvalUUID string) (*PublishedPost, error)
}
