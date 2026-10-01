package pipeline

import "time"

type PostProcessor struct{}

func NewPostProcessor() *PostProcessor {
	return &PostProcessor{}
}

func (pp *PostProcessor) Process(rawText string, archetype Archetype) GeneratedPost {
	cleanContent := ApplyAntiCringeFilter(rawText)

	return GeneratedPost{
		Content:   cleanContent,
		Archetype: archetype,
		CreatedAt: time.Now(),
	}
}