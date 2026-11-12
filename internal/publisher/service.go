package publisher

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNilLinkedInClient = errors.New("linkedin client cannot be nil")
	ErrNilStore          = errors.New("publisher store cannot be nil")
)

type Service struct {
	client    LinkedInClient
	store     Store
	authorURN string
}

func NewService(client LinkedInClient, store Store, authorURN string) (*Service, error) {
	if client == nil {
		return nil, ErrNilLinkedInClient
	}
	if store == nil {
		return nil, ErrNilStore
	}
	if authorURN == "" {
		return nil, ErrEmptyAuthorURN
	}
	return &Service{
		client:    client,
		store:     store,
		authorURN: authorURN,
	}, nil
}

func (s *Service) Publish(ctx context.Context, approvalUUID string, content string) (*PublishedPost, error) {
	if approvalUUID == "" {
		return nil, ErrEmptyApprovalUUID
	}

	pubRecord, err := NewPublishedPost(approvalUUID)
	if err != nil {
		return nil, fmt.Errorf("error creando registro de publicacion: %w", err)
	}

	shareURN, err := s.client.SharePost(ctx, s.authorURN, content)
	if err != nil {
		_ = pubRecord.MarkAsFailed(err.Error())
		_ = s.store.Save(ctx, pubRecord)
		return pubRecord, fmt.Errorf("error publicando en linkedin: %w", err)
	}

	if err := pubRecord.MarkAsPublished(shareURN); err != nil {
		return nil, fmt.Errorf("error actualizando estado a publicado: %w", err)
	}

	if err := s.store.Save(ctx, pubRecord); err != nil {
		return nil, fmt.Errorf("error persistiendo post publicado: %w", err)
	}

	return pubRecord, nil
}
