package pipeline

import "fmt"

type PromptBuilder struct{}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

func (pb *PromptBuilder) BuildPrompt(ctx RichContext, archetype Archetype) (string, error) {
	notesText := "None"
	if ctx.HumanNotes != nil && ctx.HumanNotes.Notes != "" {
		notesText = ctx.HumanNotes.Notes
	}

	tmpl, err := GetTemplateForArchetype(archetype)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf("[SYSTEM INSTRUCTION]\n%s\n\n[FEW-SHOT EXAMPLE]\nInput: %s\nOutput: %s\n\n[TARGET CONTEXT]\nEvent Type: %s\nTag: %s\nTitle: %s\nRepository: %s\nAuthor: %s\nURL: %s\nDescription:\n%s\n\n[HUMAN NOTES]\n%s",
		tmpl.SystemInstruction,
		tmpl.Example.InputContext,
		tmpl.Example.ExpectedOutput,
		ctx.GitData.EventType,
		ctx.GitData.Tag,
		ctx.GitData.Title,
		ctx.GitData.Repository,
		ctx.GitData.Author,
		ctx.GitData.URL,
		ctx.GitData.Description,
		notesText,
	)

	return prompt, nil
}
