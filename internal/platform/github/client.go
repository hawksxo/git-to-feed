package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL string = "https://api.github.com"

type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
}

type compareResponse struct {
	Commits []struct {
		SHA     string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"commits"`
}

type PullRequestFile struct {
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
}

type RateLimitStatus struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
}

type rateLimitResponse struct {
	Resources struct {
		Core struct {
			Limit     int   `json:"limit"`
			Used      int   `json:"used"`
			Remaining int   `json:"remaining"`
			Reset     int64 `json:"reset"`
		}
	}
}

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

func NewClient(token string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		token:      token,
		baseURL:    defaultBaseURL,
	}
}

func (c *Client) SetBaseURL(url string) {
	c.baseURL = strings.TrimRight(url, "/")
}

func (c *Client) newRequest(ctx context.Context, method, endpoint string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("github: error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return req, nil
}

func (c *Client) GetRateLimitStatus(ctx context.Context) (*RateLimitStatus, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/rate_limit")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: error executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: rate limit returned unexpected status: %d", resp.StatusCode)
	}

	var data rateLimitResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("github: error decoding rate limit response: %w", err)
	}

	return &RateLimitStatus{
		Limit:     data.Resources.Core.Limit,
		Remaining: data.Resources.Core.Remaining,
		ResetAt:   time.Unix(data.Resources.Core.Reset, 0).UTC(),
	}, nil
}

func (c *Client) GetPullRequestFiles(ctx context.Context, owner, repo string, pullNumber int) ([]PullRequestFile, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/pulls/%d/files", owner, repo, pullNumber)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: error executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: get pr files returned unexpected status: %d", resp.StatusCode)
	}

	var files []PullRequestFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("github: error decoding pr files response: %w", err)
	}

	return files, nil
}

func (c *Client) CompareCommits(ctx context.Context, owner, repo, base, head string) ([]Commit, error) {
	endpoint := fmt.Sprintf("/repos/%s/%s/compare/%s...%s", owner, repo, base, head)
	req, err := c.newRequest(ctx, http.MethodGet, endpoint)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github: error executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github: get compare commits returned unexpected status: %d", resp.StatusCode)
	}

	var data compareResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("github: error decoding compare commits response: %w", err)
	}

	commits := make([]Commit, 0, len(data.Commits))
	for _, item := range data.Commits {
		commits = append(commits, Commit{
			SHA: item.SHA,
			Message: item.Commit.Message,
		})
	}

	return commits, nil
}