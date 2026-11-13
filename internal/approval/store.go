package approval

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrPostNotFound = errors.New("el post aprobado no se encuentra")
)

type Store interface {
	Save(ctx context.Context, post *ApprovalPost) error
	FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error)
	FindAllPending(ctx context.Context) ([]ApprovalPost, error)
	Update(ctx context.Context, post *ApprovalPost) error
	SupersedePending(ctx context.Context, uuids []string, releaseTag string) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db: db,
	}
}

type inMemoryStore struct {
	posts map[string]*ApprovalPost
}

func NewInMemoryStore() Store {
	return &inMemoryStore{
		posts: make(map[string]*ApprovalPost),
	}
}

func (m *inMemoryStore) Save(ctx context.Context, post *ApprovalPost) error {
	m.posts[post.UUID] = post
	return nil
}

func (m *inMemoryStore) FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error) {
	p, ok := m.posts[uuid]
	if !ok {
		return nil, ErrPostNotFound
	}
	return p, nil
}

func (m *inMemoryStore) FindAllPending(ctx context.Context) ([]ApprovalPost, error) {
	var list []ApprovalPost
	for _, p := range m.posts {
		if p.Status == StatusPending {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (m *inMemoryStore) Update(ctx context.Context, post *ApprovalPost) error {
	m.posts[post.UUID] = post
	return nil
}

func (m *inMemoryStore) SupersedePending(ctx context.Context, uuids []string, releaseTag string) error {
	for _, id := range uuids {
		if p, ok := m.posts[id]; ok && p.Status == StatusPending {
			p.Status = StatusSupersededByRelease
			p.SupersededBy = releaseTag
		}
	}
	return nil
}

func (s *PostgresStore) Save(ctx context.Context, post *ApprovalPost) error {
	query := `
		INSERT INTO approval_posts(uuid, idempotency_key, content, archetype, status) 
		VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (idempotency_key) DO NOTHING
	`

	var idempotencyKey *string
	if post.IdempotencyKey != "" {
		idempotencyKey = &post.IdempotencyKey
	}

	_, err := s.db.ExecContext(ctx, query, post.UUID, idempotencyKey, post.Content, post.Archetype, post.Status)
	if err != nil {
		return fmt.Errorf("error al intentar guardar el post aprobado: %w", err)
	}

	return nil
}

func (s *PostgresStore) FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error) {
	query := `
		SELECT uuid, content, archetype, status, edited_content, reject_reason
		FROM approval_posts
		WHERE uuid = $1
	`

	var post ApprovalPost
	err := s.db.QueryRowContext(ctx, query, uuid).Scan(
		&post.UUID,
		&post.Content,
		&post.Archetype,
		&post.Status,
		&post.EditedContent,
		&post.RejectReason,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (s *PostgresStore) FindAllPending(ctx context.Context) ([]ApprovalPost, error) {
	query := `
		SELECT uuid, content, archetype, status, edited_content, reject_reason, created_at
		FROM approval_posts
		WHERE status = $1 AND created_at >= NOW() - INTERVAL '14 days'
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, "PENDING")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []ApprovalPost

	for rows.Next() {
		var p ApprovalPost

		if err := rows.Scan(&p.UUID, &p.Content, &p.Archetype, &p.Status, &p.EditedContent, &p.RejectReason, &p.CreatedAt); err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	return posts, nil
}

func (s *PostgresStore) Update(ctx context.Context, post *ApprovalPost) error {
	query := `
		UPDATE approval_posts
		SET status = $1,
			edited_content = $2,
			reject_reason = $3,
			updated_at = NOW()
		WHERE uuid = $4 AND status = 'PENDING'
	`

	res, err := s.db.ExecContext(ctx, query, post.Status, post.EditedContent, post.RejectReason, post.UUID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return ErrPostNotFound
	}

	return nil
}

func (s *PostgresStore) SupersedePending(ctx context.Context, uuids []string, releaseTag string) error {
	if len(uuids) == 0 {
		return nil
	}

	query := `
		UPDATE approval_posts
		SET status = 'SUPERSEDED_BY_RELEASE',
			superseded_by = $1,
			updated_at = NOW()
		WHERE uuid = ANY($2) AND status = 'PENDING'
	`

	_, err := s.db.ExecContext(ctx, query, releaseTag, pq.Array(uuids))
	if err != nil {
		return err
	}

	return nil
}
