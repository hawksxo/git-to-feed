package publisher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockStore struct {
	posts map[string]*PublishedPost
}

func newMockStore() *mockStore {
	return &mockStore{posts: make(map[string]*PublishedPost)}
}

func (m *mockStore) Save(ctx context.Context, post *PublishedPost) error {
	m.posts[post.ID] = post
	return nil
}

func (m *mockStore) FindByUUID(ctx context.Context, uuid string) (*PublishedPost, error) {
	if p, ok := m.posts[uuid]; ok {
		return p, nil
	}
	return nil, ErrNotFound
}

func (m *mockStore) FindByApprovalUUID(ctx context.Context, approvalUUID string) (*PublishedPost, error) {
	for _, p := range m.posts {
		if p.ApprovalPostUUID == approvalUUID {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func TestModelTransitions(t *testing.T) {
	post, err := NewPublishedPost("app-123")
	if err != nil {
		t.Fatalf("unexpected error creating post: %v", err)
	}

	if post.Status != PublishStatusPending {
		t.Errorf("expected PENDING, got %s", post.Status)
	}

	if err := post.MarkAsPublished("urn:li:share:123"); err != nil {
		t.Fatalf("unexpected error marking published: %v", err)
	}
	if post.Status != PublishStatusPublished {
		t.Errorf("expected PUBLISHED, got %s", post.Status)
	}

	// Transicion duplicada debe fallar
	if err := post.MarkAsPublished("urn:li:share:456"); err != ErrInvalidStatusTransition {
		t.Errorf("expected ErrInvalidStatusTransition, got %v", err)
	}
}

func TestServicePublishSuccess(t *testing.T) {
	mockClient := NewMockLinkedInClient("urn:li:share:test", false)
	store := newMockStore()
	svc, err := NewService(mockClient, store, "urn:li:person:author")
	if err != nil {
		t.Fatalf("unexpected error initializing service: %v", err)
	}

	pub, err := svc.Publish(context.Background(), "app-uuid-1", "Contenido a publicar")
	if err != nil {
		t.Fatalf("unexpected error publishing: %v", err)
	}

	if pub.Status != PublishStatusPublished {
		t.Errorf("expected PUBLISHED, got %s", pub.Status)
	}
	if pub.LinkedInShareURN != "urn:li:share:test" {
		t.Errorf("expected share URN, got %s", pub.LinkedInShareURN)
	}
}

func TestServicePublishFailure(t *testing.T) {
	mockClient := NewMockLinkedInClient("", true) // Falla a proposito
	store := newMockStore()
	svc, _ := NewService(mockClient, store, "urn:li:person:author")

	pub, err := svc.Publish(context.Background(), "app-uuid-2", "Contenido a fallar")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if pub.Status != PublishStatusFailed {
		t.Errorf("expected FAILED, got %s", pub.Status)
	}
}

func TestHTTPLinkedInClientSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id": "urn:li:share:created-123"}`))
	}))
	defer server.Close()

	client, err := NewHTTPLinkedInClient("test-token", server.Client())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	client.baseURL = server.URL

	urn, err := client.SharePost(context.Background(), "urn:li:person:author", "Post test")
	if err != nil {
		t.Fatalf("unexpected error sharing post: %v", err)
	}

	if urn != "urn:li:share:created-123" {
		t.Errorf("expected urn:li:share:created-123, got %s", urn)
	}
}
