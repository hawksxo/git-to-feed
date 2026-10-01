package approval

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrPostNotFound = errors.New("approval post not found")
)

type InMemoryApprovalRepository struct {
	mu    sync.RWMutex
	posts map[string]*ApprovalPost
}

func NewInMemoryApprovalRepository() *InMemoryApprovalRepository  {
	return &InMemoryApprovalRepository{
		posts: make(map[string]*ApprovalPost),
	}
}

func (r *InMemoryApprovalRepository) Save(ctx context.Context, post *ApprovalPost) error  {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.UUID] = post
	return nil
}

func (r *InMemoryApprovalRepository) FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[uuid]
	if !exists {
		return nil, ErrPostNotFound
	}

	return post, nil
}

func (r *InMemoryApprovalRepository) FindAllPending(ctx context.Context) ([]ApprovalPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var pending []ApprovalPost
	for _, post := range r.posts {
		if post.Status == StatusPending {
			pending = append(pending, *post)
		}
	}
	return pending, nil
}

func (r *InMemoryApprovalRepository) Update(ctx context.Context, post *ApprovalPost) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.posts[post.UUID]; !exists {
		return ErrPostNotFound
	}
	r.posts[post.UUID] = post
	return nil
}