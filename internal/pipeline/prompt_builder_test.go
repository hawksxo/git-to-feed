package pipeline

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestBuildPromptSuccess(t *testing.T) {
	pb := NewPromptBuilder(nil)

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

	prompt, err := pb.BuildPrompt(context.Background(), richCtx, ArchetypeRelease)
	if err != nil {
		t.Fatalf("BuildPrompt() unexpected error: %v", err)
	}

	if !strings.Contains(prompt, "[INSTRUCCIÓN DE SISTEMA]") {
		t.Errorf("Expected prompt to contain [INSTRUCCIÓN DE SISTEMA]")
	}

	if !strings.Contains(prompt, "Tipo de Evento: release") {
		t.Errorf("Expected prompt to contain Tipo de Evento: release")
	}

	if !strings.Contains(prompt, "Versión / Tag: v1.0.0") {
		t.Errorf("Expected prompt to contain Tag: v1.0.0")
	}

	if !strings.Contains(prompt, "v1.0.0 Release") {
		t.Errorf("Expected prompt to contain release title")
	}

	if !strings.Contains(prompt, "Repositorio: hawksxo/git-to-feed") {
		t.Errorf("Expected prompt to contain Repositorio")
	}

	if !strings.Contains(prompt, "Autor: hawksxo") {
		t.Errorf("Expected prompt to contain Autor: hawksxo")
	}

	if !strings.Contains(prompt, "URL de GitHub: https://github.com/hawksxo/git-to-feed") {
		t.Errorf("Expected prompt to contain URL")
	}

	if !strings.Contains(prompt, "Importante destacar el aspecto de seguridad.") {
		t.Errorf("Expected prompt to contain human notes")
	}
}

func TestBuildPromptWithoutHumanNotes(t *testing.T) {
	pb := NewPromptBuilder(nil)

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

	prompt, err := pb.BuildPrompt(context.Background(), richCtx, ArchetypeFeature)
	if err != nil {
		t.Fatalf("BuildPrompt() unexpected error: %v", err)
	}

	if !strings.Contains(prompt, "[NOTAS ADICIONALES DEL DESARROLLADOR]\nNinguna") {
		t.Errorf("Expected human notes section to be 'Ninguna'")
	}
}

func TestBuildPromptInvalidArchetype(t *testing.T) {
	pb := NewPromptBuilder(nil)

	richCtx := RichContext{
		GitData: GitData{
			Title: "Test",
		},
	}

	prompt, err := pb.BuildPrompt(context.Background(), richCtx, Archetype("unknown"))
	if err == nil {
		t.Errorf("Expected error for invalid archetype, got nil")
	}

	if prompt != "" {
		t.Errorf("Expected empty prompt string on error, got %q", prompt)
	}
}

func TestBuildPromptFewShotRepo(t *testing.T) {
	fewShotRepo := NewInMemoryFewShotRepository()
	err := fewShotRepo.Save(context.Background(), &FewShotExampleEntity{
		UUID:           "custom-01",
		Archetype:      ArchetypeRelease,
		InputContext:   "Custom Input Release Context",
		ExpectedOutput: "Custom Output: ¡Publicación Dinámica desde BD!",
		IsActive:       true,
	})
	if err != nil {
		t.Fatalf("unexpected error saving example: %v", err)
	}
	pb := NewPromptBuilder(fewShotRepo)

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

	prompt, err := pb.BuildPrompt(context.Background(), richCtx, ArchetypeRelease)
	if err != nil {
		t.Fatalf("unexpected error building prompt: %v", err)
	}

	if !strings.Contains(prompt, "Custom Output: ¡Publicación Dinámica desde BD!") {
		t.Errorf("expected prompt to contain dynamic db example, got: %s", prompt)
	}
}
