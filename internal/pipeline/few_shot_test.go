package pipeline_test

import (
	"context"
	"testing"

	"github.com/hawksxo/git-to-feed/internal/pipeline"
)

func TestInMemoryFewShotRepository(t *testing.T) {
	ctx := context.Background()
	repo := pipeline.NewInMemoryFewShotRepository()

	example := &pipeline.FewShotExampleEntity{
		Archetype:      pipeline.ArchetypeRelease,
		InputContext:   "Release v1.0.0",
		ExpectedOutput: "🚀 ¡Lanzamos v1.0.0!",
		IsActive:       true,
	}

	t.Run("Debe guardar y recuperar ejemplos activos", func(t *testing.T) {
		err := repo.Save(ctx, example)
		if err != nil {
			t.Fatalf("se esperaba error nil, se obtuvo %v", err)
		}

		active, err := repo.FindActiveByArchetype(ctx, pipeline.ArchetypeRelease)
		if err != nil {
			t.Fatalf("se esperaba error nil en FindActiveByArchetype, se obtuvo %v", err)
		}
		if len(active) == 0 {
			t.Fatalf("se esperaba al menos 1 ejemplo activo")
		}
		if active[0].ExpectedOutput != example.ExpectedOutput {
			t.Errorf("se esperaba output %s, se obtuvo %s", example.ExpectedOutput, active[0].ExpectedOutput)
		}
	})

	t.Run("Debe fallar si los campos estan vacios", func(t *testing.T) {
		invalidEx := &pipeline.FewShotExampleEntity{Archetype: pipeline.ArchetypeRelease}
		err := repo.Save(ctx, invalidEx)
		if err != pipeline.ErrEmptyExampleFields {
			t.Errorf("se esperaba error %v, se obtuvo %v", pipeline.ErrEmptyExampleFields, err)
		}
	})
}
