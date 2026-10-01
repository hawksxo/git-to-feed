package pipeline

import (
	"testing"
)

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
