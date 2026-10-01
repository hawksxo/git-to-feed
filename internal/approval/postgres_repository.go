package approval

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

type PostgresApprovalRepository struct {
	db *sql.DB
}

func NewPostgresApprovalRepository(db *sql.DB) *PostgresApprovalRepository {
	return &PostgresApprovalRepository{db: db}
}

func (r *PostgresApprovalRepository) Save(ctx context.Context, post *ApprovalPost) error {
	query := `
		INSERT INTO approval_posts (uuid, idempotency_key, content, archetype, status, edited_content, reject_reason, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (idempotency_key) DO NOTHING
	`
	var idempotencyKey *string
	if post.IdempotencyKey != "" {
		idempotencyKey = &post.IdempotencyKey
	}

	_, err := r.db.ExecContext(
		ctx, query,
		post.UUID,
		idempotencyKey,
		post.GeneratedPost.Content,
		string(post.GeneratedPost.Archetype),
		string(post.Status),
		post.EditedContent,
		post.RejectReason,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error saving approval_post to postgres: %w", err)
	}
	return nil
}

func (r *PostgresApprovalRepository) FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error) {
	query := `
		SELECT uuid, content, archetype, status, edited_content, reject_reason, created_at, updated_at
		FROM approval_posts
		WHERE uuid = $1
	`
	var post ApprovalPost
	var archetypeStr, statusStr string

	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&post.UUID,
		&post.GeneratedPost.Content,
		&archetypeStr,
		&statusStr,
		&post.EditedContent,
		&post.RejectReason,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("approval_post no encontrado: %w", err)
		}
		return nil, fmt.Errorf("error buscando approval_post por UUID: %w", err)
	}

	post.GeneratedPost.Archetype = pipeline.Archetype(archetypeStr)
	post.Status = ApprovalStatus(statusStr)
	return &post, nil
}

func (r *PostgresApprovalRepository) FindAllPending(ctx context.Context) ([]ApprovalPost, error) {
	query := `
		SELECT uuid, content, archetype, status, edited_content, reject_reason, created_at, updated_at
		FROM approval_posts
		WHERE status = 'PENDING'
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error consultando approval_posts pendientes: %w", err)
	}
	defer rows.Close()

	var posts []ApprovalPost
	for rows.Next() {
		var post ApprovalPost
		var archetypeStr, statusStr string

		if err := rows.Scan(
			&post.UUID,
			&post.GeneratedPost.Content,
			&archetypeStr,
			&statusStr,
			&post.EditedContent,
			&post.RejectReason,
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error escaneando fila de approval_post: %w", err)
		}

		post.GeneratedPost.Archetype = pipeline.Archetype(archetypeStr)
		post.Status = ApprovalStatus(statusStr)
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostgresApprovalRepository) Update(ctx context.Context, post *ApprovalPost) error {
	query := `
		UPDATE approval_posts
		SET status = $1, edited_content = $2, reject_reason = $3, updated_at = $4
		WHERE uuid = $5
	`
	res, err := r.db.ExecContext(
		ctx, query,
		string(post.Status),
		post.EditedContent,
		post.RejectReason,
		post.UpdatedAt,
		post.UUID,
	)
	if err != nil {
		return fmt.Errorf("error actualizando approval_post en postgres: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return fmt.Errorf("approval_post no encontrado para actualizar")
	}

	return nil
}
