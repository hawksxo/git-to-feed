package pipeline

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrEmptyExampleFields = errors.New("input context and expected output cannot be empty")

type FewShotExampleEntity struct {
	UUID           string    `json:"uuid"`
	Archetype      Archetype `json:"archetype"`
	InputContext   string    `json:"input_context"`
	ExpectedOutput string    `json:"expected_output"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type FewShotRepository interface {
	Save(ctx context.Context, example *FewShotExampleEntity) error
	FindActiveByArchetype(ctx context.Context, archetype Archetype) ([]FewShotExampleEntity, error)
	FindAll(ctx context.Context) ([]FewShotExampleEntity, error)
}

type InMemoryFewShotRepository struct {
	mu       sync.RWMutex
	examples map[string]FewShotExampleEntity
}

func NewInMemoryFewShotRepository() *InMemoryFewShotRepository {
	return &InMemoryFewShotRepository{
		examples: make(map[string]FewShotExampleEntity),
	}
}

func (r *InMemoryFewShotRepository) Save(ctx context.Context, example *FewShotExampleEntity) error {
	if example == nil || example.InputContext == "" || example.ExpectedOutput == "" {
		return ErrEmptyExampleFields
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if example.UUID == "" {
		example.UUID = fmt.Sprintf("ex-%d", time.Now().UnixNano())
	}
	if example.CreatedAt.IsZero() {
		example.CreatedAt = time.Now().UTC()
	}

	r.examples[example.UUID] = *example
	return nil
}

func (r *InMemoryFewShotRepository) FindActiveByArchetype(ctx context.Context, archetype Archetype) ([]FewShotExampleEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []FewShotExampleEntity
	for _, ex := range r.examples {
		if ex.Archetype == archetype && ex.IsActive {
			result = append(result, ex)
		}
	}
	return result, nil
}

func (r *InMemoryFewShotRepository) FindAll(ctx context.Context) ([]FewShotExampleEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []FewShotExampleEntity
	for _, ex := range r.examples {
		result = append(result, ex)
	}
	return result, nil
}

type PostgresFewShotRepository struct {
	db *sql.DB
}

func NewPostgresFewShotRepository(db *sql.DB) *PostgresFewShotRepository {
	return &PostgresFewShotRepository{db: db}
}

func (r *PostgresFewShotRepository) Save(ctx context.Context, example *FewShotExampleEntity) error {
	if example == nil || example.InputContext == "" || example.ExpectedOutput == "" {
		return ErrEmptyExampleFields
	}
	if example.UUID == "" {
		example.UUID = fmt.Sprintf("ex-%d", time.Now().UnixNano())
	}
	if example.CreatedAt.IsZero() {
		example.CreatedAt = time.Now().UTC()
	}

	query := `
		INSERT INTO few_shot_examples (uuid, archetype, input_context, expected_output, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (uuid) DO UPDATE
		SET archetype = EXCLUDED.archetype,
		    input_context = EXCLUDED.input_context,
		    expected_output = EXCLUDED.expected_output,
		    is_active = EXCLUDED.is_active
	`
	_, err := r.db.ExecContext(
		ctx, query,
		example.UUID,
		string(example.Archetype),
		example.InputContext,
		example.ExpectedOutput,
		example.IsActive,
		example.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("error guardando few_shot_example en postgres: %w", err)
	}
	return nil
}

func (r *PostgresFewShotRepository) FindActiveByArchetype(ctx context.Context, archetype Archetype) ([]FewShotExampleEntity, error) {
	query := `
		SELECT uuid, archetype, input_context, expected_output, is_active, created_at
		FROM few_shot_examples
		WHERE archetype = $1 AND is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, string(archetype))
	if err != nil {
		return nil, fmt.Errorf("error consultando ejemplos activos: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []FewShotExampleEntity
	for rows.Next() {
		var ex FewShotExampleEntity
		var archStr string

		if err := rows.Scan(&ex.UUID, &archStr, &ex.InputContext, &ex.ExpectedOutput, &ex.IsActive, &ex.CreatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando fila de few_shot_example: %w", err)
		}
		ex.Archetype = Archetype(archStr)
		result = append(result, ex)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando ejemplos activos: %w", err)
	}
	return result, nil
}

func (r *PostgresFewShotRepository) FindAll(ctx context.Context) ([]FewShotExampleEntity, error) {
	query := `
		SELECT uuid, archetype, input_context, expected_output, is_active, created_at
		FROM few_shot_examples
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error consultando todos los ejemplos: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var result []FewShotExampleEntity
	for rows.Next() {
		var ex FewShotExampleEntity
		var archStr string

		if err := rows.Scan(&ex.UUID, &archStr, &ex.InputContext, &ex.ExpectedOutput, &ex.IsActive, &ex.CreatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando fila de few_shot_example: %w", err)
		}
		ex.Archetype = Archetype(archStr)
		result = append(result, ex)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando filas: %w", err)
	}
	return result, nil
}
