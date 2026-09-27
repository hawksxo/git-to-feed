package pipeline

import (
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
