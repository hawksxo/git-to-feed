package approval

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidTransition  = errors.New("invalid approval status transition")
	ErrEmptyRejectReason  = errors.New("reject reason cannot be empty")
	ErrEmptyEditedContent = errors.New("edited content cannot be empty")
	ErrPostNotPending     = errors.New("approval post is not in pending status")
	DraftTTL              = 14 * 24 * time.Hour
	ErrDraftExpired       = errors.New("approval draft has expired (older than 14 days)")
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) validatePending(ctx context.Context, post *ApprovalPost) error {
	if post.Status != StatusPending {
		return ErrInvalidTransition
	}

	if time.Since(post.CreatedAt) > DraftTTL {
		post.Status = StatusExpired
		post.UpdatedAt = time.Now()
		_ = s.store.Update(ctx, post)
		return ErrDraftExpired
	}
	
	return nil
}

func (s *Service) SubmitForApproval(ctx context.Context, content string, archetype string, idempotencyKey string) (*ApprovalPost, error) {
	approvalPost := &ApprovalPost{
		UUID:           fmt.Sprintf("app-%d", time.Now().UnixNano()),
		IdempotencyKey: idempotencyKey,
		Content:        content,
		Archetype:      archetype,
		Status:         StatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.store.Save(ctx, approvalPost); err != nil {
		return nil, err
	}

	return approvalPost, nil
}

func (s *Service) Approve(ctx context.Context, uuid string) (*ApprovalPost, error) {
	post, err := s.store.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if err := s.validatePending(ctx, post); err != nil {
		return nil, err
	}

	post.Status = StatusApproved
	post.UpdatedAt = time.Now()

	if err := s.store.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) Reject(ctx context.Context, uuid string, reason string) (*ApprovalPost, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, ErrEmptyRejectReason
	}

	post, err := s.store.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if err := s.validatePending(ctx, post); err != nil {
		return nil, err
	}

	post.Status = StatusRejected
	post.RejectReason = reason
	post.UpdatedAt = time.Now()

	if err := s.store.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) EditAndApprove(ctx context.Context, uuid string, newContent string) (*ApprovalPost, error) {
	if strings.TrimSpace(newContent) == "" {
		return nil, ErrEmptyEditedContent
	}

	post, err := s.store.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if err := s.validatePending(ctx, post); err != nil {
		return nil, err
	}

	post.EditedContent = newContent
	post.Status = StatusApproved
	post.UpdatedAt = time.Now()

	if err := s.store.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Service) ListPending(ctx context.Context) ([]ApprovalPost, error) {
	return s.store.FindAllPending(ctx)
}

func (s *Service) SupersedeDrafts(ctx context.Context, uuids []string, releaseTag string) error {
	return s.store.SupersedePending(ctx, uuids, releaseTag)
}
