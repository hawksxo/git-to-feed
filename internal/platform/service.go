package platform

import (
	"context"
	"os"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/publisher"
	"github.com/hawksxo/git-to-feed/internal/webhook"
)

type EventNotifier interface {
	SendApprovalNotification(post *approval.ApprovalPost) error
}

type EventOrchestrator struct {
	approvalUseCase *approval.ApprovalUseCase
	postProcessor   *pipeline.PostProcessor
	publishUseCase  *publisher.PublishApprovedPostUseCase
	contextBuilder  *pipeline.ContextBuilder
	generator       *pipeline.DefaultGenerator
	notifier        EventNotifier
}

func NewEventOrchestrator(approvalUseCase *approval.ApprovalUseCase, postProcessor *pipeline.PostProcessor, publishUseCase *publisher.PublishApprovedPostUseCase) *EventOrchestrator {
	var geminiClient *pipeline.GeminiLLMClient
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	if geminiApiKey != "" {
		geminiClient, _ = pipeline.NewGeminiLLMClient(geminiApiKey, os.Getenv("GEMINI_MODEL_NAME"))
	}

	return &EventOrchestrator{
		approvalUseCase: approvalUseCase,
		postProcessor:   postProcessor,
		publishUseCase:  publishUseCase,
		contextBuilder:  pipeline.NewContextBuilder(),
		generator:       pipeline.NewDefaultGenerator(geminiClient),
	}
}

func (eo *EventOrchestrator) SetNotifier(notifier EventNotifier) {
	eo.notifier = notifier
}

func (eo *EventOrchestrator) ProcessGitHubPayload(ctx context.Context, payload webhook.GitHubPayload) (*approval.ApprovalPost, error) {
	var richCtx pipeline.RichContext
	var archetype pipeline.Archetype

	if payload.EventType == "release" || payload.Release.TagName != "" {
		richCtx = eo.contextBuilder.FromRelease(payload.Release, payload.Repository, payload.Sender, nil)
		archetype = pipeline.ArchetypeRelease
	} else {
		richCtx = eo.contextBuilder.FromPullRequest(payload.PullRequest, payload.Repository, payload.Sender, nil)
		archetype = pipeline.ArchetypeFeature
	}

	generatedPost, err := eo.generator.Generate(richCtx, archetype)
	if err != nil {
		return nil, err
	}

	approvalPost, err := eo.approvalUseCase.SubmitForApproval(ctx, generatedPost)
	if err != nil {
		return nil, err
	}

	if eo.notifier != nil {
		_ = eo.notifier.SendApprovalNotification(approvalPost)
	}

	return approvalPost, nil
}

func (eo *EventOrchestrator) ProcessGitHubEvent(ctx context.Context, rawText string, archetype pipeline.Archetype) (*approval.ApprovalPost, error) {
	generatedPost := eo.postProcessor.Process(rawText, archetype)
	approvalPost, err := eo.approvalUseCase.SubmitForApproval(ctx, generatedPost)
	if err != nil {
		return nil, err
	}

	if eo.notifier != nil {
		_ = eo.notifier.SendApprovalNotification(approvalPost)
	}

	return approvalPost, nil
}

func (eo *EventOrchestrator) ApproveAndPublish(ctx context.Context, uuid string) (*publisher.PublishedPost, error) {
	approvedPost, err := eo.approvalUseCase.Approve(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return eo.publishUseCase.Execute(ctx, approvedPost)
}