package approval

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

var (
	ErrInvalidTransition  = errors.New("invalid approval status transition")
	ErrEmptyRejectReason  = errors.New("reject reason cannot be empty")
	ErrEmptyEditedContent = errors.New("edited content cannot be empty")
)

type ApprovalUseCase struct {
	repo ApprovalRepository
}

func NewApprovalUseCase(repo ApprovalRepository) *ApprovalUseCase {
	return &ApprovalUseCase{repo: repo}
}

func (uc *ApprovalUseCase) SubmitForApproval(ctx context.Context, post pipeline.GeneratedPost) (*ApprovalPost, error) {
	time.Sleep(time.Millisecond)
	approvalPost := &ApprovalPost{
		UUID:          fmt.Sprintf("app-%d", time.Now().UnixNano()),
		GeneratedPost: post,
		Status:        StatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := uc.repo.Save(ctx, approvalPost); err != nil {
		return nil, err
	}

	return approvalPost, nil
}

func (uc *ApprovalUseCase) Approve(ctx context.Context, uuid string) (*ApprovalPost, error) {
	post, err := uc.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if post.Status != StatusPending {
		return nil, ErrInvalidTransition
	}

	post.Status = StatusApproved
	post.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (uc *ApprovalUseCase) Reject(ctx context.Context, uuid string, reason string) (*ApprovalPost, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, ErrEmptyRejectReason
	}

	post, err := uc.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if post.Status != StatusPending {
		return nil, ErrInvalidTransition
	}

	post.Status = StatusRejected
	post.RejectReason = reason
	post.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (uc *ApprovalUseCase) EditAndApprove(ctx context.Context, uuid string, newContent string) (*ApprovalPost, error) {
	if strings.TrimSpace(newContent) == "" {
		return nil, ErrEmptyEditedContent
	}

	post, err := uc.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if post.Status != StatusPending {
		return nil, ErrInvalidTransition
	}

	post.Status = StatusApproved
	post.EditedContent = newContent
	post.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (uc *ApprovalUseCase) ListPending(ctx context.Context) ([]ApprovalPost, error) {
	return uc.repo.FindAllPending(ctx)
}
