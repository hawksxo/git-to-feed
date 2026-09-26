package webhook

type OwnerInfo struct {
	Login string `json:"login"`
}

type RepositoryInfo struct {
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Owner       OwnerInfo `json:"owner"`
	Description string    `json:"description"`
	HTMLURL     string    `json:"html_url"`
}

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
}

type PullRequestInfo struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	Merged  bool   `json:"merged"`
	HTMLURL string `json:"html_url"`
}

type SenderInfo struct {
	User string `json:"login"`
}

type GitHubPayload struct {
	Action      string          `json:"action"`
	EventType   string          `json:"-"`
	Release     ReleaseInfo     `json:"release"`
	PullRequest PullRequestInfo `json:"pull_request"`
	Repository  RepositoryInfo  `json:"repository"`
	Sender      SenderInfo      `json:"sender"`
}
