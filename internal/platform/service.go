package platform

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hawksxo/git-to-feed/internal/approval"
	"github.com/hawksxo/git-to-feed/internal/pipeline"
	"github.com/hawksxo/git-to-feed/internal/platform/github"
	"github.com/hawksxo/git-to-feed/internal/publisher"
	"github.com/hawksxo/git-to-feed/internal/webhook"
)

type EventNotifier interface {
	SendApprovalNotification(post *approval.ApprovalPost) error
	SendSupersededNotification(post *approval.ApprovalPost, releaseTag string) error
}

type PullRequestFilesInspector interface {
	GetPullRequestFiles(ctx context.Context, owner, repo string, pullNumber int) ([]github.PullRequestFile, error)
}

type ReleaseComparator interface {
	CompareCommits(ctx context.Context, owner, repo, base, head string) ([]github.Commit, error)
}

type EventOrchestrator struct {
	approvalUseCase   *approval.ApprovalUseCase
	postProcessor     *pipeline.PostProcessor
	publishUseCase    *publisher.PublishApprovedPostUseCase
	contextBuilder    *pipeline.ContextBuilder
	generator         *pipeline.DefaultGenerator
	notifier          EventNotifier
	prInspector       PullRequestFilesInspector
	releaseComparator ReleaseComparator
}

func NewEventOrchestrator(approvalUseCase *approval.ApprovalUseCase, postProcessor *pipeline.PostProcessor,
	publishUseCase *publisher.PublishApprovedPostUseCase, fewShotRepo pipeline.FewShotRepository,
	prInspector PullRequestFilesInspector, releaseComparator ...ReleaseComparator) *EventOrchestrator {
	var geminiClient *pipeline.GeminiLLMClient
	geminiApiKey := os.Getenv("GEMINI_API_KEY")
	if geminiApiKey != "" {
		geminiClient, _ = pipeline.NewGeminiLLMClient(geminiApiKey, os.Getenv("GEMINI_MODEL_NAME"))
	}

	var comp ReleaseComparator
	if len(releaseComparator) > 0 {
		comp = releaseComparator[0]
	}

	return &EventOrchestrator{
		approvalUseCase:   approvalUseCase,
		postProcessor:     postProcessor,
		publishUseCase:    publishUseCase,
		contextBuilder:    pipeline.NewContextBuilder(),
		generator:         pipeline.NewDefaultGenerator(geminiClient, fewShotRepo),
		prInspector:       prInspector,
		releaseComparator: comp,
	}
}

func (eo *EventOrchestrator) SetNotifier(notifier EventNotifier) {
	eo.notifier = notifier
}

func (eo *EventOrchestrator) SetReleaseComparator(comp ReleaseComparator) {
	eo.releaseComparator = comp
}

func (eo *EventOrchestrator) ProcessGitHubPayload(ctx context.Context, payload webhook.GitHubPayload) (*approval.ApprovalPost, error) {
	var richCtx pipeline.RichContext
	var archetype pipeline.Archetype
	var idempotencyKey string

	if payload.EventType == "release" || payload.Release.TagName != "" {
		richCtx = eo.contextBuilder.FromRelease(payload.Release, payload.Repository, payload.Sender, nil)
		archetype = pipeline.ArchetypeRelease
		if payload.Repository.FullName != "" && payload.Release.TagName != "" {
			idempotencyKey = fmt.Sprintf("release:%s:%s", payload.Repository.FullName, payload.Release.TagName)

			if eo.releaseComparator != nil && payload.Release.TagName != "" {
				pendingPosts, err := eo.approvalUseCase.ListPending(ctx)
				if err == nil && len(pendingPosts) > 0 {
					var uuids []string
					for _, p := range pendingPosts {
						uuids = append(uuids, p.UUID)
					}
					_ = eo.approvalUseCase.SupersedeDrafts(ctx, uuids, payload.Release.TagName)
					if eo.notifier != nil {
						for _, p := range pendingPosts {
							_ = eo.notifier.SendSupersededNotification(&p, payload.Release.TagName)
						}
					}
				}
			}
		}
	} else {
		var diffStats *pipeline.DiffStats

		if eo.prInspector != nil && payload.PullRequest.Number > 0 {
			parts := strings.Split(payload.Repository.FullName, "/")

			if len(parts) == 2 {
				owner := parts[0]
				repo := parts[1]

				files, err := eo.prInspector.GetPullRequestFiles(ctx, owner, repo, payload.PullRequest.Number)
				if err == nil && len(files) > 0 {
					var totalAdd, totalDel int
					var fileDiffs []pipeline.FileDiff
					for _, f := range files {
						totalAdd += f.Additions
						totalDel += f.Deletions
						fileDiffs = append(fileDiffs, pipeline.FileDiff{
							Path:      f.Filename,
							Additions: f.Additions,
							Deletions: f.Deletions,
						})
					}
					diffStats = &pipeline.DiffStats{
						TotalFiles: len(files),
						Additions:  totalAdd,
						Deletions:  totalDel,
						Files:      fileDiffs,
					}
				}
			}
		}

		richCtx = eo.contextBuilder.FromPullRequest(payload.PullRequest, payload.Repository, payload.Sender, diffStats, nil)
		archetype = pipeline.ArchetypeFeature

		if payload.Repository.FullName != "" && payload.PullRequest.Title != "" {
			idempotencyKey = fmt.Sprintf("pr:%s:%s", payload.Repository.FullName, payload.PullRequest.Title)
		}
	}

	generatedPost, err := eo.generator.Generate(richCtx, archetype)
	if err != nil {
		return nil, err
	}

	approvalPost, err := eo.approvalUseCase.SubmitForApproval(ctx, generatedPost, idempotencyKey)
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

func (eo *EventOrchestrator) DiscardDraft(ctx context.Context, uuid string, reason string) (*approval.ApprovalPost, error) {
	rejectedPost, err := eo.approvalUseCase.Reject(ctx, uuid, reason)
	if err != nil {
		return nil, err
	}
	return rejectedPost, nil
}
