package publisher

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrNotFound      = errors.New("published post not found")
	ErrAlreadyExists = errors.New("published post record already exists")
)

type InMemoryPublisherRepository struct {
	mu    sync.RWMutex
	posts map[string]*PublishedPost
}

func NewInMemoryPublisherRepository() *InMemoryPublisherRepository {
	return &InMemoryPublisherRepository{
		posts: make(map[string]*PublishedPost),
	}
}

func (r *InMemoryPublisherRepository) Save(ctx context.Context, post *PublishedPost) error {
	if post == nil {
		return errors.New("cannot save nil post")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.posts[post.ID] = post
	return nil
}

func (r *InMemoryPublisherRepository) FindByUUID(ctx context.Context, uuid string) (*PublishedPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	post, exists := r.posts[uuid]
	if !exists {
		return nil, ErrNotFound
	}
	return post, nil
}

func (r *InMemoryPublisherRepository) FindByApprovalUUID(ctx context.Context, approvalUUID string) (*PublishedPost, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, post := range r.posts {
		if post.ApprovalPostUUID == approvalUUID {
			return post, nil
		}
	}
	return nil, ErrNotFound
}
