package pipeline

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

type GeneratedPost struct {
	Content   string
	Archetype Archetype
	CreatedAt time.Time
}
