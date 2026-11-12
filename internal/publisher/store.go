package publisher

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("publicacion no encontrada en la base de datos")
)

type Store interface {
	Save(ctx context.Context, post *PublishedPost) error
	FindByUUID(ctx context.Context, uuid string) (*PublishedPost, error)
	FindByApprovalUUID(ctx context.Context, approvalUUID string) (*PublishedPost, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Save(ctx context.Context, post *PublishedPost) error {
	query := `
		INSERT INTO published_posts (id, approval_post_uuid, linkedin_share_urn, status, error_message, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE
		SET linkedin_share_urn = EXCLUDED.linkedin_share_urn,
		    status = EXCLUDED.status,
		    error_message = EXCLUDED.error_message,
		    published_at = EXCLUDED.published_at
	`
	_, err := s.db.ExecContext(
		ctx, query,
		post.ID,
		post.ApprovalPostUUID,
		post.LinkedInShareURN,
		string(post.Status),
		post.ErrorMessage,
		post.PublishedAt,
	)
	if err != nil {
		return fmt.Errorf("error guardando published_post en postgres: %w", err)
	}
	return nil
}

func (s *PostgresStore) FindByUUID(ctx context.Context, uuid string) (*PublishedPost, error) {
	query := `
		SELECT id, approval_post_uuid, linkedin_share_urn, status, error_message, published_at
		FROM published_posts
		WHERE id = $1
	`
	var post PublishedPost
	var statusStr string

	err := s.db.QueryRowContext(ctx, query, uuid).Scan(
		&post.ID,
		&post.ApprovalPostUUID,
		&post.LinkedInShareURN,
		&statusStr,
		&post.ErrorMessage,
		&post.PublishedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error buscando published_post: %w", err)
	}

	post.Status = PublishStatus(statusStr)
	return &post, nil
}

func (s *PostgresStore) FindByApprovalUUID(ctx context.Context, approvalUUID string) (*PublishedPost, error) {
	query := `
		SELECT id, approval_post_uuid, linkedin_share_urn, status, error_message, published_at
		FROM published_posts
		WHERE approval_post_uuid = $1
	`
	var post PublishedPost
	var statusStr string

	err := s.db.QueryRowContext(ctx, query, approvalUUID).Scan(
		&post.ID,
		&post.ApprovalPostUUID,
		&post.LinkedInShareURN,
		&statusStr,
		&post.ErrorMessage,
		&post.PublishedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error buscando published_post por approval_uuid: %w", err)
	}

	post.Status = PublishStatus(statusStr)
	return &post, nil
}
