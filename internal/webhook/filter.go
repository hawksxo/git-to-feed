package webhook

import (
	"fmt"
	"strings"
)

func ShouldProcessPullRequest(pr PullRequestInfo) (bool, string)  {
	cleanTitle := strings.TrimSpace(strings.ToLower(pr.Title))
	isValidPrefix := strings.HasPrefix(cleanTitle, "feat:") || strings.HasPrefix(cleanTitle, "refactor:") || strings.HasPrefix(cleanTitle, "perf:")
	
	if pr.Base.Ref != "develop" {
		return false, fmt.Sprintf("target branch '%s' is not develop", pr.Base.Ref)
	}

	if !isValidPrefix {
		return false, fmt.Sprintf("PR title prefix '%s' is not deliverable", pr.Title)
	}

	return true, ""
}