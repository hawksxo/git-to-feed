package webhook_test

import (
	"testing"
	"github.com/hawksxo/git-to-feed/internal/webhook"
)

func TestExtractFeedNotes(t *testing.T) () {
	input := "Hola equipo\n/feed-notes Soporte para OAuth2\nPor favor revisar"
	got := webhook.ExtractFeedNotes(input)
	if got != "Soporte para OAuth2" {
		t.Errorf("got %q, want %q", got, "Soporte para OAuth2")
	}
}