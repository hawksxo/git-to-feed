package platform

import (
	"context"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/publisher"
)

type EventOrchestrator struct {
	approvalUseCase *approval.ApprovalUseCase
	postProcessor   *pipeline.PostProcessor
	publishUseCase  *publisher.PublishApprovedPostUseCase
}

func NewEventOrchestrator(approvalUseCase *approval.ApprovalUseCase, postProcessor *pipeline.PostProcessor, publishUseCase *publisher.PublishApprovedPostUseCase) *EventOrchestrator {
	return &EventOrchestrator{approvalUseCase: approvalUseCase, postProcessor: postProcessor, publishUseCase: publishUseCase}
}

func (eo *EventOrchestrator) ProcessGitHubEvent(ctx context.Context, rawText string, archetype pipeline.Archetype) (*approval.ApprovalPost, error) {
	generatedPost := eo.postProcessor.Process(rawText, archetype)
	return eo.approvalUseCase.SubmitForApproval(ctx, generatedPost)
}

func (eo *EventOrchestrator) ApproveAndPublish(ctx context.Context, uuid string) (*publisher.PublishedPost, error) {
	approvedPost, err := eo.approvalUseCase.Approve(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return eo.publishUseCase.Execute(ctx, approvedPost)
}