package publisher

import (
	"context"
	"errors"
	"fmt"

	"github.com/hawksxo/git-to-feed/internal/approval"
)

var (
	ErrNilApprovalPost    = errors.New("approval post cannot be nil")
	ErrPostNotApproved    = errors.New("only approved posts can be published to linkedin")
	ErrNilLinkedInClient  = errors.New("linkedin client cannot be nil")
	ErrNilRepository      = errors.New("publisher repository cannot be nil")
)

type PublishApprovedPostUseCase struct {
	client    LinkedInClient
	repo      PublisherRepository
	authorURN string
}

func NewPublishApprovedPostUseCase(client LinkedInClient, repo PublisherRepository, authorURN string) (*PublishApprovedPostUseCase, error) {
	if client == nil {
		return nil, ErrNilLinkedInClient
	}
	if repo == nil {
		return nil, ErrNilRepository
	}
	if authorURN == "" {
		return nil, ErrEmptyAuthorURN
	}
	return &PublishApprovedPostUseCase{
		client:    client,
		repo:      repo,
		authorURN: authorURN,
	}, nil
}

func (uc *PublishApprovedPostUseCase) Execute(ctx context.Context, approvalPost *approval.ApprovalPost) (*PublishedPost, error) {
	if approvalPost == nil {
		return nil, ErrNilApprovalPost
	}
	if approvalPost.Status != approval.StatusApproved {
		return nil, ErrPostNotApproved
	}

	pubRecord, err := NewPublishedPost(approvalPost.UUID)
	if err != nil {
		return nil, fmt.Errorf("error creando registro de publicacion: %w", err)
	}

	fullText := approvalPost.GeneratedPost.Content
	if approvalPost.EditedContent != "" {
		fullText = approvalPost.EditedContent
	}

	shareURN, err := uc.client.SharePost(ctx, uc.authorURN, fullText)
	if err != nil {
		_ = pubRecord.MarkAsFailed(err.Error())
		_ = uc.repo.Save(ctx, pubRecord)
		return pubRecord, fmt.Errorf("error publicando en linkedin: %w", err)
	}

	if err := pubRecord.MarkAsPublished(shareURN); err != nil {
		return nil, fmt.Errorf("error actualizando estado a publicado: %w", err)
	}

	if err := uc.repo.Save(ctx, pubRecord); err != nil {
		return nil, fmt.Errorf("error guardando registro de publicacion: %w", err)
	}

	return pubRecord, nil
}
