package pipeline

import (
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
}

func NewDefaultGenerator() *DefaultGenerator {
	return &DefaultGenerator{
		promptBuilder: NewPromptBuilder(),
		postProcessor: NewPostProcessor(),
	}
}

func (g *DefaultGenerator) Generate(ctx RichContext, archetype Archetype) (GeneratedPost, error) {
	var sb strings.Builder

	switch archetype {
	case ArchetypeRelease:
		sb.WriteString(fmt.Sprintf("🚀 ¡Lanzamos nueva versión de %s! (%s)\n\n", ctx.GitData.Repository, ctx.GitData.Tag))
		sb.WriteString(fmt.Sprintf("Nos alegra presentar %s.\n\n", ctx.GitData.Title))
		if ctx.GitData.Description != "" {
			sb.WriteString(fmt.Sprintf("📌 Cambios principales:\n%s\n\n", ctx.GitData.Description))
		}
		if ctx.GitData.URL != "" {
			sb.WriteString(fmt.Sprintf("💻 Revisa los detalles del release aquí: %s\n\n", ctx.GitData.URL))
		}
		sb.WriteString("#golang #cleanarchitecture #backend #opensource #devcommunity")

	case ArchetypeFeature:
		sb.WriteString(fmt.Sprintf("✨ Nueva funcionalidad implementada en %s!\n\n", ctx.GitData.Repository))
		sb.WriteString(fmt.Sprintf("Mejora clave: %s.\n\n", ctx.GitData.Title))
		if ctx.GitData.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", ctx.GitData.Description))
		}
		if ctx.GitData.URL != "" {
			sb.WriteString(fmt.Sprintf("🔗 Pull Request: %s\n\n", ctx.GitData.URL))
		}
		sb.WriteString("#golang #backend #cleanarchitecture #feature #softwareengineering")

	default:
		sb.WriteString(fmt.Sprintf("🛠️ Actualización en el proyecto %s: %s\n\n", ctx.GitData.Repository, ctx.GitData.Title))
		if ctx.GitData.Description != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", ctx.GitData.Description))
		}
		if ctx.GitData.URL != "" {
			sb.WriteString(fmt.Sprintf("🔗 Enlace: %s\n\n", ctx.GitData.URL))
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
