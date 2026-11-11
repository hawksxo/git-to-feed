package generator

import "time"

type GitData struct {
	EventType   string
	Tag         string
	Title       string
	Description string
	Author      string
	Repository  string
	URL         string
}

type HumanNotes struct {
	Notes string
}

type FileDiff struct {
	Path      string
	Additions int
	Deletions int
}

type DiffStats struct {
	TotalFiles int
	Additions  int
	Deletions  int
	Files      []FileDiff
}

type RichContext struct {
	GitData         GitData
	DiffStats       *DiffStats
	HumanNotes      *HumanNotes
	ProjectIdentity string
	CreatedAt       time.Time
}

type Archetype string

const (
	ArchetypeRelease  Archetype = "release"
	ArchetypeFeature  Archetype = "feature"
	ArchetypeRefactor Archetype = "refactor"
)

type FewShotExample struct {
	InputContext   string
	ExpectedOutput string
}

type FewShotExampleEntity struct {
	UUID           string    `json:"uuid"`
	Archetype      Archetype `json:"archetype"`
	InputContext   string    `json:"input_context"`
	ExpectedOutput string    `json:"expected_output"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

type PromptTemplate struct {
	Archetype         Archetype
	SystemInstruction string
	Example           FewShotExample
}

type GeneratedPost struct {
	Content   string
	Archetype Archetype
	CreatedAt time.Time
}

type clicheReplacement struct {
	phrase      string
	replacement string
}
