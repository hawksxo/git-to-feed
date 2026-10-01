package webhook

import (
	"fmt"
)

type ProcessWebhookUseCase struct {

}

func (u *ProcessWebhookUseCase) ProcessEvent(payload GitHubPayload) error {

	// STEP A: Validate rules for release
	if payload.EventType == "release" && payload.Action != "published" {
		return nil
	}

	// STEP B: Validate rules for pull_request
	if payload.EventType == "pull_request" && (payload.Action != "closed" || !payload.PullRequest.Merged) {
		return nil
	}

	// STEP C: Log captured event data
	fmt.Printf("Successfully received event data")

	return nil
}