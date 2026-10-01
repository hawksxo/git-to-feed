package publisher

import (
	"context"
	"database/sql"
	"fmt"
)

type PostgresPublisherRepository struct {
	db *sql.DB
}

func NewPostgresPublisherRepository(db *sql.DB) *PostgresPublisherRepository {
	return &PostgresPublisherRepository{db: db}
}

func (r *PostgresPublisherRepository) Save(ctx context.Context, post *PublishedPost) error {
	query := `
		INSERT INTO published_posts (id, approval_post_uuid, linkedin_share_urn, status, error_message, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE
		SET linkedin_share_urn = EXCLUDED.linkedin_share_urn,
		    status = EXCLUDED.status,
		    error_message = EXCLUDED.error_message,
		    published_at = EXCLUDED.published_at
	`
	_, err := r.db.ExecContext(
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

func (r *PostgresPublisherRepository) FindByUUID(ctx context.Context, uuid string) (*PublishedPost, error) {
	query := `
		SELECT id, approval_post_uuid, linkedin_share_urn, status, error_message, published_at
		FROM published_posts
		WHERE id = $1
	`
	var post PublishedPost
	var statusStr string

	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&post.ID,
		&post.ApprovalPostUUID,
		&post.LinkedInShareURN,
		&statusStr,
		&post.ErrorMessage,
		&post.PublishedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error buscando published_post por UUID: %w", err)
	}

	post.Status = PublishStatus(statusStr)
	return &post, nil
}

func (r *PostgresPublisherRepository) FindByApprovalUUID(ctx context.Context, approvalUUID string) (*PublishedPost, error) {
	query := `
		SELECT id, approval_post_uuid, linkedin_share_urn, status, error_message, published_at
		FROM published_posts
		WHERE approval_post_uuid = $1
	`
	var post PublishedPost
	var statusStr string

	err := r.db.QueryRowContext(ctx, query, approvalUUID).Scan(
		&post.ID,
		&post.ApprovalPostUUID,
		&post.LinkedInShareURN,
		&statusStr,
		&post.ErrorMessage,
		&post.PublishedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error buscando published_post por ApprovalUUID: %w", err)
	}

	post.Status = PublishStatus(statusStr)
	return &post, nil
}
