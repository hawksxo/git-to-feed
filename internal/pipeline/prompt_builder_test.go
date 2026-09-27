package pipeline

import (
	"strings"
	"testing"
	"time"
)

func TestBuildPromptSuccess(t *testing.T) {
	pb := NewPromptBuilder()

	richCtx := RichContext{
		GitData: GitData{
			EventType:   "release",
			Tag:         "v1.0.0",
			Title:       "v1.0.0 Release",
			Description: "HMAC signature verification implemented.",
			Author:      "hawksxo",
			Repository:  "hawksxo/git-to-feed",
			URL:         "https://github.com/hawksxo/git-to-feed",
		},
		HumanNotes: &HumanNotes{
			Notes: "Importante destacar el aspecto de seguridad.",
		},
		CreatedAt: time.Now(),
	}

	prompt, err := pb.BuildPrompt(richCtx, ArchetypeRelease)
	if err != nil {
		t.Fatalf("BuildPrompt() unexpected error: %v", err)
	}

	if !strings.Contains(prompt, "[SYSTEM INSTRUCTION]") {
		t.Errorf("Expected prompt to contain [SYSTEM INSTRUCTION]")
	}

	if !strings.Contains(prompt, "v1.0.0 Release") {
		t.Errorf("Expected prompt to contain release title")
	}

	if !strings.Contains(prompt, "Importante destacar el aspecto de seguridad.") {
		t.Errorf("Expected prompt to contain human notes")
	}
}

func TestBuildPromptWithoutHumanNotes(t *testing.T) {
	pb := NewPromptBuilder()

	richCtx := RichContext{
		GitData: GitData{
			EventType:   "pull_request",
			Tag:         "",
			Title:       "Refactor HTTP router",
			Description: "Cleaned up route definitions.",
			Author:      "hawksxo",
			Repository:  "hawksxo/git-to-feed",
			URL:         "https://github.com/hawksxo/git-to-feed",
		},
		HumanNotes: nil,
		CreatedAt:  time.Now(),
	}

	prompt, err := pb.BuildPrompt(richCtx, ArchetypeFeature)
	if err != nil {
		t.Fatalf("BuildPrompt() unexpected error: %v", err)
	}

	if !strings.Contains(prompt, "[HUMAN NOTES]\nNone") {
		t.Errorf("Expected human notes section to be 'None'")
	}
}

func TestBuildPromptInvalidArchetype(t *testing.T) {
	pb := NewPromptBuilder()

	richCtx := RichContext{
		GitData: GitData{
			Title: "Test",
		},
	}

	prompt, err := pb.BuildPrompt(richCtx, Archetype("unknown"))
	if err == nil {
		t.Errorf("Expected error for invalid archetype, got nil")
	}

	if prompt != "" {
		t.Errorf("Expected empty prompt string on error, got %q", prompt)
	}
}
