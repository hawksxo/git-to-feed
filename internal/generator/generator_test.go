package generator

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestApplyAntiCringeFilterCliches(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected string
	}{
		{
			name:     "Debe reemplazar 'I am thrilled to announce' por 'Lanzamos'",
			raw:      "I am thrilled to announce the new feature.",
			expected: "Lanzamos the new feature.",
		},
		{
			name:     "Debe reemplazar 'game-changer' por 'mejora clave'",
			raw:      "This PR is a game-changer for performance.",
			expected: "This PR is a mejora clave for performance.",
		},
		{
			name:     "Debe reemplazar 'synergy' por 'integración'",
			raw:      "Better synergy between modules.",
			expected: "Better integración between modules.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApplyAntiCringeFilter(tt.raw)
			if got != tt.expected {
				t.Errorf("ApplyAntiCringeFilter() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestApplyAntiCringeFilterEmojiDensity(t *testing.T) {
	raw := "🚀🚀🚀🚀 Lanzamos la nueva versión 🎉🎉🎉"
	got := ApplyAntiCringeFilter(raw)

	if strings.Contains(got, "🚀🚀🚀") {
		t.Errorf("Expected emoji density to be capped at 2, got %q", got)
	}

	if strings.Contains(got, "🎉🎉🎉") {
		t.Errorf("Expected emoji density to be capped at 2, got %q", got)
	}
}

func TestSanitizeText(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected string
	}{
		{
			name:     "Debe eliminar comentarios HTML",
			raw:      "Hello <!-- comment --> world",
			expected: "Hello world",
		},
		{
			name:     "Debe eliminar checkboxes sin marcar",
			raw:      "Feature description\n- [ ] Task pending\n- [X] Task done",
			expected: "Feature description\n- [X] Task done",
		},
		{
			name:     "Debe eliminar espacios excesivos",
			raw:      "   Text with spaces   \n",
			expected: "Text with spaces",
		},
		{
			name:     "Debe normalizar múltiples saltos de línea consecutivos",
			raw:      "Paragraph 1\n\n\n\nParagraph 2",
			expected: "Paragraph 1\n\nParagraph 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeInput(tt.raw)
			if got != tt.expected {
				t.Errorf("SanitizeText() =%q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetTemplateForArchetype(t *testing.T) {
	tests := []struct {
		name      string
		archetype Archetype
		expectErr bool
	}{
		{
			name:      "Debe retornar plantilla valida para ArchetypeRelease",
			archetype: ArchetypeRelease,
			expectErr: false,
		},
		{
			name:      "Debe retornar plantilla valida para ArchetypeFeature",
			archetype: ArchetypeFeature,
			expectErr: false,
		},
		{
			name:      "Debe retornar plantilla valida para ArchetypeRefactor",
			archetype: ArchetypeRefactor,
			expectErr: false,
		},
		{
			name:      "Debe retornar error para arquetipo invalido",
			archetype: Archetype("invalid_archetype"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := GetTemplateForArchetype(tt.archetype)

			if tt.expectErr {
				if err == nil {
					t.Errorf("GetTemplateForArchetype(%q) expected error, got nil", tt.archetype)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetTemplateForArchetype(%q) unexpected error: %v", tt.archetype, err)
			}

			if tmpl.Archetype != tt.archetype {
				t.Errorf("tmpl.Archetype = %q, want %q", tmpl.Archetype, tt.archetype)
			}

			if tmpl.SystemInstruction == "" {
				t.Errorf("Expected SystemInstruction to be non-empty")
			}

			if tmpl.Example.InputContext == "" || tmpl.Example.ExpectedOutput == "" {
				t.Errorf("Expected FewShotExample input and output to be non-empty")
			}
		})
	}
}

func TestLiveGeminiGeneration(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("Saltando test de llamada real a Gemini: GEMINI_API_KEY no definida")
	}

	modelName := os.Getenv("GEMINI_MODEL_NAME")
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}

	client, err := NewGeminiLLMClient(apiKey, modelName)
	if err != nil {
		t.Fatalf("Error inicializando GeminiLLMClient: %v", err)
	}

	gen := New(client)

	richCtx := RichContext{
		GitData: GitData{
			EventType:   "release",
			Tag:         "v2.0.0",
			Title:       "Arquitectura limpia e idiomática en Go (Package-by-Feature)",
			Description: "Reestructuración integral del backend eliminando capas redundantes de Clean Architecture a favor de Package-by-Feature, dependencias directas y funciones puras.",
			Author:      "hawksxo",
			Repository:  "hawksxo/git-to-feed",
			URL:         "https://github.com/hawksxo/git-to-feed/releases/tag/v2.0.0",
		},
		ProjectIdentity: "Microservicio orientado a eventos que transforma actividad de GitHub en publicaciones técnicas de LinkedIn.",
	}

	post, err := gen.Generate(context.Background(), richCtx, ArchetypeRelease)
	if err != nil {
		t.Fatalf("Error en gen.Generate con Gemini real: %v", err)
	}

	if post.Content == "" {
		t.Fatal("El post generado no debería estar vacío")
	}

	t.Logf("\n=================== POST GENERADO EN VIVO POR GEMINI ===================\n%s\n========================================================================\n", post.Content)
}
