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

type RichContext struct {
	GitData    GitData
	HumanNotes *HumanNotes
	CreatedAt  time.Time
}

type GeneratedPost struct {
	Content string
	Archetype Archetype
	CreatedAt time.Time
}