package approval

import (
	"context"
	"testing"
	"time"
)

type mockStore struct {
	posts map[string]*ApprovalPost
}

func newMockStore() *mockStore {
	return &mockStore{posts: make(map[string]*ApprovalPost)}
}

func (m *mockStore) Save(ctx context.Context, post *ApprovalPost) error {
	m.posts[post.UUID] = post
	return nil
}

func (m *mockStore) FindByUUID(ctx context.Context, uuid string) (*ApprovalPost, error) {
	p, ok := m.posts[uuid]
	if !ok {
		return nil, ErrPostNotFound
	}
	return p, nil
}

func (m *mockStore) FindAllPending(ctx context.Context) ([]ApprovalPost, error) {
	var list []ApprovalPost
	for _, p := range m.posts {
		if p.Status == StatusPending {
			list = append(list, *p)
		}
	}
	return list, nil
}

func (m *mockStore) Update(ctx context.Context, post *ApprovalPost) error {
	m.posts[post.UUID] = post
	return nil
}

func (m *mockStore) SupersedePending(ctx context.Context, uuids []string, releaseTag string) error {
	for _, id := range uuids {
		if p, ok := m.posts[id]; ok && p.Status == StatusPending {
			p.Status = StatusSupersededByRelease
			p.SupersededBy = releaseTag
		}
	}
	return nil
}

func TestSubmitForApproval(t *testing.T) {
	store := newMockStore()
	svc := NewService(store)
	ctx := context.Background()

	post, err := svc.SubmitForApproval(ctx, "Contenido de prueba", "feature", "idemp-1")
	if err != nil {
		t.Fatalf("SubmitForApproval() error inesperado: %v", err)
	}

	if post.UUID == "" {
		t.Errorf("se esperaba un UUID no vacio")
	}
	if post.Status != StatusPending {
		t.Errorf("se esperaba estado PENDING, se obtuvo %v", post.Status)
	}
}

func TestApproveSuccessAndTransitionError(t *testing.T) {
	store := newMockStore()
	svc := NewService(store)
	ctx := context.Background()

	post, _ := svc.SubmitForApproval(ctx, "Contenido", "release", "idemp-2")

	approved, err := svc.Approve(ctx, post.UUID)
	if err != nil {
		t.Fatalf("Approve() error inesperado: %v", err)
	}
	if approved.Status != StatusApproved {
		t.Errorf("se esperaba estado APPROVED, se obtuvo %v", approved.Status)
	}

	// Segundo approve debe fallar por transicion invalida
	_, err = svc.Approve(ctx, post.UUID)
	if err != ErrInvalidTransition {
		t.Errorf("se esperaba ErrInvalidTransition, se obtuvo: %v", err)
	}
}

func TestRejectWithReason(t *testing.T) {
	store := newMockStore()
	svc := NewService(store)
	ctx := context.Background()

	post, _ := svc.SubmitForApproval(ctx, "Contenido a rechazar", "feature", "idemp-3")

	// Rechazo con motivo vacio debe fallar
	_, err := svc.Reject(ctx, post.UUID, "   ")
	if err != ErrEmptyRejectReason {
		t.Errorf("se esperaba ErrEmptyRejectReason, se obtuvo: %v", err)
	}

	// Rechazo con motivo valido
	rejected, err := svc.Reject(ctx, post.UUID, "No cumple con las directrices")
	if err != nil {
		t.Fatalf("Reject() error inesperado: %v", err)
	}
	if rejected.Status != StatusRejected {
		t.Errorf("se esperaba estado REJECTED, se obtuvo %v", rejected.Status)
	}
}

func TestEditAndApprove(t *testing.T) {
	store := newMockStore()
	svc := NewService(store)
	ctx := context.Background()

	post, _ := svc.SubmitForApproval(ctx, "Contenido original", "feature", "idemp-4")

	edited, err := svc.EditAndApprove(ctx, post.UUID, "Contenido mejorado manualmente")
	if err != nil {
		t.Fatalf("EditAndApprove() error inesperado: %v", err)
	}
	if edited.Status != StatusApproved {
		t.Errorf("se esperaba estado APPROVED, se obtuvo %v", edited.Status)
	}
	if edited.EditedContent != "Contenido mejorado manualmente" {
		t.Errorf("el contenido editado no coincide: %s", edited.EditedContent)
	}
}

func TestDraftExpiration(t *testing.T) {
	store := newMockStore()
	svc := NewService(store)
	ctx := context.Background()

	post, _ := svc.SubmitForApproval(ctx, "Post antiguo", "feature", "idemp-5")
	// Simulamos que pasaron 15 dias
	post.CreatedAt = time.Now().Add(-15 * 24 * time.Hour)

	_, err := svc.Approve(ctx, post.UUID)
	if err != ErrDraftExpired {
		t.Errorf("se esperaba ErrDraftExpired, se obtuvo: %v", err)
	}
}

func TestBotChannelResolution(t *testing.T) {
	bot := &Bot{
		channels: ChannelConfig{
			DefaultChannelID:  "chan-default",
			ReleasesChannelID: "chan-releases",
			FeaturesChannelID: "chan-features",
		},
	}

	if ch := bot.resolveChannel("", "general"); ch != "chan-default" {
		t.Errorf("se esperaba fallback a chan-default, se obtuvo %s", ch)
	}
	if ch := bot.resolveChannel("chan-custom", "custom"); ch != "chan-custom" {
		t.Errorf("se esperaba canal configurado chan-custom, se obtuvo %s", ch)
	}
}
