package generator

import (
	"context"
	"errors"
	"time"
)

type Generator struct {
	client *GeminiLLMClient
}

func New(client *GeminiLLMClient) *Generator {
	return &Generator{
		client: client,
	}
}

func (g *Generator) Generate(ctx context.Context, richContext RichContext, archetype Archetype) (GeneratedPost, error) {
	template, err := GetTemplateForArchetype(archetype)
	if err != nil {
		return GeneratedPost{}, err
	}

	prompt, err := BuildPrompt(richContext, template)
	if err != nil {
		return GeneratedPost{}, err
	}

	rawContent, err := g.client.GeneratePostWithLLM(ctx, prompt)
	if err != nil {
		return GeneratedPost{}, err
	}

	postGenerated := ApplyAntiCringeFilter(rawContent)
	if postGenerated == "" {
		return GeneratedPost{}, errors.New("el post generado quedó vacío tras el filtro")
	}

	return GeneratedPost{
		Content: postGenerated, 
		Archetype: template.Archetype, 
		CreatedAt: time.Now(),
	}, nil
}