package pipeline

import (
	"strings"
	"testing"
)

func TestPostProcessorProcess(t *testing.T) {
	pp := NewPostProcessor()

	rawInput := "🚀🚀🚀 I am thrilled to announce this game-changer feature for our module!"
	post := pp.Process(rawInput, ArchetypeRelease)

	if strings.Contains(post.Content, "I am thrilled to announce") {
		t.Errorf("Expected cliché to be filtered, got %q", post.Content)
	}

	if strings.Contains(post.Content, "🚀🚀🚀") {
		t.Errorf("Expected emoji density to be capped, got %q", post.Content)
	}

	if post.Archetype != ArchetypeRelease {
		t.Errorf("Archetype = %q, want %q", post.Archetype, ArchetypeRelease)
	}

	if post.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt timestamp to be set")
	}
}