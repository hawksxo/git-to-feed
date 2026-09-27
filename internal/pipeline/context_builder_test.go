package pipeline

import (
	"testing"

	"github.com/hawksxo/git-to-feed/internal/webhook"
)

func TestContextBuilderFromRelease(t *testing.T) {
	cb := NewContextBuilder()

	release := webhook.ReleaseInfo{
		Name: "v1.0.0",
		Body:    "<!-- comment --> Release Body",
		HTMLURL: "https://github.com/hawksxo/git-to-feed",
	}

	repository := webhook.RepositoryInfo{
		FullName: "hawksxo/git-to-feed",
	}

	sender := webhook.SenderInfo{
		User: "hawksxo",
	}

	note := &HumanNotes{
		Notes: "Nota importante de negocio",
	}

	ctx := cb.FromRelease(release, repository, sender, note)
	
	if ctx.GitData.Title != "v1.0.0" {
		t.Errorf("Title = %q, want %q", ctx.GitData.Title, "v1.0.0")
	}

	if ctx.GitData.Description != "Release Body" {
		t.Errorf("Description = %q, want %q", ctx.GitData.Description, "Release Body")
	}

	if ctx.GitData.Author != "hawksxo" {
		t.Logf("Coincidencia!")
	}

	if ctx.GitData.Repository != "hawksxo/git-to-feed" {
		t.Errorf("Repository = %q, want %q", ctx.GitData.Repository, "hawksxo/git-to-feed")
	}

	if ctx.GitData.URL != "https://github.com/hawksxo/git-to-feed" {
		t.Errorf("URL = %q, want %q", ctx.GitData.URL, "https://github.com/hawksxo/git-to-feed")
	}

	if ctx.HumanNotes == nil && ctx.HumanNotes.Notes != "Nota importante de negocio" {
		t.Errorf("HumanNotes unexpected value")
	}

	if ctx.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set, got zero time")
	}
}

func TestContextBuilderFromPullRequest(t *testing.T)  {
	cb := NewContextBuilder()

	pr := webhook.PullRequestInfo{
		Title: "Integration pipeline module",
		Body: "- [ ] Task pending\n- [X] Task done",
		HTMLURL: "https://github.com/hawksxo/git-to-feed",
	}

	repository := webhook.RepositoryInfo{
		FullName: "hawksxo/git-to-feed",
	}

	sender := webhook.SenderInfo{
		User: "hawksxo",
	}

	var note *HumanNotes = nil

	ctx := cb.FromPullRequest(pr, repository, sender, note)

	if ctx.GitData.Title != "Integration pipeline module" {
		t.Errorf("Title = %q, want %q", ctx.GitData.Title, "Integration pipeline module")
	}

	if ctx.GitData.Description != "- [X] Task done" {
		t.Errorf("Description = %q, want %q", ctx.GitData.Description, "- [X] Task done")
	}

	if ctx.GitData.URL != "https://github.com/hawksxo/git-to-feed" {
		t.Errorf("URL = %q, want %q", ctx.GitData.URL, "https://github.com/hawksxo/git-to-feed")
	}

	if ctx.HumanNotes != nil {
		t.Errorf("Expected HumanNotes to be present, got nil")
	}

	if ctx.CreatedAt.IsZero() {
		t.Errorf("Expected CreatedAt to be set, but got zero time")
	}
}
