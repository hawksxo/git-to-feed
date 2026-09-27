package pipeline

import "testing"

func TestSanitizeText(t *testing.T) {
	tests := []struct {
		name string
		raw string
		expected string
	} {
		{
			name: "Debe eliminar comentarios HTML",
			raw: "Hello <!-- comment --> world",
			expected: "Hello world",
		},
		{
			name: "Debe eliminar checkboxes sin marcar",
			raw: "Feature description\n- [ ] Task pending\n- [X] Task done",
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
			got := SanitizeText(tt.raw)
			if got != tt.expected {
				t.Errorf("SanitizeText() =%q, want %q", got, tt.expected)
			}
		})
	}
}