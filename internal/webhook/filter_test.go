package webhook

import "testing"

func TestShouldProcessPullRequest_RejectMain(t *testing.T) {
	pr := PullRequestInfo{
		Title: "feat: test",
		Base:  BaseInfo{Ref: "main"},
	}

	got, _ := ShouldProcessPullRequest(pr)

	if got != false {
		t.Errorf("expected false for main branch, got true")
	}
}

func TestShouldProcessPullRequest_RejectNonDeliverablePrefix(t *testing.T) {
	pr := PullRequestInfo{
		Title: "docs: update README",
		Base:  BaseInfo{Ref: "develop"},
	}

	got, _ := ShouldProcessPullRequest(pr)

	if got != false {
		t.Errorf("expected false for docs prefix, got true")
	}
}

func TestShouldProcessPullRequest_AcceptValidPR(t *testing.T)  {
	pr := PullRequestInfo{
		Title: "feat: new test integration",
		Base: BaseInfo{Ref: "develop"},
	}

	got, _ := ShouldProcessPullRequest(pr)

	if got != true {
		t.Errorf("expected true for valid feat PR, got false")
	}
}