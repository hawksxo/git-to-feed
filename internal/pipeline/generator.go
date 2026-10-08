package pipeline

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Generator interface {
	Generate(ctx RichContext, archetype Archetype) (GeneratedPost, error)
}

type DefaultGenerator struct {
	promptBuilder *PromptBuilder
	postProcessor *PostProcessor
	geminiLLMClient * GeminiLLMClient
}

func NewDefaultGenerator(geminiLLMClient *GeminiLLMClient) *DefaultGenerator {
	return &DefaultGenerator{
		promptBuilder: NewPromptBuilder(),
		postProcessor: NewPostProcessor(),
		geminiLLMClient: geminiLLMClient,
	}
}

func (g *DefaultGenerator) Generate(ctx RichContext, archetype Archetype) (GeneratedPost, error) {
	if g.geminiLLMClient != nil {
		prompt, err := g.promptBuilder.BuildPrompt(ctx, archetype)
		if err == nil {
			rawContent, err := g.geminiLLMClient.GeneratePostWithLLM(context.Background(), prompt)
			if err == nil && strings.TrimSpace(rawContent) != "" {
				cleanContent := ApplyAntiCringeFilter(rawContent)
				return GeneratedPost{
					Content:   cleanContent,
					Archetype: archetype,
					CreatedAt: time.Now(),
				}, nil
			}
			if err != nil {
				fmt.Printf("⚠️ Excepción al invocar Gemini LLM: %v\n", err)
			}
		}
	}
	
	var sb strings.Builder

	switch archetype {
	case ArchetypeRelease:
		fmt.Fprintf(&sb, "🚀 ¡Lanzamos nueva versión de %s! (%s)\n\n", ctx.GitData.Repository, ctx.GitData.Tag)
		fmt.Fprintf(&sb, "Nos alegra presentar %s.\n\n", ctx.GitData.Title)
		if ctx.GitData.Description != "" {
			fmt.Fprintf(&sb, "📌 Cambios principales:\n%s\n\n", ctx.GitData.Description)
		}
		if ctx.GitData.URL != "" {
			fmt.Fprintf(&sb, "💻 Revisa los detalles del release aquí: %s\n\n", ctx.GitData.URL)
		}
		sb.WriteString("#golang #cleanarchitecture #backend #opensource #devcommunity")

	case ArchetypeFeature:
		fmt.Fprintf(&sb, "✨ Nueva funcionalidad implementada en %s!\n\n", ctx.GitData.Repository)
		fmt.Fprintf(&sb, "Mejora clave: %s.\n\n", ctx.GitData.Title)
		if ctx.GitData.Description != "" {
			fmt.Fprintf(&sb, "%s\n\n", ctx.GitData.Description)
		}
		if ctx.GitData.URL != "" {
			fmt.Fprintf(&sb, "🔗 Pull Request: %s\n\n", ctx.GitData.URL)
		}
		sb.WriteString("#golang #backend #cleanarchitecture #feature #softwareengineering")

	default:
		fmt.Fprintf(&sb, "🛠️ Actualización en el proyecto %s: %s\n\n", ctx.GitData.Repository, ctx.GitData.Title)
		if ctx.GitData.Description != "" {
			fmt.Fprintf(&sb, "%s\n\n", ctx.GitData.Description)
		}
		if ctx.GitData.URL != "" {
			fmt.Fprintf(&sb, "🔗 Enlace: %s\n\n", ctx.GitData.URL)
		}
		sb.WriteString("#golang #backend #code #dev")
	}

	formattedContent := sb.String()
	cleanContent := ApplyAntiCringeFilter(formattedContent)

	return GeneratedPost{
		Content:   cleanContent,
		Archetype: archetype,
		CreatedAt: time.Now(),
	}, nil
}
