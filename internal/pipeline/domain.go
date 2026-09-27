package pipeline

import "time"

type GitData struct {
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