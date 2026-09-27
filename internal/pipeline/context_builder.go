package pipeline

import (
	"time"

	"github.com/hawksxo/git-to-feed/internal/webhook"
)

type ContextBuilder struct{}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{}
}

func (cb *ContextBuilder) FromRelease(release webhook.ReleaseInfo, repo webhook.RepositoryInfo, sender webhook.SenderInfo, notes *HumanNotes) RichContext {
	cleanDescription := SanitizeText(release.Body)

	return RichContext{
		GitData: GitData{
			Title:       release.Name,
			Description: cleanDescription,
			Author:      sender.User,
			Repository:  repo.FullName,
			URL:         release.HTMLURL,
		},
		HumanNotes: notes,
		CreatedAt:  time.Now(),
	}
}

func (cb *ContextBuilder) FromPullRequest(pr webhook.PullRequestInfo, repo webhook.RepositoryInfo, sender webhook.SenderInfo, notes *HumanNotes) RichContext  {
	cleanDescription := SanitizeText(pr.Body)
	
	return RichContext{
		GitData: GitData{
			Title: pr.Title,
			Description: cleanDescription,
			Author: sender.User,
			Repository: repo.FullName,
			URL: pr.HTMLURL,
		},
		HumanNotes: notes,
		CreatedAt: time.Now(),
	}
}
