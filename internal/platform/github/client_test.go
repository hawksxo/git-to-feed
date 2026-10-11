package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetRateLimitStatus(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"resources": {
				"core": {
					"limit": 5000,
					"remaining": 4999,
					"reset": 1728590000
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", 2*time.Second)
	client.SetBaseURL(server.URL)

	status, err := client.GetRateLimitStatus(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if status.Limit != 5000 {
		t.Errorf("expected limit 5000, got: %d", status.Limit)
	}

	if status.Remaining != 4999 {
		t.Errorf("expected remaining 4999, got: %d", status.Remaining)
	}
}

func TestGetPullRequestFiles_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{
				"filename": "file1.txt",
				"status": "added",
				"additions": 103,
				"deletions": 21,
				"changes": 124
			}
		]`))
	}))
	defer server.Close()

	client := NewClient("test-token", 2*time.Second)
	client.SetBaseURL(server.URL)

	files, err := client.GetPullRequestFiles(context.Background(), "owner", "repo", 7)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 file, got: %d", len(files))
	}

	if files[0].Filename != "file1.txt" {
		t.Errorf("expected filename file1.txt, got %s", files[0].Filename)
	}

	if files[0].Additions != 103 {
		t.Errorf("expected 103 additions, got %d", files[0].Additions)
	}
}

func TestCompareCommits_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"commits": [
				{
					"sha": "abc1234",
					"commit": {
						"message": "feat: test message"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient("test-token", 2*time.Second)
	client.SetBaseURL(server.URL)

	commits, err := client.CompareCommits(context.Background(), "owner", "repo", "v1.0.0", "v1.1.0")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(commits) != 1 {
		t.Errorf("expected 1 commit, got: %d", len(commits))
	}

	if commits[0].SHA != "abc1234" {
		t.Errorf("expected sha abc1234, got %s", commits[0].SHA)
	}

	if commits[0].Message != "feat: test message" {
		t.Errorf("expected 'feat: test message' message, got %s", commits[0].Message)
	}
}

func TestClient_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message": "Bad credentials"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", 2*time.Second)
	client.SetBaseURL(server.URL)

	_, err := client.GetRateLimitStatus(context.Background())
	if err == nil {
		t.Fatalf("expected error on HTTP 401, got nil")
	}
}
