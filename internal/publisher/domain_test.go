package publisher_test

import (
	"testing"

	"github.com/hawksxo/git-to-feed/internal/publisher"
)

func TestNewPublishedPost(t *testing.T) {
	t.Run("Debe crear un PublishedPost valido con estado PENDING", func(t *testing.T) {
		approvalUUID := "550e8400-e29b-41d4-a716-446655440000"
		post, err := publisher.NewPublishedPost(approvalUUID)

		if err != nil {
			t.Fatalf("se esperaba error nil, se obtuvo: %v", err)
		}
		if post.ID == "" {
			t.Errorf("se esperaba que ID no estuviera vacio")
		}
		if post.ApprovalPostUUID != approvalUUID {
			t.Errorf("se esperaba ApprovalPostUUID %s, se obtuvo %s", approvalUUID, post.ApprovalPostUUID)
		}
		if post.Status != publisher.PublishStatusPending {
			t.Errorf("se esperaba estado PENDING, se obtuvo %s", post.Status)
		}
	})

	t.Run("Debe retornar error si ApprovalPostUUID esta vacio", func(t *testing.T) {
		post, err := publisher.NewPublishedPost("")

		if err == nil {
			t.Fatalf("se esperaba error ErrEmptyApprovalUUID, se obtuvo nil")
		}
		if post != nil {
			t.Errorf("se esperaba post nil, se obtuvo: %v", post)
		}
		if err != publisher.ErrEmptyApprovalUUID {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrEmptyApprovalUUID, err)
		}
	})
}

func TestMarkAsPublished(t *testing.T) {
	t.Run("Debe marcar como PUBLISHED correctamente", func(t *testing.T) {
		post, _ := publisher.NewPublishedPost("550e8400-e29b-41d4-a716-446655440000")
		shareURN := "urn:li:share:123456789"

		err := post.MarkAsPublished(shareURN)

		if err != nil {
			t.Fatalf("se esperaba error nil, se obtuvo: %v", err)
		}
		if post.Status != publisher.PublishStatusPublished {
			t.Errorf("se esperaba estado PUBLISHED, se obtuvo %s", post.Status)
		}
		if post.LinkedInShareURN != shareURN {
			t.Errorf("se esperaba LinkedInShareURN %s, se obtuvo %s", shareURN, post.LinkedInShareURN)
		}
	})

	t.Run("Debe fallar al marcar como PUBLISHED si shareURN esta vacio", func(t *testing.T) {
		post, _ := publisher.NewPublishedPost("550e8400-e29b-41d4-a716-446655440000")

		err := post.MarkAsPublished("")

		if err != publisher.ErrEmptyShareURN {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrEmptyShareURN, err)
		}
	})

	t.Run("Debe fallar si la transicion de estado no es desde PENDING", func(t *testing.T) {
		post, _ := publisher.NewPublishedPost("550e8400-e29b-41d4-a716-446655440000")
		_ = post.MarkAsPublished("urn:li:share:123456789")

		err := post.MarkAsPublished("urn:li:share:987654321")

		if err != publisher.ErrInvalidStatusTransition {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrInvalidStatusTransition, err)
		}
	})
}

func TestMarkAsFailed(t *testing.T) {
	t.Run("Debe marcar como FAILED correctamente con razon", func(t *testing.T) {
		post, _ := publisher.NewPublishedPost("550e8400-e29b-41d4-a716-446655440000")
		reason := "LinkedIn API 401 Unauthorized"

		err := post.MarkAsFailed(reason)

		if err != nil {
			t.Fatalf("se esperaba error nil, se obtuvo: %v", err)
		}
		if post.Status != publisher.PublishStatusFailed {
			t.Errorf("se esperaba estado FAILED, se obtuvo %s", post.Status)
		}
		if post.ErrorMessage != reason {
			t.Errorf("se esperaba ErrorMessage %s, se obtuvo %s", reason, post.ErrorMessage)
		}
	})

	t.Run("Debe fallar si la razon esta vacia", func(t *testing.T) {
		post, _ := publisher.NewPublishedPost("550e8400-e29b-41d4-a716-446655440000")

		err := post.MarkAsFailed("")

		if err != publisher.ErrEmptyFailureReason {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrEmptyFailureReason, err)
		}
	})

	t.Run("Debe fallar si se intenta marcar FAILED desde un estado no PENDING", func(t *testing.T) {
		post, _ := publisher.NewPublishedPost("550e8400-e29b-41d4-a716-446655440000")
		_ = post.MarkAsFailed("error 1")

		err := post.MarkAsFailed("error 2")

		if err != publisher.ErrInvalidStatusTransition {
			t.Errorf("se esperaba error %v, se obtuvo %v", publisher.ErrInvalidStatusTransition, err)
		}
	})
}
